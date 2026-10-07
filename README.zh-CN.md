<p align="center"><img src="docs/banner.png" alt="open-lowmem：用极少的内存跑超大模型" width="720"></p>

# open-lowmem：24GB 显卡 + 32GB 内存运行 1250 亿参数大模型，实测约 75 tokens/s（Windows）

[English](README.md) | [简体中文](README.zh-CN.md) | [繁體中文](README.zh-TW.md) | [日本語](README.ja.md) | [한국어](README.ko.md) | [Español](README.es.md) | [Français](README.fr.md) | [Deutsch](README.de.md) | [Русский](README.ru.md) | [Português](README.pt-BR.md)

**用极少的内存跑超大模型。** 一组 Windows 一键启动器：在 **24GB 显存 + 32GB 内存**的普通电脑上本地运行 **1250 亿参数的 MoE 大模型 Qwen3.8-Flash-Next**，再一键接入 **Claude Code、OpenCode、Claude 桌面版、Open WebUI**。离线、隐私、免费。

作者：**廖亦辰** · 中国厦门 · GitHub [@bathfire-dev](https://github.com/bathfire-dev) · 许可证：MIT

> 如果它帮你省下了折腾的时间，欢迎点个 Star，让更多显存不大的朋友找到它。

## 这是什么？

`open-lowmem`（低内存）是一组 Windows 小启动器，用 Go 编写。每个启动器都会先启动本地推理引擎 [Strata](https://github.com/Niko1221/Strata)，它把显存、内存和 SSD 配合起来，让 1250 亿参数的混合专家模型塞进普通游戏电脑。然后启动你选的客户端，关闭窗口时自动清理所有进程。

关键词：本地大模型、低显存、低内存、24G 显存跑大模型、显存不够怎么办、MoE 混合专家、量化（IQ2_XS）、Qwen3.8-Flash-Next、Strata、Claude Code 接本地模型、OpenCode、离线 AI 编程助手、RTX 4090、Windows。

## 实测速度

在本仓库的开发机上实测（RTX 4090D 24GB、32GB 内存、IQ2_XS、Windows，2026 年 10 月 7 日）：生成 512 个 token，共测 8 次，结果为 **62–94 tokens/s，中位数约 77**。一次一个请求，提示词很短。实际速度会随上下文长度、后台负载和模型版本变化。

## 需要多少内存？

数据来自 [Strata 官方文档](https://github.com/Niko1221/Strata)，不是本仓库自己测的：

| 版本（同一个模型） | 所需内存+显存 | 说明 |
| --- | --- | --- |
| Q2_0 | 约 37.6 GB | 最快 |
| IQ2_XS | 约 39.2 GB | 推荐，本仓库默认配置 |
| IQ3_XXS / IQ3_S | 约 47 / 55 GB | 质量更好，需要更多内存 |

Strata 要求 12GB 以上的 NVIDIA 或 AMD 显卡、32GB 以上内存、约 80GB 空闲硬盘。官方在 RTX 5070（12GB）+ 64GB 内存上实测 IQ2_XS 约 79 tokens/s。24GB 显卡配 32GB 内存时，Q2_0 和 IQ2_XS 可以用 Strata 的低内存模式运行。本仓库在 RTX 4090D 24GB + 32GB 内存的机器上开发。

## 启动器

| 启动器 | 作用 |
| --- | --- |
| `claude-strata.exe` | 先启动 Strata，再用 `--settings` 运行本机的 Claude Code（不改动你的 `~/.claude`） |
| `claude-strata-ui.exe` | 弹出选文件夹窗口，启动 Strata 和 `claude-code-webui` 网页界面 |
| `claude-desktop-strata.exe` | 启动 Strata，写入 Claude 桌面版的第三方网关配置，再启动桌面版 |
| `opencode-qwen.exe` / `opencode-qwen-web.exe` | OpenCode 终端 / 网页版（内嵌 OpenCode 和配置） |
| `opencode-qwen-strata-desktop.exe` | 先启动 Strata，再打开 OpenCode 桌面版 |
| `openwebui-qwen.exe` | Ollama + Open WebUI |

关闭启动器窗口时，会通过 Windows Job Object 结束所有子进程，不会让模型一直占着显存。

## 常见问题

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

## 重要说明

- 本仓库**不包含** Strata、OpenCode、Claude Code、Claude 桌面版和模型文件，需要自行安装。
- 这是个人项目，与 Anthropic、阿里云（Qwen）、Strata、OpenCode 均无关联。

## 从源码构建

1. 安装 Go 1.25+、Node.js、PowerShell 5.1+。
2. 把 Strata 解压到 `strata\Strata-main`（或设置环境变量 `STRATA_DIR`），并按它的文档生成 `run-*.bat`。
3. 需要 OpenCode：在 `opencode\` 里执行 `npm install`；需要网页界面：在 `claudeui\` 里执行 `npm install`。
4. 运行 `.\build.ps1`，或 `.\build.ps1 -Only <名字>`，名字可选：tui, web, desktop, openwebui, strata, claude, claudewrap, claudeui, claudedesktop。产物在 `dist\`。

## 参与贡献

欢迎提 Issue 和 PR，尤其是：Claude 桌面版的测试结果、其他显卡的速度数据、Linux 移植。请附上你的显卡、内存和 Strata 模型版本。翻译有不准确的地方也欢迎指出。

## 交流群

微信扫码加入「AI时代」群聊（二维码 7 天有效，截止 10 月 14 日；过期了请在 Issues 里留言，我会更新）：

<img src="docs/wechat-group.png" alt="微信群「AI时代」二维码" width="260">

## 致谢与第三方许可

- [Strata](https://github.com/Niko1221/Strata)：推理引擎，请遵循其许可证。
- [OpenCode](https://github.com/anomalyco/opencode)：MIT。
- [claude-code-webui](https://github.com/sugyan/claude-code-webui)：MIT。
- Qwen3.8-Flash-Next 模型请遵循其发布方的许可证。
