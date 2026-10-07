<p align="center"><img src="docs/banner.png" alt="open-lowmem: run huge LLMs on tiny memory" width="720"></p>

# open-lowmem: run a 125B LLM on a 24 GB GPU and 32 GB RAM (Windows)

[English](README.md) | [简体中文](README.zh-CN.md) | [繁體中文](README.zh-TW.md) | [日本語](README.ja.md) | [한국어](README.ko.md) | [Español](README.es.md) | [Français](README.fr.md) | [Deutsch](README.de.md) | [Русский](README.ru.md) | [Português](README.pt-BR.md)

**Run huge LLMs on tiny memory.** One-click Windows launchers that run a **125-billion-parameter MoE model (Qwen3.8-Flash-Next)** locally on a consumer PC with a **24 GB graphics card and 32 GB of RAM**, then plug it into **Claude Code, OpenCode, the Claude desktop app and Open WebUI**. Offline, private, free.

Author: **廖亦辰 (Liao Yichen)** · Xiamen, China · GitHub [@bathfire-dev](https://github.com/bathfire-dev) · License: MIT

> If this saves you a weekend of setup, a star helps other people with small GPUs find it.

## What is this?

`open-lowmem` is a set of small Go launchers for Windows. Each one starts the local inference engine [Strata](https://github.com/Niko1221/Strata), which makes a 125B-parameter mixture-of-experts model fit into a normal gaming PC by sharing the work between GPU VRAM, system RAM and the SSD. The launcher then starts the client you picked and cleans everything up when you close it.

Keywords: local LLM, low VRAM, low memory, run large language models on 24 GB GPU, MoE, quantization (IQ2_XS), Qwen3.8-Flash-Next, Strata, Claude Code with a local model, OpenCode, offline AI coding assistant, RTX 4090, Windows.

## How much memory does it need?

Numbers come from the [Strata documentation](https://github.com/Niko1221/Strata), not from benchmarks in this repo:

| Model size (same model) | RAM + VRAM needed | Notes |
| --- | --- | --- |
| Q2_0 | about 37.6 GB | fastest |
| IQ2_XS | about 39.2 GB | recommended; this repo's default config |
| IQ3_XXS / IQ3_S | about 47 / 55 GB | better quality, needs more RAM |

Strata needs a 12 GB or larger NVIDIA or AMD card, 32 GB or more RAM and about 80 GB of free disk. Its authors measured 79 tokens/s on an RTX 5070 (12 GB) with 64 GB RAM for IQ2_XS. With a 24 GB card and 32 GB RAM, Q2_0 and IQ2_XS run in Strata's low-RAM mode. This repo is developed on an RTX 4090D 24 GB with 32 GB RAM.

## Launchers

| Launcher | What it does |
| --- | --- |
| `claude-strata.exe` | Starts Strata, then runs your installed Claude Code with `--settings`. Your `~/.claude` is not touched. |
| `claude-strata-ui.exe` | Folder picker, then Strata and the `claude-code-webui` browser UI. |
| `claude-desktop-strata.exe` | Starts Strata, writes a third-party gateway profile for the Claude desktop app, then starts the app. |
| `opencode-qwen.exe`, `opencode-qwen-web.exe` | OpenCode terminal and web UI with the config embedded. |
| `opencode-qwen-strata-desktop.exe` | Starts Strata, then the OpenCode desktop app. |
| `openwebui-qwen.exe` | Ollama plus Open WebUI. |

Closing the launcher window kills every child process through a Windows Job Object, so no model keeps eating your VRAM.

## FAQ

**Can I run a 125B-parameter model on a 24 GB GPU?**
Yes, with a mixture-of-experts model and an engine that splits the experts across VRAM, RAM and SSD. That is what Strata does; this repo only wires it to your tools. Expect to need about 39 GB of combined RAM and VRAM for IQ2_XS.

**How do I use Claude Code with a local model?**
Strata exposes an Anthropic-compatible endpoint (`/v1/messages` on `127.0.0.1:8090`). `claude-strata.exe` starts Strata and launches Claude Code with a settings overlay that points at it, using a Claude model name such as `claude-sonnet-4-6` as an alias. Claude Code itself is not included; install it from Anthropic.

**Does the Claude desktop app work with it?**
Not verified. The profile is written, but it is unknown whether the app accepts a plain `http://127.0.0.1:8090` gateway. Treat `claude-desktop-strata.exe` as experimental.

**Is it quality-equal to cloud models?**
No. Low-bit quantization trades quality for memory. On public leaderboards Qwen3.8-Flash-Next is a strong open model, but not a frontier closed model, and these scores are partly vendor-reported.

**Does it work on Linux or macOS?**
Not now. The launchers use Windows APIs.

## Important notes

- Strata, OpenCode, Claude Code, the Claude desktop app and the model weights are **not included**. Install them yourself.
- Independent personal project, not affiliated with Anthropic, Alibaba Cloud (Qwen), Strata or OpenCode.

## Build from source

1. Install Go 1.25+, Node.js and PowerShell 5.1+.
2. Put Strata in `strata\Strata-main` (or set `STRATA_DIR`) and generate its `run-*.bat`.
3. Run `npm install` in `opencode\` (OpenCode launchers) and in `claudeui\` (web UI).
4. Run `.\build.ps1`, or `.\build.ps1 -Only <name>` with one of: tui, web, desktop, openwebui, strata, claude, claudewrap, claudeui, claudedesktop. Output goes to `dist\`.

## Contributing

Issues and pull requests are welcome, especially: a test result of the Claude desktop path, speed numbers on other GPUs, and a Linux port. Please include your GPU, RAM and the Strata model size. Translations can be improved the same way.

## Chat group

WeChat group "AI时代" (the QR code expires on Oct 14; open an issue if it has expired and I will refresh it):

<img src="docs/wechat-group.png" alt="WeChat group QR code" width="260">

## Credits and third-party licenses

[Strata](https://github.com/Niko1221/Strata), [OpenCode](https://github.com/anomalyco/opencode) (MIT), [claude-code-webui](https://github.com/sugyan/claude-code-webui) (MIT) and the Qwen model each keep their own licenses.
