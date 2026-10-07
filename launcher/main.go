// opencode-qwen: single-file launcher that bundles opencode and wires it to a
// local Ollama model (Qwen3.8 27B abliterated, 80K context).
package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"
	"unsafe"
)

// mode is injected at build time (-X main.mode=desktop). Empty = CLI/web launcher.
var mode string

// backend is injected at build time (-X main.backend=strata). Empty = Ollama (huihui-80k),
// "strata" = Qwen3.8-Flash-Next served by Strata on strataAddr.
var backend string

const (
	strataAddr  = "127.0.0.1:8090"
	strataModel = "qwen3.8-flash-next"
	strataCtx   = 131072
)

const (
	appVersion  = "1.18.35-1"
	ollamaAddr  = "127.0.0.1:11435" // dedicated port, independent of the desktop app (11434)
	legacyAddr  = "127.0.0.1:11434"
	modelName   = "huihui-80k"
	baseModel   = "huihui_ai/Qwen3.8-abliterated:27b"
	numCtx      = 81920
	keepAlive   = "30m"
	modelfileTP = `FROM %s
PARAMETER num_ctx %d
PARAMETER num_predict -1
PARAMETER temperature 0.7
PARAMETER top_p 0.8
PARAMETER top_k 20
PARAMETER min_p 0
PARAMETER repeat_penalty 1
`
)

var cleanup func()

// defaultArgs is injected at build time (-X main.defaultArgs=web,--port,4096)
// and used when the exe is started without arguments.
var defaultArgs string

// ---- Windows job object: child processes die with the launcher, even if the
// console window is simply closed. ----
type jobBasicLimit struct {
	PerProcessUserTimeLimit int64
	PerJobUserTimeLimit     int64
	LimitFlags              uint32
	MinimumWorkingSetSize   uintptr
	MaximumWorkingSetSize   uintptr
	ActiveProcessLimit      uint32
	Affinity                uintptr
	PriorityClass           uint32
	SchedulingClass         uint32
}
type jobIO struct{ a, b, c, d, e, f uint64 }
type jobExtLimit struct {
	Basic                                                                jobBasicLimit
	IO                                                                   jobIO
	ProcessMemoryLimit, JobMemoryLimit, PeakProcessMemory, PeakJobMemory uintptr
}

var job syscall.Handle

func initJob() {
	k := syscall.NewLazyDLL("kernel32.dll")
	h, _, _ := k.NewProc("CreateJobObjectW").Call(0, 0)
	if h == 0 {
		return
	}
	var info jobExtLimit
	info.Basic.LimitFlags = 0x2000 // JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE
	r, _, _ := k.NewProc("SetInformationJobObject").Call(h, 9, uintptr(unsafe.Pointer(&info)), unsafe.Sizeof(info))
	if r == 0 {
		return
	}
	job = syscall.Handle(h)
}

func addToJob(cmd *exec.Cmd) {
	if job == 0 || cmd == nil || cmd.Process == nil {
		return
	}
	h, err := syscall.OpenProcess(0x1F0FFF, false, uint32(cmd.Process.Pid)) // PROCESS_ALL_ACCESS
	if err != nil {
		return
	}
	defer syscall.CloseHandle(h)
	syscall.NewLazyDLL("kernel32.dll").NewProc("AssignProcessToJobObject").Call(uintptr(job), uintptr(h))
}

// freeStalePort kills a leftover opencode.exe (our own extracted copy) that still
// holds the web UI port, e.g. after the previous window was force-closed.
func freeStalePort(args []string, bin string) {
	port := ""
	for i, a := range args {
		if a == "--port" && i+1 < len(args) {
			port = args[i+1]
		}
	}
	if port == "" || port == "0" {
		return
	}
	conn, err := net.DialTimeout("tcp", "127.0.0.1:"+port, time.Second)
	if err != nil {
		return
	}
	conn.Close()
	say("端口 %s 已被占用,尝试清理上次残留的 opencode 进程", port)
	ps := fmt.Sprintf(`Get-CimInstance Win32_Process -Filter "Name='opencode.exe'" | Where-Object { $_.ExecutablePath -eq '%s' } | ForEach-Object { Stop-Process -Id $_.ProcessId -Force }`, bin)
	cmd := exec.Command("powershell", "-NoProfile", "-Command", ps)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	cmd.Run()
	time.Sleep(time.Second)
}

func setTitle(t string) {
	p, err := syscall.UTF16PtrFromString(t)
	if err == nil {
		syscall.NewLazyDLL("kernel32.dll").NewProc("SetConsoleTitleW").Call(uintptr(unsafe.Pointer(p)))
	}
}

func say(format string, a ...any) { fmt.Printf("[opencode-qwen] "+format+"\n", a...) }

func die(format string, a ...any) {
	say("错误: "+format, a...)
	if cleanup != nil {
		cleanup()
	}
	fmt.Println("按回车键退出...")
	bufio.NewReader(os.Stdin).ReadString('\n')
	os.Exit(1)
}

func setUTF8Console() {
	k := syscall.NewLazyDLL("kernel32.dll")
	k.NewProc("SetConsoleOutputCP").Call(65001)
	k.NewProc("SetConsoleCP").Call(65001)
}

