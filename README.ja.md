<p align="center"><img src="docs/banner.png" alt="open-lowmem: 小さなメモリで巨大な LLM を動かす" width="720"></p>

# open-lowmem: VRAM 24 GB + RAM 32 GB で 1250 億パラメータの LLM を動かす（Windows）

[English](README.md) | [简体中文](README.zh-CN.md) | [繁體中文](README.zh-TW.md) | [日本語](README.ja.md) | [한국어](README.ko.md) | [Español](README.es.md) | [Français](README.fr.md) | [Deutsch](README.de.md) | [Русский](README.ru.md) | [Português](README.pt-BR.md)

**小さなメモリで巨大な LLM を。** **VRAM 24 GB のグラフィックスカードと 32 GB の RAM** を積んだ一般的な PC で、**1250 億パラメータの MoE モデル（Qwen3.8-Flash-Next）** をローカル実行し、**Claude Code、OpenCode、Claude デスクトップアプリ、Open WebUI** につなぐ Windows 用ワンクリックランチャーです。オフライン、プライバシー重視、無料。

作者: **廖亦辰 (Liao Yichen)** · 中国・厦門 · GitHub [@bathfire-dev](https://github.com/bathfire-dev) · ライセンス: MIT

> セットアップの手間を減らせたなら、Star をいただけると、小さな GPU の人が見つけやすくなります。

## これは何？

`open-lowmem` は Go で書かれた Windows 用の小さなランチャー集です。どれも最初にローカル推論エンジン [Strata](https://github.com/Niko1221/Strata) を起動します。Strata は GPU の VRAM、システム RAM、SSD で処理を分担し、1250 億パラメータの Mixture-of-Experts モデルを一般的なゲーミング PC に収めます。そのあと選んだクライアントを起動し、ウィンドウを閉じるとすべてのプロセスを片付けます。

キーワード: ローカル LLM、低 VRAM、低メモリ、24 GB GPU で大規模言語モデル、MoE、量子化（IQ2_XS）、Qwen3.8-Flash-Next、Strata、Claude Code をローカルモデルで使う、OpenCode、オフライン AI コーディング、RTX 4090、Windows。

## 必要なメモリは？

数値は [Strata のドキュメント](https://github.com/Niko1221/Strata) によるもので、このリポジトリで計測したものではありません。

| サイズ（同じモデル） | 必要な RAM + VRAM | 備考 |
| --- | --- | --- |
| Q2_0 | 約 37.6 GB | 最速 |
| IQ2_XS | 約 39.2 GB | 推奨。このリポジトリの既定設定 |
| IQ3_XXS / IQ3_S | 約 47 / 55 GB | 高品質、RAM がより必要 |

Strata は 12 GB 以上の NVIDIA / AMD カード、32 GB 以上の RAM、約 80 GB の空きディスクを必要とします。公式の計測では、RTX 5070（12 GB）と 64 GB RAM で IQ2_XS が約 79 tokens/s でした。24 GB カードと 32 GB RAM の組み合わせでは、Q2_0 と IQ2_XS が Strata の低 RAM モードで動きます。このリポジトリは RTX 4090D 24 GB + 32 GB RAM で開発しています。

## ランチャー

| ランチャー | 内容 |
| --- | --- |
| `claude-strata.exe` | Strata を起動し、インストール済みの Claude Code を `--settings` 付きで実行。`~/.claude` は変更しません。 |
| `claude-strata-ui.exe` | フォルダ選択のあと、Strata と `claude-code-webui` のブラウザ UI を起動。 |
| `claude-desktop-strata.exe` | Strata を起動し、Claude デスクトップアプリ用のサードパーティゲートウェイ設定を書き込んでアプリを起動。 |
| `opencode-qwen.exe`, `opencode-qwen-web.exe` | 設定を内蔵した OpenCode のターミナル版と Web 版。 |
| `opencode-qwen-strata-desktop.exe` | Strata を起動し、OpenCode デスクトップアプリを起動。 |
| `openwebui-qwen.exe` | Ollama と Open WebUI。 |

ランチャーのウィンドウを閉じると、Windows の Job Object ですべての子プロセスが終了し、モデルが VRAM を占有し続けることはありません。

## FAQ

**24 GB の GPU で 1250 億パラメータのモデルを動かせますか？**
はい。Mixture-of-Experts モデルと、エキスパートを VRAM・RAM・SSD に振り分けるエンジンがあれば可能です。それを行うのが Strata で、このリポジトリはそれをツールにつなぐだけです。IQ2_XS では RAM と VRAM の合計で約 39 GB が目安です。

**Claude Code をローカルモデルで使うには？**
Strata は Anthropic 互換のエンドポイント（`127.0.0.1:8090` の `/v1/messages`）を提供します。`claude-strata.exe` は Strata を起動し、そこを指す設定オーバーレイ付きで Claude Code を起動します。モデル名には `claude-sonnet-4-6` のような Claude のモデル名をエイリアスとして使います。Claude Code 自体は含まれません。Anthropic から入手してください。

**Claude デスクトップアプリで動きますか？**
未検証です。設定は書き込まれますが、アプリが `http://127.0.0.1:8090` というプレーンなゲートウェイを受け付けるかは不明です。`claude-desktop-strata.exe` は実験的なものとして扱ってください。

**クラウドモデルと同じ品質ですか？**
いいえ。低ビット量子化は品質とメモリを引き換えにします。公開リーダーボード上で Qwen3.8-Flash-Next は有力なオープンモデルですが、最先端のクローズドモデルではなく、スコアの一部はベンダー自己申告です。

**Linux や macOS で動きますか？**
今は動きません。ランチャーは Windows API を使っています。

## 重要な注意

- Strata、OpenCode、Claude Code、Claude デスクトップアプリ、モデルの重みは**含まれていません**。各自でインストールしてください。
- 個人の独立プロジェクトで、Anthropic、Alibaba Cloud（Qwen）、Strata、OpenCode とは無関係です。

## ソースからのビルド

1. Go 1.25 以上、Node.js、PowerShell 5.1 以上をインストールします。
2. Strata を `strata\Strata-main` に置き（または環境変数 `STRATA_DIR` を設定）、`run-*.bat` を生成します。
3. `opencode\`（OpenCode 用）と `claudeui\`（Web UI 用）で `npm install` を実行します。
4. `.\build.ps1` を実行します。1 つだけビルドする場合は `.\build.ps1 -Only <名前>`（tui, web, desktop, openwebui, strata, claude, claudewrap, claudeui, claudedesktop）。出力は `dist\` です。

## コントリビュート

Issue と Pull Request を歓迎します。特に、Claude デスクトップ経路のテスト結果、他の GPU での速度、Linux 対応が助かります。GPU、RAM、Strata のモデルサイズを添えてください。翻訳の改善も歓迎です。

## 参考: WeChat グループ

中国語圏向けの WeChat グループ「AI时代」です（QR コードは 10 月 14 日まで有効。期限切れの場合は Issue でお知らせください）。

<img src="docs/wechat-group.png" alt="WeChat グループの QR コード" width="260">

## クレジットとサードパーティのライセンス

[Strata](https://github.com/Niko1221/Strata)、[OpenCode](https://github.com/anomalyco/opencode)（MIT）、[claude-code-webui](https://github.com/sugyan/claude-code-webui)（MIT）、Qwen モデルは、それぞれのライセンスに従います。
