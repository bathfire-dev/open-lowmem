# Build the launchers into dist\. All are single files with the config embedded. ASCII only on purpose (PowerShell 5.1).
#   build.ps1                 build everything
#   build.ps1 -Only strata    build only dist\opencode-qwen-strata-desktop.exe (others may be running/locked)
# Names for -Only: tui, web, desktop, openwebui, strata, claude, claudewrap, claudeui, claudedesktop
param([string]$Only = '')
$ErrorActionPreference = 'Stop'
$root   = $PSScriptRoot
$assets = Join-Path $root 'launcher\assets'
New-Item -ItemType Directory -Force $assets, (Join-Path $root 'dist') | Out-Null
function Want([string]$n) { return ($Only -eq '' -or $Only -eq $n) }

$src = Join-Path $root 'opencode\node_modules\opencode-windows-x64\bin\opencode.exe'
$needCli = (Want 'tui') -or (Want 'web')
if ($needCli -and -not (Test-Path $src)) { throw "Missing $src - run 'npm install opencode-ai' in the opencode folder first" }

if ($needCli) { Copy-Item $src (Join-Path $assets 'opencode.exe') -Force }
Copy-Item (Join-Path $root 'config\opencode.json') (Join-Path $assets 'opencode.json') -Force
Copy-Item (Join-Path $root 'config\opencode.strata.json') (Join-Path $assets 'opencode.strata.json') -Force
Copy-Item (Join-Path $root 'config\claude.strata.settings.json') (Join-Path $assets 'claude.strata.settings.json') -Force
Copy-Item (Join-Path $root 'config\claude.desktop.strata.json') (Join-Path $assets 'claude.desktop.strata.json') -Force

Push-Location (Join-Path $root 'launcher')
try {
    if (-not (Test-Path go.mod)) { go mod init opencode-qwen 2>$null | Out-Null }
    $env:GOOS = 'windows'; $env:GOARCH = 'amd64'; $env:CGO_ENABLED = '0'

    if (Want 'tui') {
        # Terminal version: no arguments -> TUI
        go build -trimpath -ldflags '-s -w' -o (Join-Path $root 'dist\opencode-qwen.exe') .
        if ($LASTEXITCODE -ne 0) { throw 'go build (terminal) failed' }
    }
    if (Want 'web') {
        # Web version: no arguments -> "opencode web --port 4096" (opens the browser)
        go build -trimpath -ldflags '-s -w -X main.defaultArgs=web,--port,4096' -o (Join-Path $root 'dist\opencode-qwen-web.exe') .
        if ($LASTEXITCODE -ne 0) { throw 'go build (web) failed' }
    }
    if (Want 'desktop') {
        # Desktop version: starts Ollama, then the installed OpenCode desktop app, cleans up on exit.
        go build -tags desktoponly -trimpath -ldflags '-s -w -X main.mode=desktop' -o (Join-Path $root 'dist\opencode-qwen-desktop.exe') .
        if ($LASTEXITCODE -ne 0) { throw 'go build (desktop) failed' }
    }
    if (Want 'openwebui') {
        # Open WebUI version: starts Ollama, then openwebui\venv\Scripts\open-webui.exe, opens the browser.
        go build -tags desktoponly -trimpath -ldflags '-s -w -X main.mode=openwebui' -o (Join-Path $root 'dist\openwebui-qwen.exe') .
        if ($LASTEXITCODE -ne 0) { throw 'go build (openwebui) failed' }
    }
    if (Want 'strata') {
        # Strata version: unloads Ollama, starts strata\Strata-main\run-*.bat (Qwen3.8-Flash-Next), then the desktop app.
        go build -tags desktoponly -trimpath -ldflags '-s -w -X main.mode=desktop -X main.backend=strata' -o (Join-Path $root 'dist\opencode-qwen-strata-desktop.exe') .
        if ($LASTEXITCODE -ne 0) { throw 'go build (strata) failed' }
    }
    if (Want 'claude') {
        # Claude Code version: starts Strata, then runs the local "claude" CLI with --settings (Strata endpoint).
        go build -tags desktoponly -trimpath -ldflags '-s -w -X main.mode=claude -X main.backend=strata' -o (Join-Path $root 'dist\claude-strata.exe') .
        if ($LASTEXITCODE -ne 0) { throw 'go build (claude) failed' }
    }
    if (Want 'claudewrap') {
        # Helper spawned by the web UI. Injects --settings and forwards to the real claude.exe.
        # Must not start Strata and must not write anything of its own to stdout.
        go build -tags desktoponly -trimpath -ldflags '-s -w -X main.mode=claudewrap' -o (Join-Path $root 'dist\claude-wrap.exe') .
        if ($LASTEXITCODE -ne 0) { throw 'go build (claudewrap) failed' }
    }
    if (Want 'claudeui') {
        # Browser UI: folder picker, then Strata, then claude-code-webui (needs claude-wrap.exe beside it).
        go build -tags desktoponly -trimpath -ldflags '-s -w -X main.mode=claudeui -X main.backend=strata' -o (Join-Path $root 'dist\claude-strata-ui.exe') .
        if ($LASTEXITCODE -ne 0) { throw 'go build (claudeui) failed' }
    }
    if (Want 'claudedesktop') {
        # Official Claude Desktop (MSIX): starts Strata, points the 3P gateway at it, then launches the app.
        go build -tags desktoponly -trimpath -ldflags '-s -w -X main.mode=claudedesktop -X main.backend=strata' -o (Join-Path $root 'dist\claude-desktop-strata.exe') .
        if ($LASTEXITCODE -ne 0) { throw 'go build (claudedesktop) failed' }
    }
}
finally {
    Pop-Location
    Remove-Item (Join-Path $assets 'opencode.exe') -Force -ErrorAction SilentlyContinue
}

Get-ChildItem (Join-Path $root 'dist') -Filter *.exe | Format-Table Name, @{n='MB';e={[math]::Round($_.Length/1MB,1)}}