func findOllama() string {
	if p := os.Getenv("OLLAMA_EXE"); p != "" {
		return p
	}
	if la := os.Getenv("LOCALAPPDATA"); la != "" {
		p := filepath.Join(la, "Programs", "Ollama", "ollama.exe")
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	if p, err := exec.LookPath("ollama.exe"); err == nil {
		return p
	}
	return ""
}

func modelsDir() string {
	if d := os.Getenv("OLLAMA_MODELS"); d != "" {
		return d
	}
	// user-level env var (set via System Properties) may not be inherited by this process
	if out, err := exec.Command("reg", "query", `HKCU\Environment`, "/v", "OLLAMA_MODELS").Output(); err == nil {
		for _, line := range strings.Split(string(out), "\n") {
			f := strings.Fields(line)
			if len(f) >= 3 && f[0] == "OLLAMA_MODELS" {
				return strings.TrimSpace(strings.Join(f[2:], " "))
			}
		}
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".ollama", "models")
}

var httpc = &http.Client{Timeout: 5 * time.Second}

func up(addr string) bool {
	r, err := httpc.Get("http://" + addr + "/api/version")
	if err != nil {
		return false
	}
	r.Body.Close()
	return r.StatusCode == 200
}

func hasModel(name string) bool {
	r, err := httpc.Get("http://" + ollamaAddr + "/api/tags")
	if err != nil {
		return false
	}
	defer r.Body.Close()
	var t struct {
		Models []struct {
			Name string `json:"name"`
		} `json:"models"`
	}
	if json.NewDecoder(r.Body).Decode(&t) != nil {
		return false
	}
	for _, m := range t.Models {
		if m.Name == name || m.Name == name+":latest" {
			return true
		}
	}
	return false
}

func envFor(extra ...string) []string {
	return append(os.Environ(), extra...)
}

func ollamaEnv(models string) []string {
	return envFor(
		"OLLAMA_HOST="+ollamaAddr,
		"OLLAMA_MODELS="+models,
		"OLLAMA_FLASH_ATTENTION=1",
		"OLLAMA_KV_CACHE_TYPE=q8_0",
		"OLLAMA_KEEP_ALIVE="+keepAlive,
		"OLLAMA_NUM_PARALLEL=1",
		"OLLAMA_MAX_LOADED_MODELS=1",
	)
}

// freeLegacyServer unloads models held by another Ollama server (e.g. the
// desktop app on 11434) so that the 27B model fits fully into VRAM.
func freeLegacyServer() { freeServer(legacyAddr) }

// freeServer unloads every model held by the Ollama server at addr.
func freeServer(addr string) {
	r, err := httpc.Get("http://" + addr + "/api/ps")
	if err != nil {
		return
	}
	defer r.Body.Close()
	var ps struct {
		Models []struct {
			Name string `json:"name"`
		} `json:"models"`
	}
	if json.NewDecoder(r.Body).Decode(&ps) != nil {
		return
	}
	for _, m := range ps.Models {
		say("释放 %s 上已加载的模型 %s (避免显存不足)", addr, m.Name)
		body, _ := json.Marshal(map[string]any{"model": m.Name, "keep_alive": 0})
		resp, err := httpc.Post("http://"+addr+"/api/generate", "application/json", bytes.NewReader(body))
		if err == nil {
			resp.Body.Close()
		}
	}
	time.Sleep(2 * time.Second)
}

func startOllama(exe, models string) *exec.Cmd {
	cmd := exec.Command(exe, "serve")
	cmd.Env = ollamaEnv(models)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true, CreationFlags: 0x08000000} // CREATE_NO_WINDOW
	if err := cmd.Start(); err != nil {
		die("无法启动 Ollama: %v", err)
	}
	addToJob(cmd)
	for i := 0; i < 60; i++ {
		if up(ollamaAddr) {
			return cmd
		}
		time.Sleep(time.Second)
	}
	killTree(cmd)
	die("Ollama 服务 60 秒内没有就绪")
	return nil
}

func killTree(cmd *exec.Cmd) {
	if cmd == nil || cmd.Process == nil {
		return
	}
	exec.Command("taskkill", "/PID", fmt.Sprint(cmd.Process.Pid), "/T", "/F").Run()
}

func ensureModel(exe, models string) {
	if hasModel(modelName) {
		return
	}
	say("未找到模型 %s,准备创建 (上下文 %d)", modelName, numCtx)
	if !hasModel(baseModel) {
		fmt.Printf("基础模型 %s 不存在,需要下载约 17GB。现在下载吗? [y/N] ", baseModel)
		ans, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		if !strings.HasPrefix(strings.ToLower(strings.TrimSpace(ans)), "y") {
			die("没有基础模型,无法继续")
		}
		c := exec.Command(exe, "pull", baseModel)
		c.Env = ollamaEnv(models)
		c.Stdout, c.Stderr, c.Stdin = os.Stdout, os.Stderr, os.Stdin
		if err := c.Run(); err != nil {
			die("下载模型失败: %v", err)
		}
	}
	tmp := filepath.Join(os.TempDir(), "Modelfile.huihui-80k")
	if err := os.WriteFile(tmp, []byte(fmt.Sprintf(modelfileTP, baseModel, numCtx)), 0o644); err != nil {
		die("写 Modelfile 失败: %v", err)
	}
	defer os.Remove(tmp)
	c := exec.Command(exe, "create", modelName, "-f", tmp)
	c.Env = ollamaEnv(models)
	c.Stdout, c.Stderr = os.Stdout, os.Stderr
	if err := c.Run(); err != nil {
		die("创建模型失败: %v", err)
	}
}

func preload() {
	say("正在把模型加载进显存 (首次约 10~20 秒)...")
	body, _ := json.Marshal(map[string]any{"model": modelName, "keep_alive": keepAlive})
	c := &http.Client{Timeout: 5 * time.Minute}
	resp, err := c.Post("http://"+ollamaAddr+"/api/generate", "application/json", bytes.NewReader(body))
	if err != nil {
		say("预加载失败(可忽略,首次提问时会自动加载): %v", err)
		return
	}
	io.Copy(io.Discard, resp.Body)
	resp.Body.Close()
}

