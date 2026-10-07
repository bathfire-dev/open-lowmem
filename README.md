# open-lowmem: run a 125B LLM on a 24 GB GPU and 32 GB RAM (Windows)

**Run huge LLMs on tiny memory.** One-click Windows launchers that run a **125-billion-parameter MoE model (Qwen3.8-Flash-Next)** locally on a consumer PC with a **24 GB graphics card and 32 GB of RAM**, then plug it into **Claude Code, OpenCode, the Claude desktop app and Open WebUI**. Offline, private, free.

[English](#english) | [中文](#中文)

Author: **廖亦辰 (Liao Yichen)** · Xiamen, China · GitHub [@bathfire-dev](https://github.com/bathfire-dev) · License: MIT

> If this saves you a weekend of setup, a star helps other people with small GPUs find it.

---

## English

### What is this?

`open-lowmem` is a set of small Go launchers for Windows. Each one starts the local inference engine [Strata](https://github.com/Niko1221/Strata), which makes a 125B-parameter mixture-of-experts model fit into a normal gaming PC by sharing the work between GPU VRAM, system RAM and the SSD. The launcher then starts the client you picked and cleans everything up when you close it.

Keywords: local LLM, low VRAM, low memory, run large language models on 24 GB GPU, MoE, quantization (IQ2_XS), Qwen3.8-Flash-Next, Strata, Claude Code with a local model, OpenCode, offline AI coding assistant, RTX 4090, Windows.

### How much memory does it need?

Numbers come from the [Strata documentation](https://github.com/Niko1221/Strata), not from benchmarks in this repo:

| Model size (same model) | RAM + VRAM needed | Notes |
| --- | --- | --- |
| Q2_0 | about 37.6 GB | fastest |
| IQ2_XS | about 39.2 GB | recommended; this repo's default config |
| IQ3_XXS / IQ3_S | about 47 / 55 GB | better quality, needs more RAM |

Strata needs a 12 GB or larger NVIDIA or AMD card, 32 GB or more RAM and about 80 GB of free disk. Its authors measured 79 tokens/s on an RTX 5070 (12 GB) with 64 GB RAM for IQ2_XS. With a 24 GB card and 32 GB RAM, Q2_0 and IQ2_XS run in Strata's low-RAM mode. This repo is developed on an RTX 4090D 24 GB with 32 GB RAM.

### Launchers

| Launcher | What it does |
| --- | --- |
| `claude-strata.exe` | Starts Strata, then runs your installed Claude Code with `--settings`. Your `~/.claude` is not touched. |
| `claude-strata-ui.exe` | Folder picker, then Strata and the `claude-code-webui` browser UI. |
| `claude-desktop-strata.exe` | Starts Strata, writes a third-party gateway profile for the Claude desktop app, then starts the app. |
| `opencode-qwen.exe`, `opencode-qwen-web.exe` | OpenCode terminal and web UI with the config embedded. |
| `opencode-qwen-strata-desktop.exe` | Starts Strata, then the OpenCode desktop app. |
| `openwebui-qwen.exe` | Ollama plus Open WebUI. |

Closing the launcher window kills every child process through a Windows Job Object, so no model keeps eating your VRAM.

### FAQ

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

### Important notes

- Strata, OpenCode, Claude Code, the Claude desktop app and the model weights are **not included**. Install them yourself.
- Independent personal project, not affiliated with Anthropic, Alibaba Cloud (Qwen), Strata or OpenCode.

### Build from source

1. Install Go 1.25+, Node.js and PowerShell 5.1+.
2. Put Strata in `strata\Strata-main` (or set `STRATA_DIR`) and generate its `run-*.bat`.
3. Run `npm install` in `opencode\` (OpenCode launchers) and in `claudeui\` (web UI).
4. Run `.\build.ps1`, or `.\build.ps1 -Only <name>` with one of: tui, web, desktop, openwebui, strata, claude, claudewrap, claudeui, claudedesktop. Output goes to `dist\`.

### Contributing

Issues and pull requests are welcome, especially: a test result of the Claude desktop path, speed numbers on other GPUs, and a Linux port. Please include your GPU, RAM and the Strata model size.

### Chat group

WeChat group "AI时代" (the QR code expires on Oct 14; open an issue if it has expired and I will refresh it):

<img src="docs/wechat-group.png" alt="WeChat group QR code" width="260">

### Credits and third-party licenses

[Strata](https://github.com/Niko1221/Strata), [OpenCode](https://github.com/anomalyco/opencode) (MIT), [claude-code-webui](https://github.com/sugyan/claude-code-webui) (MIT) and the Qwen model each keep their own licenses.

---

## 中文

### 这是什么？

`open-lowmem`（低内存）是一组 Windows 小启动器：用 [Strata](https://github.com/Niko1221/Strata) 推理引擎，在**24GB 显存 + 32GB 内存**的普通电脑上本地运行 **1250 亿参数的 MoE 大模型 Qwen3.8-Flash-Next**，再一键接入 **Claude Code、OpenCode、Claude 桌面版、Open WebUI**。离线、隐私、免费。

关键词：本地大模型、低显存、低内存、24G 显存跑大模型、显存不够怎么办、MoE 混合专家、量化（IQ2_XS）、Qwen3.8-Flash-Next、Strata、Claude Code 接本地模型、OpenCode、离线 AI 编程助手、RTX 4090、Windows。

### 需要多少内存？

数据来自 [Strata 官方文档](https://github.com/Niko1221/Strata)，不是本仓库自己测的：

| 版本（同一个模型） | 所需内存+显存 | 说明 |
| --- | --- | --- |
| Q2_0 | 约 37.6 GB | 最快 |
| IQ2_XS | 约 39.2 GB | 推荐，本仓库默认配置 |
| IQ3_XXS / IQ3_S | 约 47 / 55 GB | 质量更好，需要更多内存 |

Strata 要求 12GB 以上的 NVIDIA 或 AMD 显卡、32GB 以上内存、约 80GB 空闲硬盘。官方在 RTX 5070（12GB）+ 64GB 内存上实测 IQ2_XS 约 79 tokens/s。24GB 显卡配 32GB 内存时，Q2_0 和 IQ2_XS 可以用 Strata 的低内存模式运行。本仓库在 RTX 4090D 24GB + 32GB 内存的机器上开发。

### 启动器

| 启动器 | 作用 |
| --- | --- |
| `claude-strata.exe` | 先启动 Strata，再用 `--settings` 运行本机的 Claude Code（不改动你的 `~/.claude`） |
| `claude-strata-ui.exe` | 弹出选文件夹窗口，启动 Strata 和 `claude-code-webui` 网页界面 |
| `claude-desktop-strata.exe` | 启动 Strata，写入 Claude 桌面版的第三方网关配置，再启动桌面版 |
| `opencode-qwen.exe` / `opencode-qwen-web.exe` | OpenCode 终端 / 网页版（内嵌 OpenCode 和配置） |
| `opencode-qwen-strata-desktop.exe` | 先启动 Strata，再打开 OpenCode 桌面版 |
| `openwebui-qwen.exe` | Ollama + Open WebUI |

关闭启动器窗口时，会通过 Windows Job Object 结束所有子进程，不会让模型一直占着显存。

### 常见问题

**24G 显卡能跑 1250 亿参数的模型吗？**
可以，前提是 MoE 模型，加上能把专家分到显存、内存和 SSD 的引擎。Strata 做的就是这件事，本仓库只负责把它接到你的工具上。IQ2_XS 大约需要 39GB 的内存加显存。

**Claude Code 怎么接本地模型？**
Strata 提供兼容 Anthropic 的接口（`127.0.0.1:8090` 的 `/v1/messages`）。`claude-strata.exe` 会启动 Strata，再用一份设置叠加文件启动 Claude Code，把请求指向它，模型名用 `claude-sonnet-4-6` 这类别名。Claude Code 本身不在仓库里，请自行从 Anthropic 安装。

**Claude 桌面版能用吗？**
没有验证。配置已经写入，但桌面版是否接受 `http://127.0.0.1:8090` 作为网关地址还不确定，请把 `claude-desktop-strata.exe` 当作实验功能。

**效果和云端模型一样吗？**
不一样。低比特量化是用质量换内存。Qwen3.8-Flash-Next 在公开榜单上是不错的开源模型，但不是顶级闭源模型，而且部分分数来自厂商自报。

**支持 Linux 或 macOS 吗？**
目前不支持，启动器用了 Windows 接口。

### 重要说明

- 本仓库**不包含** Strata、OpenCode、Claude Code、Claude 桌面版和模型文件，需要自行安装。
- 这是个人项目，与 Anthropic、阿里云（Qwen）、Strata、OpenCode 均无关联。

### 从源码构建

1. 安装 Go 1.25+、Node.js、PowerShell 5.1+。
2. 把 Strata 解压到 `strata\Strata-main`（或设置环境变量 `STRATA_DIR`），并按它的文档生成 `run-*.bat`。
3. 需要 OpenCode：在 `opencode\` 里执行 `npm install`；需要网页界面：在 `claudeui\` 里执行 `npm install`。
4. 运行 `.\build.ps1`，或 `.\build.ps1 -Only <名字>`，名字可选：tui, web, desktop, openwebui, strata, claude, claudewrap, claudeui, claudedesktop。产物在 `dist\`。

### 参与贡献

欢迎提 Issue 和 PR，尤其是：Claude 桌面版的测试结果、其他显卡的速度数据、Linux 移植。请附上你的显卡、内存和 Strata 模型版本。

### 交流群

微信扫码加入「AI时代」群聊（二维码 7 天有效，截止 10 月 14 日；过期了请在 Issues 里留言，我会更新）：

<img src="docs/wechat-group.png" alt="微信群「AI时代」二维码" width="260">

### 致谢与第三方许可

- [Strata](https://github.com/Niko1221/Strata)：推理引擎，请遵循其许可证。
- [OpenCode](https://github.com/anomalyco/opencode)：MIT。
- [claude-code-webui](https://github.com/sugyan/claude-code-webui)：MIT。
- Qwen3.8-Flash-Next 模型请遵循其发布方的许可证。
