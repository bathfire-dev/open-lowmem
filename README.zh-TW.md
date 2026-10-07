<p align="center"><img src="docs/banner.png" alt="open-lowmem：用極少的記憶體跑超大模型" width="720"></p>

# open-lowmem：24GB 顯示卡 + 32GB 記憶體執行 1250 億參數大模型，實測約 75 tokens/s（Windows）

[English](README.md) | [简体中文](README.zh-CN.md) | [繁體中文](README.zh-TW.md) | [日本語](README.ja.md) | [한국어](README.ko.md) | [Español](README.es.md) | [Français](README.fr.md) | [Deutsch](README.de.md) | [Русский](README.ru.md) | [Português](README.pt-BR.md)

**用極少的記憶體跑超大模型。** 一組 Windows 一鍵啟動器：在 **24GB 顯示卡 + 32GB 記憶體**的一般電腦上本機執行 **1250 億參數的 MoE 大模型 Qwen3.8-Flash-Next**，再一鍵接上 **Claude Code、OpenCode、Claude 桌面版、Open WebUI**。離線、隱私、免費。

作者：**廖亦辰** · 中國廈門 · GitHub [@bathfire-dev](https://github.com/bathfire-dev) · 授權：MIT

> 如果它幫你省下了折騰的時間，歡迎給個 Star，讓更多顯示卡不大的朋友找到它。

## 這是什麼？

`open-lowmem`（低記憶體）是一組用 Go 撰寫的 Windows 小型啟動器。每個啟動器會先啟動本機推論引擎 [Strata](https://github.com/Niko1221/Strata)，它讓顯示卡的 VRAM、系統記憶體和 SSD 分工合作，使 1250 億參數的混合專家模型能放進一般電競電腦。接著啟動你選的用戶端，關閉視窗時會自動清理所有程序。

關鍵字：本機大模型、低顯存、低記憶體、24G 顯存跑大模型、MoE 混合專家、量化（IQ2_XS）、Qwen3.8-Flash-Next、Strata、Claude Code 接本機模型、OpenCode、離線 AI 程式助手、RTX 4090、Windows。

## 實測速度

在本專案的開發機上實測（RTX 4090D 24GB、32GB 記憶體、IQ2_XS、Windows，2026 年 10 月 7 日）：生成 512 個 token，共測 8 次，結果為 **62–94 tokens/s，中位數約 77**。一次一個請求，提示詞很短。實際速度會隨上下文長度、背景負載與模型版本而變化。

## 需要多少記憶體？

數據來自 [Strata 官方文件](https://github.com/Niko1221/Strata)，不是本專案自己測的：

| 版本（同一個模型） | 所需記憶體 + 顯存 | 說明 |
| --- | --- | --- |
| Q2_0 | 約 37.6 GB | 最快 |
| IQ2_XS | 約 39.2 GB | 建議，本專案預設設定 |
| IQ3_XXS / IQ3_S | 約 47 / 55 GB | 品質更好，需要更多記憶體 |

Strata 需要 12GB 以上的 NVIDIA 或 AMD 顯示卡、32GB 以上記憶體、約 80GB 可用硬碟空間。官方在 RTX 5070（12GB）+ 64GB 記憶體上實測 IQ2_XS 約 79 tokens/s。24GB 顯示卡搭配 32GB 記憶體時，Q2_0 與 IQ2_XS 可用 Strata 的低記憶體模式執行。本專案在 RTX 4090D 24GB + 32GB 記憶體的機器上開發。

## 啟動器

| 啟動器 | 作用 |
| --- | --- |
| `claude-strata.exe` | 先啟動 Strata，再用 `--settings` 執行本機的 Claude Code（不會改動你的 `~/.claude`） |
| `claude-strata-ui.exe` | 跳出選擇資料夾視窗，啟動 Strata 與 `claude-code-webui` 網頁介面 |
| `claude-desktop-strata.exe` | 啟動 Strata，寫入 Claude 桌面版的第三方閘道設定，再啟動桌面版 |
| `opencode-qwen.exe` / `opencode-qwen-web.exe` | OpenCode 終端機 / 網頁版（內嵌 OpenCode 與設定） |
| `opencode-qwen-strata-desktop.exe` | 先啟動 Strata，再開啟 OpenCode 桌面版 |
| `openwebui-qwen.exe` | Ollama + Open WebUI |

關閉啟動器視窗時，會透過 Windows Job Object 結束所有子程序，不會讓模型一直佔著顯存。

## 常見問題

**24G 顯示卡能跑 1250 億參數的模型嗎？**
可以，前提是 MoE 模型，加上能把專家分配到顯存、記憶體和 SSD 的引擎。Strata 做的就是這件事，本專案只負責把它接到你的工具上。IQ2_XS 大約需要 39GB 的記憶體加顯存。

**Claude Code 怎麼接本機模型？**
Strata 提供相容 Anthropic 的介面（`127.0.0.1:8090` 的 `/v1/messages`）。`claude-strata.exe` 會啟動 Strata，再用一份設定疊加檔啟動 Claude Code，把請求指向它，模型名稱使用 `claude-sonnet-4-6` 這類別名。Claude Code 本身不在專案裡，請自行從 Anthropic 安裝。

**Claude 桌面版能用嗎？**
尚未驗證。設定已寫入，但桌面版是否接受 `http://127.0.0.1:8090` 作為閘道位址還不確定，請把 `claude-desktop-strata.exe` 當作實驗功能。

**效果和雲端模型一樣嗎？**
不一樣。低位元量化是用品質換記憶體。Qwen3.8-Flash-Next 在公開排行榜上是不錯的開源模型，但不是頂級閉源模型，而且部分分數來自廠商自報。

**支援 Linux 或 macOS 嗎？**
目前不支援，啟動器使用了 Windows 介面。

## 重要說明

- 本專案**不包含** Strata、OpenCode、Claude Code、Claude 桌面版與模型檔案，需要自行安裝。
- 這是個人專案，與 Anthropic、阿里雲（Qwen）、Strata、OpenCode 均無關聯。

## 從原始碼建置

1. 安裝 Go 1.25+、Node.js、PowerShell 5.1+。
2. 把 Strata 解壓縮到 `strata\Strata-main`（或設定環境變數 `STRATA_DIR`），並依它的文件產生 `run-*.bat`。
3. 需要 OpenCode：在 `opencode\` 執行 `npm install`；需要網頁介面：在 `claudeui\` 執行 `npm install`。
4. 執行 `.\build.ps1`，或 `.\build.ps1 -Only <名稱>`，名稱可選：tui, web, desktop, openwebui, strata, claude, claudewrap, claudeui, claudedesktop。輸出在 `dist\`。

## 參與貢獻

歡迎提交 Issue 與 PR，尤其是：Claude 桌面版的測試結果、其他顯示卡的速度數據、Linux 移植。請附上你的顯示卡、記憶體與 Strata 模型版本。翻譯有不準確之處也歡迎指正。

## 交流群

微信掃碼加入「AI时代」群聊（QR code 7 天有效，截止 10 月 14 日；過期請在 Issues 留言，我會更新）：

<img src="docs/wechat-group.png" alt="微信群「AI时代」QR code" width="260">

## 致謝與第三方授權

- [Strata](https://github.com/Niko1221/Strata)：推論引擎，請遵循其授權。
- [OpenCode](https://github.com/anomalyco/opencode)：MIT。
- [claude-code-webui](https://github.com/sugyan/claude-code-webui)：MIT。
- Qwen3.8-Flash-Next 模型請遵循其發布者的授權。