func extractFile(dst string, data []byte) error {
	if st, err := os.Stat(dst); err == nil && st.Size() == int64(len(data)) {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	tmp := dst + ".tmp"
	if err := os.WriteFile(tmp, data, 0o755); err != nil {
		return err
	}
	os.Remove(dst)
	return os.Rename(tmp, dst)
}

// findDesktop locates the OpenCode desktop app (Electron).
func findDesktop() string {
	var c []string
	if p := os.Getenv("OPENCODE_DESKTOP_EXE"); p != "" {
		c = append(c, p)
	}
	if self, err := os.Executable(); err == nil {
		d := filepath.Dir(self)
		c = append(c,
			filepath.Join(d, "desktop", "app", "OpenCode.exe"),
			filepath.Join(d, "..", "desktop", "app", "OpenCode.exe"),
			filepath.Join(d, "OpenCode.exe"),
		)
	}
	if la := os.Getenv("LOCALAPPDATA"); la != "" {
		c = append(c, filepath.Join(la, "Programs", "OpenCode", "OpenCode.exe"))
	}
	for _, p := range c {
		if st, err := os.Stat(p); err == nil && !st.IsDir() {
			abs, _ := filepath.Abs(p)
			return abs
		}
	}
	return ""
}

// desktopRunning reports whether any process of the given desktop exe is alive.
// (Electron is single-instance: a second launch hands over to the first and exits.)
func desktopRunning(exe string) bool {
	ps := fmt.Sprintf(`if (Get-Process -Name OpenCode -ErrorAction SilentlyContinue | Where-Object { $_.Path -eq '%s' }) { 'yes' }`, exe)
	cmd := exec.Command("powershell", "-NoProfile", "-Command", ps)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	out, err := cmd.Output()
	return err == nil && strings.Contains(string(out), "yes")
}

// findOpenWebUI returns the open-webui executable and its root folder.
func findOpenWebUI() (exe, root string) {
	var roots []string
	if d := os.Getenv("OPENWEBUI_DIR"); d != "" {
		roots = append(roots, d)
	}
	if self, err := os.Executable(); err == nil {
		d := filepath.Dir(self)
		roots = append(roots, filepath.Join(d, "openwebui"), filepath.Join(d, "..", "openwebui"))
	}
	for _, r := range roots {
		p := filepath.Join(r, "venv", "Scripts", "open-webui.exe")
		if st, err := os.Stat(p); err == nil && !st.IsDir() {
			abs, _ := filepath.Abs(r)
			return filepath.Join(abs, "venv", "Scripts", "open-webui.exe"), abs
		}
	}
	return "", ""
}

func portBusy(port string) bool {
	c, err := net.DialTimeout("tcp", "127.0.0.1:"+port, 500*time.Millisecond)
	if err != nil {
		return false
	}
	c.Close()
	return true
}

func openBrowser(url string) {
	cmd := exec.Command("cmd", "/c", "start", "", url)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	cmd.Start()
}

func runOpenWebUI() int {
	exe, root := findOpenWebUI()
	if exe == "" {
		die("没有找到 Open WebUI (openwebui\\venv\\Scripts\\open-webui.exe)。请先安装,或设置环境变量 OPENWEBUI_DIR")
	}
	data := filepath.Join(root, "data")
	os.MkdirAll(data, 0o755)

	port := os.Getenv("OPENWEBUI_PORT")
	if port == "" {
		port = "8080"
		for i := 0; i < 10 && portBusy(port); i++ {
			port = fmt.Sprint(8081 + i)
		}
	}
	url := "http://127.0.0.1:" + port

	say("启动 Open WebUI (%s),首次启动可能需要 1~3 分钟...", url)
	say("关闭此窗口即停止 Open WebUI 和 Ollama")
	cmd := exec.Command(exe, "serve", "--host", "127.0.0.1", "--port", port)
	cmd.Dir = data // the generated secret key lives here
	cmd.Env = envFor(
		"DATA_DIR="+data,
		"OLLAMA_BASE_URL=http://"+ollamaAddr,
		"ENABLE_OPENAI_API=False",
		"WEBUI_AUTH=False", // single local user, bound to 127.0.0.1 only
		"DEFAULT_MODELS="+modelName+":latest",
		"ENABLE_EVALUATION_ARENA_MODELS=False",
		"WEBUI_NAME=Open WebUI - Qwen",
		"ENABLE_VERSION_UPDATE_CHECK=False",
		"ANONYMIZED_TELEMETRY=False",
		"SCARF_NO_ANALYTICS=true",
		"DO_NOT_TRACK=true",
		"HF_HUB_DISABLE_TELEMETRY=1",
		"PYTHONUTF8=1",
		"PYTHONIOENCODING=utf-8",
	)
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	if err := cmd.Start(); err != nil {
		say("启动 Open WebUI 失败: %v", err)
		return 1
	}
	addToJob(cmd)

	done := make(chan struct{})
	go func() {
		client := &http.Client{Timeout: 2 * time.Second}
		for i := 0; i < 600; i++ {
			select {
			case <-done:
				return
			default:
			}
			if r, err := client.Get(url + "/health"); err == nil {
				r.Body.Close()
				if r.StatusCode == 200 {
					say("Open WebUI 已就绪: %s (正在打开浏览器)", url)
					openBrowser(url)
					return
				}
			}
			time.Sleep(time.Second)
		}
		say("等待 Open WebUI 就绪超时,请查看上方日志")
	}()

	err := cmd.Wait()
	close(done)
	if ee, ok := err.(*exec.ExitError); ok {
		return ee.ExitCode()
	}
	return 0
}

// findStrata locates the Strata checkout (strata\Strata-main) and its generated run-<model>.bat.
func findStrata() (root, script string) {
	var roots []string
	if d := os.Getenv("STRATA_DIR"); d != "" {
		roots = append(roots, d)
	}
	if self, err := os.Executable(); err == nil {
		d := filepath.Dir(self)
		roots = append(roots, filepath.Join(d, "strata", "Strata-main"), filepath.Join(d, "..", "strata", "Strata-main"))
	}
	for _, r := range roots {
		if st, err := os.Stat(filepath.Join(r, "setup.py")); err != nil || st.IsDir() {
			continue
		}
		abs, _ := filepath.Abs(r)
		if s := os.Getenv("STRATA_RUN"); s != "" {
			return abs, s
		}
		m, _ := filepath.Glob(filepath.Join(abs, "run-*.bat"))
		for _, p := range m {
			if strings.HasPrefix(strings.ToLower(filepath.Base(p)), "run-opencode-") {
				continue // our own generated copies
			}
			if strings.Contains(strings.ToLower(filepath.Base(p)), "iq2_xs") {
				return abs, filepath.Base(p)
			}
		}
		if len(m) > 0 {
			return abs, filepath.Base(m[0])
		}
		return abs, ""
	}
	return "", ""
}

// strataHealth reports whether the Strata server answers and has the model loaded.
func strataHealth() (ok, loaded bool) {
	c := &http.Client{Timeout: 3 * time.Second}
	r, err := c.Get("http://" + strataAddr + "/health")
	if err != nil {
		return false, false
	}
	defer r.Body.Close()
	if r.StatusCode != 200 {
		return false, false
	}
	var h struct {
		Loaded bool `json:"loaded"`
	}
	json.NewDecoder(r.Body).Decode(&h)
	return true, h.Loaded
}

// runStrata makes room in VRAM (unloads Ollama models) and starts the Strata server in the
// background; returns the process (nil when an already running server is reused).
func runStrata() *exec.Cmd {
	root, script := findStrata()
	if root == "" {
		die("没有找到 Strata (strata\\Strata-main)。请先安装,或设置环境变量 STRATA_DIR")
	}
	if ok, loaded := strataHealth(); ok && loaded {
		say("复用已运行的 Strata %s", strataAddr)
		return nil
	} else if ok || portBusy("8090") {
		die("端口 8090 已被占用,但不是就绪的 Strata。请先关闭占用它的程序")
	}
	if script == "" {
		die("Strata 还没有生成启动脚本 run-*.bat,请先运行安装 (START-HERE.bat --setup)")
	}
	// Strata fills the GPU with experts: no Ollama model may stay loaded.
	freeServer(ollamaAddr)
	freeServer(legacyAddr)

	// Setup's run script opens the browser (--open); make a copy without it for our own use.
	if b, err := os.ReadFile(filepath.Join(root, script)); err == nil {
		s := strings.ReplaceAll(string(b), ` "--open"`, "")
		quiet := "run-opencode-" + strings.TrimPrefix(script, "run-")
		if os.WriteFile(filepath.Join(root, quiet), []byte(s), 0o644) == nil {
			script = quiet
		}
	}
	logPath := filepath.Join(root, "launcher-strata.log")
	lf, _ := os.Create(logPath)
	say("启动 Strata (%s),首次加载约 1~3 分钟,期间电脑可能变卡,请不要关闭窗口", script)
	cmd := exec.Command("cmd", "/c", script)
	cmd.Dir = root
	cmd.Env = envFor("PYTHONUTF8=1")
	if lf != nil {
		cmd.Stdout, cmd.Stderr = lf, lf
	}
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true, CreationFlags: 0x08000000} // CREATE_NO_WINDOW
	if err := cmd.Start(); err != nil {
		die("无法启动 Strata: %v", err)
	}
	addToJob(cmd)

	exited := make(chan struct{})
	go func() { cmd.Wait(); close(exited) }()
	start := time.Now()
	for tick := 1; ; tick++ {
		select {
		case <-exited:
			die("Strata 意外退出,日志: %s", logPath)
		case <-time.After(2 * time.Second):
		}
		if _, loaded := strataHealth(); loaded {
			say("Strata 已就绪 (用时 %d 秒): http://%s/v1", int(time.Since(start).Seconds()), strataAddr)
			return cmd
		}
		if time.Since(start) > 15*time.Minute {
			killTree(cmd)
			die("Strata 15 分钟内没有就绪,日志: %s", logPath)
		}
		if tick%10 == 0 {
			say("等待 Strata 加载模型... %d 秒", int(time.Since(start).Seconds()))
		}
	}
}

