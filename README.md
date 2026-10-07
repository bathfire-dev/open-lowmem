# open-strata-launcher

[中文](#中文) | [English](#english)

Author: **bathfire-dev** · License: MIT

---

## 中文

在 Windows 上一键启动本地大模型 **Qwen3.8-Flash-Next**（通过 [Strata](https://github.com/Niko1221/Strata) 推理引擎），并接到几种客户端：

| 启动器 | 作用 |
| --- | --- |
| `opencode-qwen.exe` / `opencode-qwen-web.exe` | OpenCode 终端 / 网页版（内嵌 OpenCode 和配置） |
| `opencode-qwen-strata-desktop.exe` | 先启动 Strata，再打开 OpenCode 桌面版 |
| `claude-strata.exe` | 先启动 Strata，再用 `--settings` 运行本机的 Claude Code（不改动你的 `~/.claude`） |
| `claude-strata-ui.exe` | 弹出选文件夹窗口，启动 Strata 和 `claude-code-webui` 网页界面 |
| `claude-desktop-strata.exe` | 启动 Strata，写入 Claude 桌面版的第三方网关配置，再启动桌面版 |
| `openwebui-qwen.exe` | Ollama + Open WebUI |

关闭启动器窗口时，会通过 Windows Job Object 一并结束它拉起的所有子进程。

### 重要说明（请先读）

- **Claude 桌面版这条路径还没有端到端测试过。** 配置写入已经实现，但桌面版是否接受 `http://127.0.0.1:8090` 作为网关地址（有的教程说必须 HTTPS）尚未验证。
- 仅支持 Windows。作者的机器是 RTX 4090D 24GB + 32GB 内存。
- 本仓库**不包含** Strata、OpenCode、Claude Code、Claude 桌面版、模型文件，需要你自己安装，见下。
- 这是个人项目，与 Anthropic、阿里云（Qwen）、Strata、OpenCode 均无关联。

### 使用

1. 安装 Go 1.25+、Node.js、PowerShell 5.1+。
2. 把 Strata 解压到 `strata\Strata-main`（或设置环境变量 `STRATA_DIR`），并按它的文档生成 `run-*.bat`。
3. 需要 OpenCode：在 `opencode\` 里执行 `npm install`；需要网页界面：在 `claudeui\` 里执行 `npm install`。
4. 构建：

```powershell
.\build.ps1                    # 全部
.\build.ps1 -Only claude       # 只构建某一个：tui, web, desktop, openwebui, strata, claude, claudewrap, claudeui, claudedesktop
```

产物在 `dist\`。构建脚本只含 ASCII，以兼容 PowerShell 5.1。

### 目录

- `launcher/` Go 源码和内嵌配置
- `config/` 各客户端的配置模板
- `icon/` 图标和生成脚本
- `build.ps1` 构建脚本

### 交流群

微信扫码加入「AI时代」群聊（二维码 7 天有效，截止 10 月 14 日；过期了请在 Issues 里留言，我会更新）：

<img src="docs/wechat-group.png" alt="微信群「AI时代」二维码" width="260">

### 致谢与第三方许可

- [Strata](https://github.com/Niko1221/Strata)：推理引擎，请遵循其许可证。
- [OpenCode](https://github.com/anomalyco/opencode)：MIT。
- [claude-code-webui](https://github.com/sugyan/claude-code-webui)：MIT。
- Qwen3.8-Flash-Next 模型请遵循其发布方的许可证。

---

## English

One-click Windows launchers that start a local **Qwen3.8-Flash-Next** through the [Strata](https://github.com/Niko1221/Strata) inference engine and attach it to several clients: OpenCode (terminal, web, desktop), Claude Code (via `--settings`, leaving your `~/.claude` untouched), a folder-picker web UI (`claude-code-webui`), the Claude desktop app (third-party gateway profile), and Open WebUI. When you close the launcher, every child process is killed through a Windows Job Object.

### Read this first

- **The Claude desktop path has not been tested end to end.** The profile is written, but it is unverified whether the app accepts a plain `http://127.0.0.1:8090` gateway (some guides say HTTPS is required).
- Windows only. Developed on an RTX 4090D 24 GB with 32 GB RAM.
- Strata, OpenCode, Claude Code, the Claude desktop app and the model weights are **not** included. Install them yourself.
- Independent personal project, not affiliated with Anthropic, Alibaba Cloud (Qwen), Strata or OpenCode.

### Build

1. Install Go 1.25+, Node.js and PowerShell 5.1+.
2. Put Strata in `strata\Strata-main` (or set `STRATA_DIR`) and generate its `run-*.bat`.
3. `npm install` in `opencode\` (OpenCode launchers) and in `claudeui\` (web UI).
4. Run `.\build.ps1`, or `.\build.ps1 -Only <name>` with one of: tui, web, desktop, openwebui, strata, claude, claudewrap, claudeui, claudedesktop. Output goes to `dist\`.

### Chat group

WeChat group "AI时代" (the QR code expires on Oct 14; open an issue if it has expired and I will refresh it):

<img src="docs/wechat-group.png" alt="WeChat group QR code" width="260">

### Credits and third-party licenses

[Strata](https://github.com/Niko1221/Strata), [OpenCode](https://github.com/anomalyco/opencode) (MIT), [claude-code-webui](https://github.com/sugyan/claude-code-webui) (MIT) and the Qwen model each keep their own licenses.