// strataSettingsPath returns the Claude Code settings file that points at Strata.
// A claude.strata.settings.json next to the exe wins; otherwise the embedded copy
// is written under %LOCALAPPDATA%. The user's ~/.claude/settings.json is never touched.
func strataSettingsPath() (string, bool, error) {
	const cfgName = "claude.strata.settings.json"
	if self, err := os.Executable(); err == nil {
		ext := filepath.Join(filepath.Dir(self), cfgName)
		if st, err := os.Stat(ext); err == nil && !st.IsDir() {
			return ext, true, nil
		}
	}
	base := filepath.Join(os.Getenv("LOCALAPPDATA"), "opencode-qwen", appVersion)
	if err := os.MkdirAll(base, 0o755); err != nil {
		return "", false, err
	}
	cfg := filepath.Join(base, cfgName)
	if err := os.WriteFile(cfg, claudeConfig, 0o644); err != nil {
		return "", false, err
	}
	return cfg, false, nil
}

func isClaudeWrap() bool {
	if mode == "claudewrap" {
		return true
	}
	return strings.Contains(strings.ToLower(filepath.Base(os.Args[0])), "claude-wrap")
}

// findClaudeNative locates the real Claude Code executable, never this wrapper.
func findClaudeNative() string {
	self, _ := os.Executable()
	self, _ = filepath.Abs(self)
	var cands []string
	if p := os.Getenv("CLAUDE_EXE"); p != "" {
		cands = append(cands, p)
	}
	if app := os.Getenv("APPDATA"); app != "" {
		cands = append(cands, filepath.Join(app, "npm", "node_modules", "@anthropic-ai", "claude-code", "bin", "claude.exe"))
	}
	if home, err := os.UserHomeDir(); err == nil {
		cands = append(cands,
			filepath.Join(home, ".local", "bin", "claude.exe"),
			filepath.Join(home, ".claude", "local", "claude.exe"),
		)
	}
	for _, p := range cands {
		st, err := os.Stat(p)
		if err != nil || st.IsDir() {
			continue
		}
		abs, _ := filepath.Abs(p)
		if self != "" && strings.EqualFold(abs, self) {
			continue
		}
		return abs
	}
	return ""
}

// runClaudeWrap is spawned by the web UI in place of claude. It prepends
// --settings and forwards every byte of stdio. Nothing else may be written
// to stdout: the web UI parses stdout as JSON.
func runClaudeWrap() int {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "claude-wrap is a helper. Use claude-strata-ui.exe or claude-strata.exe.")
		return 1
	}
	cfg, _, err := strataSettingsPath()
	if err != nil {
		fmt.Fprintf(os.Stderr, "claude-wrap: %v\n", err)
		return 1
	}
	exe := findClaudeNative()
	if exe == "" {
		fmt.Fprintln(os.Stderr, "claude-wrap: claude.exe not found. Install Claude Code with: npm i -g @anthropic-ai/claude-code")
		return 1
	}
	args := os.Args[1:]
	// The web UI probes with --version. Skip --settings so stdout stays a single version line.
	onlyVersion := len(args) == 1 && (args[0] == "--version" || args[0] == "-v" || args[0] == "-V")
	if !onlyVersion {
		args = append([]string{"--settings", cfg}, args...)
	}
	streamJSON := false
	for _, a := range args {
		if a == "stream-json" {
			streamJSON = true
			break
		}
	}
	cmd := exec.Command(exe, args...)
	cmd.Stdin, cmd.Stderr = os.Stdin, os.Stderr
	cmd.Env = envFor("NO_PROXY=127.0.0.1,localhost", "no_proxy=127.0.0.1,localhost")
	if streamJSON {
		// Claude Code emits one system message per thinking tick. The web UI
		// draws every system message, so drop those ticks and keep the reply.
		stdout, pipeErr := cmd.StdoutPipe()
		if pipeErr != nil {
			fmt.Fprintf(os.Stderr, "claude-wrap: %v\n", pipeErr)
			return 1
		}
		if err = cmd.Start(); err != nil {
			fmt.Fprintf(os.Stderr, "claude-wrap: %v\n", err)
			return 1
		}
		forwardClaudeStdout(os.Stdout, stdout)
		err = cmd.Wait()
	} else {
		cmd.Stdout = os.Stdout
		err = cmd.Run()
	}
	if ee, ok := err.(*exec.ExitError); ok {
		return ee.ExitCode()
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "claude-wrap: %v\n", err)
		return 1
	}
	return 0
}

// forwardClaudeStdout copies Claude's JSON lines, dropping thinking-progress ticks.
func forwardClaudeStdout(dst io.Writer, src io.Reader) {
	sc := bufio.NewScanner(src)
	sc.Buffer(make([]byte, 0, 64*1024), 16*1024*1024)
	for sc.Scan() {
		line := sc.Bytes()
		if skipThinkingTick(line) {
			continue
		}
		dst.Write(line)
		dst.Write([]byte("\n"))
	}
	if err := sc.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "claude-wrap: read stdout: %v\n", err)
	}
}

func skipThinkingTick(line []byte) bool {
	if !bytes.Contains(line, []byte(`"thinking_tokens"`)) {
		return false
	}
	var m struct {
		Type    string `json:"type"`
		Subtype string `json:"subtype"`
	}
	if json.Unmarshal(line, &m) != nil {
		return false
	}
	return m.Type == "system" && m.Subtype == "thinking_tokens"
}

func findNodeExe() string {
	if p := os.Getenv("CLAUDE_NODE"); p != "" {
		if st, err := os.Stat(p); err == nil && !st.IsDir() {
			return p
		}
	}
	var cands []string
	if pf := os.Getenv("ProgramFiles"); pf != "" {
		cands = append(cands, filepath.Join(pf, "nodejs", "node.exe"))
	}
	if la := os.Getenv("LOCALAPPDATA"); la != "" {
		cands = append(cands, filepath.Join(la, "Programs", "nodejs", "node.exe"))
	}
	for _, p := range cands {
		if st, err := os.Stat(p); err == nil && !st.IsDir() {
			return p
		}
	}
	out, err := exec.Command("where.exe", "node.exe").Output()
	if err == nil {
		for _, line := range strings.Split(string(out), "\n") {
			line = strings.TrimSpace(line)
			if strings.HasSuffix(strings.ToLower(line), "node.exe") {
				return line
			}
		}
	}
	return ""
}

func findWebUIScript() string {
	if p := os.Getenv("CLAUDE_WEBUI"); p != "" {
		if st, err := os.Stat(p); err == nil && !st.IsDir() {
			abs, _ := filepath.Abs(p)
			return abs
		}
	}
	rel := filepath.Join("claudeui", "node_modules", "claude-code-webui", "dist", "cli", "node.js")
	var roots []string
	if self, err := os.Executable(); err == nil {
		d := filepath.Dir(self)
		roots = append(roots, d, filepath.Dir(d))
	}
	for _, r := range roots {
		p := filepath.Join(r, rel)
		if st, err := os.Stat(p); err == nil && !st.IsDir() {
			abs, _ := filepath.Abs(p)
			return abs
		}
	}
	return ""
}

func findWrapExe() string {
	self, err := os.Executable()
	if err != nil {
		return ""
	}
	p := filepath.Join(filepath.Dir(self), "claude-wrap.exe")
	st, err := os.Stat(p)
	if err != nil || st.IsDir() {
		return ""
	}
	abs, _ := filepath.Abs(p)
	return abs
}

// pickProjectFolder opens a Windows folder dialog. Cancel returns "".
func pickProjectFolder() string {
	out := filepath.Join(os.TempDir(), "claude-strata-folder.txt")
	_ = os.Remove(out)
	script := filepath.Join(os.TempDir(), "claude-strata-folder.ps1")
	body := "\uFEFF" + `$ErrorActionPreference = 'Stop'
Add-Type -AssemblyName System.Windows.Forms
$dlg = New-Object System.Windows.Forms.FolderBrowserDialog
$dlg.Description = '选择 Claude Code 的项目文件夹。取消则打开最近用过的项目。'
$dlg.ShowNewFolderButton = $true
if ($env:PICK_START -and (Test-Path -LiteralPath $env:PICK_START)) { $dlg.SelectedPath = $env:PICK_START }
$r = $dlg.ShowDialog()
if ($r -ne [System.Windows.Forms.DialogResult]::OK) { exit 2 }
$utf8 = New-Object System.Text.UTF8Encoding $false
[System.IO.File]::WriteAllText($env:PICK_OUT, $dlg.SelectedPath, $utf8)
`
	if err := os.WriteFile(script, []byte(body), 0o644); err != nil {
		return ""
	}
	defer os.Remove(script)
	start := os.Getenv("USERPROFILE")
	cmd := exec.Command("powershell.exe", "-NoProfile", "-STA", "-ExecutionPolicy", "Bypass", "-File", script)
	cmd.Env = envFor("PICK_OUT="+out, "PICK_START="+start)
	if err := cmd.Run(); err != nil {
		return ""
	}
	b, err := os.ReadFile(out)
	_ = os.Remove(out)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(strings.TrimPrefix(string(b), "\uFEFF"))
}

// claudeProjectURL builds the web UI route for a Windows path.
// The UI treats /projects/H:/dir as drive H:.
func claudeProjectURL(base, dir string) string {
	dir = filepath.ToSlash(filepath.Clean(dir))
	parts := strings.Split(dir, "/")
	for i, p := range parts {
		if !(i == 0 && len(p) == 2 && p[1] == ':') {
			parts[i] = url.PathEscape(p)
		}
	}
	return strings.TrimRight(base, "/") + "/projects/" + strings.Join(parts, "/")
}

// runClaudeUI starts claude-code-webui against claude-wrap.exe and opens a browser.
// Strata is already running. Closing this process stops the UI; the caller stops Strata.
func runClaudeUI(project string) int {
	node := findNodeExe()
	if node == "" {
		die("没有找到 node.exe, 请先安装 Node.js")
	}
	script := findWebUIScript()
	if script == "" {
		die("没有找到 claude-code-webui。请在 claudeui 目录执行: npm install claude-code-webui")
	}
	wrap := findWrapExe()
	if wrap == "" {
		die("没有找到 claude-wrap.exe (应和本程序在同一目录)。请运行 build.ps1 -Only claudewrap")
	}

	port := os.Getenv("CLAUDE_UI_PORT")
	if port == "" {
		port = "3456"
		for i := 0; i < 10 && portBusy(port); i++ {
			port = fmt.Sprint(3457 + i)
		}
	}
	baseURL := "http://127.0.0.1:" + port
	openURL := baseURL
	if project != "" {
		openURL = claudeProjectURL(baseURL, project)
	}

	say("启动网页界面 %s", baseURL)
	say("关闭此窗口即停止网页界面和 Strata, 并释放显存")
	cmd := exec.Command(node, script, "--host", "127.0.0.1", "--port", port, "--claude-path", wrap)
	cmd.Env = envFor("NO_PROXY=127.0.0.1,localhost", "no_proxy=127.0.0.1,localhost")
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
	if err := cmd.Start(); err != nil {
		say("启动网页界面失败: %v", err)
		return 1
	}
	addToJob(cmd)

	done := make(chan struct{})
	go func() {
		client := &http.Client{Timeout: 2 * time.Second}
		for i := 0; i < 60; i++ {
			select {
			case <-done:
				return
			default:
			}
			if r, err := client.Get(baseURL + "/"); err == nil {
				r.Body.Close()
				if r.StatusCode < 500 {
					say("网页界面已就绪, 正在打开浏览器")
					say("地址: %s", openURL)
					openBrowser(openURL)
					return
				}
			}
			time.Sleep(500 * time.Millisecond)
		}
		say("网页界面 30 秒内没有就绪, 请查看上方日志")
	}()

	err := cmd.Wait()
	close(done)
	if ee, ok := err.(*exec.ExitError); ok {
		return ee.ExitCode()
	}
	return 0
}

// runClaude starts the local Claude Code CLI in this console, layered with the Strata
// settings file (--settings), so the user's global ~/.claude/settings.json stays untouched.
func runClaude(args []string) int {
	if _, err := exec.LookPath("claude"); err != nil {
		die("没有找到 claude 命令,请先安装 Claude Code (npm i -g @anthropic-ai/claude-code)")
	}
	cfg, external, err := strataSettingsPath()
	if err != nil {
		die("写配置失败: %v", err)
	}
	if external {
		say("使用外部配置 %s", cfg)
	}

	// Ctrl+C belongs to Claude Code; we only clean up after it exits.
	sig := make(chan os.Signal, 4)
	signal.Notify(sig, os.Interrupt)
	go func() {
		for range sig {
		}
	}()

	say("启动 Claude Code (模型 %s, 上下文 %d, 走 Strata /v1/messages)", strataModel, strataCtx)
	cargs := append([]string{"/c", "claude", "--settings", cfg}, args...)
	cmd := exec.Command("cmd", cargs...)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	cmd.Env = envFor("NO_PROXY=127.0.0.1,localhost", "no_proxy=127.0.0.1,localhost")
	err = cmd.Run()
	if ee, ok := err.(*exec.ExitError); ok {
		return ee.ExitCode()
	} else if err != nil {
		say("运行 claude 失败: %v", err)
		return 1
	}
	return 0
}

const claudeDesktopStrataID = "7c9e8a12-4b3f-4d2e-9f01-a1b2c3d4e5f6"

func findClaudeDesktopApp() string {
	if p := os.Getenv("CLAUDE_DESKTOP_EXE"); p != "" {
		if st, err := os.Stat(p); err == nil && !st.IsDir() {
			abs, _ := filepath.Abs(p)
			return abs
		}
	}
	if la := os.Getenv("LOCALAPPDATA"); la != "" {
		p := filepath.Join(la, "Microsoft", "WindowsApps", "claude-desktop.exe")
		if st, err := os.Stat(p); err == nil && !st.IsDir() {
			return p
		}
	}
	if p, err := exec.LookPath("claude-desktop.exe"); err == nil {
		return p
	}
	return ""
}

func claudeDesktopRunning() bool {
	ps := `$p = Get-CimInstance Win32_Process | Where-Object { $_.Name -eq 'Claude.exe' -and $_.ExecutablePath -like '*WindowsApps*Claude*' }; if ($p) { 'yes' }`
	cmd := exec.Command("powershell", "-NoProfile", "-Command", ps)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	out, err := cmd.Output()
	return err == nil && strings.Contains(string(out), "yes")
}

// applyClaudeDesktopStrataConfig writes a Claude Desktop 3P gateway profile for Strata.
// The official desktop app ignores ANTHROPIC_BASE_URL and ~/.claude/settings.json.
func applyClaudeDesktopStrataConfig() error {
	la := os.Getenv("LOCALAPPDATA")
	if la == "" {
		return fmt.Errorf("LOCALAPPDATA is empty")
	}
	lib := filepath.Join(la, "Claude-3p", "configLibrary")
	if err := os.MkdirAll(lib, 0o755); err != nil {
		return err
	}
	cfg := claudeDesktopConfig
	if self, err := os.Executable(); err == nil {
		ext := filepath.Join(filepath.Dir(self), "claude.desktop.strata.json")
		if b, err := os.ReadFile(ext); err == nil && len(b) > 0 {
			cfg = b
			say("使用外部配置 %s", ext)
		}
	}
	path := filepath.Join(lib, claudeDesktopStrataID+".json")
	if err := os.WriteFile(path, cfg, 0o644); err != nil {
		return err
	}
	metaPath := filepath.Join(lib, "_meta.json")
	meta := map[string]any{
		"appliedId": claudeDesktopStrataID,
		"entries": []map[string]string{
			{"id": "a4055f68-7f46-47c8-a1b8-12b76ec92276", "name": "Ollama huihui-64k"},
			{"id": claudeDesktopStrataID, "name": "Strata Qwen3.8-Flash-Next"},
		},
	}
	if b, err := os.ReadFile(metaPath); err == nil {
		var existing struct {
			Entries []struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			} `json:"entries"`
		}
		if json.Unmarshal(b, &existing) == nil {
			seen := map[string]bool{claudeDesktopStrataID: true, "a4055f68-7f46-47c8-a1b8-12b76ec92276": true}
			entries := []map[string]string{
				{"id": "a4055f68-7f46-47c8-a1b8-12b76ec92276", "name": "Ollama huihui-64k"},
				{"id": claudeDesktopStrataID, "name": "Strata Qwen3.8-Flash-Next"},
			}
			for _, e := range existing.Entries {
				if !seen[e.ID] {
					entries = append(entries, map[string]string{"id": e.ID, "name": e.Name})
					seen[e.ID] = true
				}
			}
			meta["entries"] = entries
		}
	}
	out, err := json.MarshalIndent(meta, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(metaPath, append(out, '\n'), 0o644)
}

func runClaudeDesktopApp() int {
	if err := applyClaudeDesktopStrataConfig(); err != nil {
		die("写入 Claude 桌面第三方网关配置失败: %v", err)
	}
	exe := findClaudeDesktopApp()
	if exe == "" {
		die("没有找到 Claude 桌面客户端 (claude-desktop.exe)。请先安装官方 Claude Desktop")
	}
	say("启动 Claude 桌面客户端 (网关 http://%s, 模型 %s)", strataAddr, strataModel)
	say("关闭桌面窗口后会自动停止 Strata")
	cmd := exec.Command(exe)
	cmd.Env = envFor("NO_PROXY=127.0.0.1,localhost", "no_proxy=127.0.0.1,localhost")
	if err := cmd.Start(); err != nil {
		say("启动 Claude 桌面客户端失败: %v", err)
		return 1
	}
	cmd.Wait()
	time.Sleep(3 * time.Second)
	for claudeDesktopRunning() {
		time.Sleep(3 * time.Second)
	}
	return 0
}

func runDesktop(cfg string) int {
	exe := findDesktop()
	if exe == "" {
		die("没有找到 OpenCode 桌面版 (OpenCode.exe)。请先安装,或设置环境变量 OPENCODE_DESKTOP_EXE")
	}
	if backend == "strata" {
		say("启动 OpenCode 桌面版 (模型 %s, 上下文 %d)", strataModel, strataCtx)
		say("关闭桌面窗口后会自动停止 Strata")
	} else {
		say("启动 OpenCode 桌面版 (模型 %s, 上下文 %d)", modelName, numCtx)
		say("关闭桌面窗口后会自动停止 Ollama")
	}
	cmd := exec.Command(exe)
	cmd.Env = envFor(
		"OPENCODE_CONFIG="+cfg,
		"OPENCODE_DISABLE_AUTOUPDATE=1",
	)
	if err := cmd.Start(); err != nil {
		say("启动桌面版失败: %v", err)
		return 1
	}
	cmd.Wait()
	// the main process may exit early when another instance already runs
	time.Sleep(3 * time.Second)
	for desktopRunning(exe) {
		time.Sleep(3 * time.Second)
	}
	return 0
}

func main() {
	if isClaudeWrap() {
		os.Exit(runClaudeWrap())
	}
	setUTF8Console()
	initJob()
	args := os.Args[1:]
	if len(args) == 0 && defaultArgs != "" {
		args = strings.Split(defaultArgs, ",")
	}
	if mode == "claudedesktop" {
		setTitle("Claude Desktop + Qwen3.8-Flash-Next (Strata) - 关闭此窗口即停止服务")
	} else if mode == "claudeui" {
		setTitle("Claude Code Web + Qwen3.8-Flash-Next (Strata) - 关闭此窗口即停止服务")
	} else if mode == "claude" {
		setTitle("Claude Code + Qwen3.8-Flash-Next (Strata) - 退出 Claude 后自动停止 Strata")
	} else if mode == "openwebui" {
		setTitle("Open WebUI + Qwen - 关闭此窗口即停止服务")
	} else if backend == "strata" {
		setTitle("OpenCode + Qwen3.8-Flash-Next (Strata) - 关闭此窗口即停止服务")
	} else if len(args) > 0 && args[0] == "web" {
		setTitle("opencode-qwen (Web) - 关闭此窗口即停止服务")
	} else {
		setTitle("opencode-qwen")
	}

	var project string
	if mode == "claudeui" {
		say("请选择项目文件夹。取消则打开最近用过的项目列表")
		project = pickProjectFolder()
		if project != "" {
			say("项目文件夹: %s", project)
		} else {
			say("未选择文件夹, 浏览器里可以选最近用过的项目")
		}
	}

	if backend == "strata" {
		if srv := runStrata(); srv != nil && os.Getenv("QWEN_KEEP_SERVER") == "" {
			cleanup = func() {
				say("关闭 Strata 并释放显存")
				killTree(srv)
			}
		}
	} else if os.Getenv("QWEN_SKIP_OLLAMA") == "" {
		exe := findOllama()
		if exe == "" {
			die("没有找到 Ollama,请先安装 https://ollama.com/download")
		}
		models := modelsDir()
		say("Ollama: %s", exe)
		say("模型目录: %s", models)
		freeLegacyServer()
		var server *exec.Cmd
		if up(ollamaAddr) {
			say("复用已运行的服务 %s", ollamaAddr)
		} else {
			say("启动 Ollama 服务 %s (Flash Attention + q8_0 KV 缓存)", ollamaAddr)
			server = startOllama(exe, models)
		}
		if server != nil && os.Getenv("QWEN_KEEP_SERVER") == "" {
			cleanup = func() {
				say("关闭 Ollama 服务并释放显存")
				killTree(server)
				exec.Command("taskkill", "/IM", "llama-server.exe", "/F").Run()
			}
		}
		ensureModel(exe, models)
		preload()
	}

	if mode == "claude" {
		code := runClaude(args)
		if cleanup != nil {
			cleanup()
		}
		os.Exit(code)
	}
	if mode == "claudeui" {
		code := runClaudeUI(project)
		if cleanup != nil {
			cleanup()
		}
		os.Exit(code)
	}
	if mode == "claudedesktop" {
		code := runClaudeDesktopApp()
		if cleanup != nil {
			cleanup()
		}
		os.Exit(code)
	}

	base := filepath.Join(os.Getenv("LOCALAPPDATA"), "opencode-qwen", appVersion)
	bin := filepath.Join(base, "opencode.exe")
	if mode != "desktop" && mode != "openwebui" && mode != "claudedesktop" {
		if err := extractFile(bin, opencodeBin); err != nil {
			die("释放 opencode 失败: %v", err)
		}
	} else if err := os.MkdirAll(base, 0o755); err != nil {
		die("创建目录失败: %v", err)
	}

	// config: opencode.json next to the launcher wins over the embedded one
	cfgName, cfgData := "opencode.json", defaultConfig
	if backend == "strata" {
		cfgName, cfgData = "opencode.strata.json", strataConfig
	}
	cfg := filepath.Join(base, cfgName)
	selfDir := ""
	if self, err := os.Executable(); err == nil {
		selfDir = filepath.Dir(self)
	}
	if selfDir != "" {
		if _, err := os.Stat(filepath.Join(selfDir, cfgName)); err == nil {
			cfg = filepath.Join(selfDir, cfgName)
			say("使用外部配置 %s", cfg)
		}
	}
	if cfg == filepath.Join(base, cfgName) {
		if err := os.WriteFile(cfg, cfgData, 0o644); err != nil {
			die("写配置失败: %v", err)
		}
	}

	// Let opencode own Ctrl+C; we only clean up after it exits.
	sig := make(chan os.Signal, 4)
	signal.Notify(sig, os.Interrupt)
	go func() {
		for range sig {
		}
	}()

	if mode == "desktop" {
		code := runDesktop(cfg)
		if cleanup != nil {
			cleanup()
		}
		os.Exit(code)
	}
	if mode == "openwebui" {
		code := runOpenWebUI()
		if cleanup != nil {
			cleanup()
		}
		os.Exit(code)
	}

	if len(args) > 0 && args[0] == "web" {
		freeStalePort(args, bin)
	}

	if backend == "strata" {
		say("启动 opencode (模型 %s, 上下文 %d)", strataModel, strataCtx)
	} else {
		say("启动 opencode (模型 %s, 上下文 %d)", modelName, numCtx)
	}
	cmd := exec.Command(bin, args...)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	cmd.Env = envFor(
		"OPENCODE_CONFIG="+cfg,
		"OPENCODE_DISABLE_AUTOUPDATE=1",
	)
	err := cmd.Run() // not placed in the job: it may spawn the user's browser
	code := 0
	if ee, ok := err.(*exec.ExitError); ok {
		code = ee.ExitCode()
	} else if err != nil {
		say("运行 opencode 失败: %v", err)
		code = 1
	}
	if cleanup != nil {
		cleanup()
	}
	os.Exit(code)
}
