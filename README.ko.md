<p align="center"><img src="docs/banner.png" alt="open-lowmem: 작은 메모리로 거대한 LLM 실행" width="720"></p>

# open-lowmem: VRAM 24 GB + RAM 32 GB로 1250억 파라미터 LLM을 실측 약 75 tokens/s로 실행 (Windows)

[English](README.md) | [简体中文](README.zh-CN.md) | [繁體中文](README.zh-TW.md) | [日本語](README.ja.md) | [한국어](README.ko.md) | [Español](README.es.md) | [Français](README.fr.md) | [Deutsch](README.de.md) | [Русский](README.ru.md) | [Português](README.pt-BR.md)

**작은 메모리로 거대한 LLM을.** **VRAM 24 GB 그래픽카드와 RAM 32 GB**를 가진 일반 PC에서 **1250억 파라미터 MoE 모델(Qwen3.8-Flash-Next)**을 로컬로 실행하고, **Claude Code, OpenCode, Claude 데스크톱 앱, Open WebUI**에 연결해 주는 Windows용 원클릭 런처입니다. 오프라인, 프라이버시 보호, 무료.

작성자: **廖亦辰 (Liao Yichen)** · 중국 샤먼 · GitHub [@bathfire-dev](https://github.com/bathfire-dev) · 라이선스: MIT

> 설정 시간을 아껴 드렸다면 Star를 눌러 주세요. 작은 GPU를 가진 분들이 이 프로젝트를 찾는 데 도움이 됩니다.

## 이게 무엇인가요?

`open-lowmem`은 Go로 작성한 Windows용 작은 런처 모음입니다. 각 런처는 먼저 로컬 추론 엔진 [Strata](https://github.com/Niko1221/Strata)를 실행합니다. Strata는 GPU VRAM, 시스템 RAM, SSD가 작업을 나눠 맡게 해서 1250억 파라미터 혼합 전문가(MoE) 모델을 일반 게이밍 PC에 올립니다. 그다음 선택한 클라이언트를 실행하고, 창을 닫으면 모든 프로세스를 정리합니다.

키워드: 로컬 LLM, 저 VRAM, 저메모리, 24 GB GPU로 대형 언어 모델 실행, MoE, 양자화(IQ2_XS), Qwen3.8-Flash-Next, Strata, 로컬 모델로 Claude Code 사용, OpenCode, 오프라인 AI 코딩 도우미, RTX 4090, Windows.

## 실측 속도

이 저장소의 개발 PC(RTX 4090D 24 GB, RAM 32 GB, IQ2_XS, Windows, 2026년 10월 7일)에서 측정했습니다. 512 토큰 생성을 8번 측정한 결과 **62~94 tokens/s, 중앙값 약 77**이었습니다. 요청은 한 번에 하나, 프롬프트는 짧습니다. 실제 속도는 컨텍스트 길이, 백그라운드 부하, 모델 크기에 따라 달라집니다.

## 메모리는 얼마나 필요한가요?

수치는 [Strata 문서](https://github.com/Niko1221/Strata)에서 가져온 것이며 이 저장소에서 직접 측정한 값이 아닙니다.

| 크기(같은 모델) | 필요한 RAM + VRAM | 비고 |
| --- | --- | --- |
| Q2_0 | 약 37.6 GB | 가장 빠름 |
| IQ2_XS | 약 39.2 GB | 권장, 이 저장소의 기본 설정 |
| IQ3_XXS / IQ3_S | 약 47 / 55 GB | 품질이 더 좋지만 RAM이 더 필요 |

Strata는 12 GB 이상의 NVIDIA 또는 AMD 그래픽카드, 32 GB 이상의 RAM, 약 80 GB의 여유 디스크가 필요합니다. 공식 측정으로는 RTX 5070(12 GB)과 64 GB RAM에서 IQ2_XS가 약 79 tokens/s였습니다. 24 GB 카드와 32 GB RAM 조합에서는 Q2_0과 IQ2_XS가 Strata의 저RAM 모드로 실행됩니다. 이 저장소는 RTX 4090D 24 GB + 32 GB RAM에서 개발했습니다.

## 런처

| 런처 | 동작 |
| --- | --- |
| `claude-strata.exe` | Strata를 실행한 뒤 설치된 Claude Code를 `--settings`와 함께 실행합니다. `~/.claude`는 건드리지 않습니다. |
| `claude-strata-ui.exe` | 폴더 선택 창 후 Strata와 `claude-code-webui` 브라우저 UI를 실행합니다. |
| `claude-desktop-strata.exe` | Strata를 실행하고 Claude 데스크톱 앱용 서드파티 게이트웨이 설정을 쓴 뒤 앱을 실행합니다. |
| `opencode-qwen.exe`, `opencode-qwen-web.exe` | 설정이 내장된 OpenCode 터미널 및 웹 UI. |
| `opencode-qwen-strata-desktop.exe` | Strata를 실행한 뒤 OpenCode 데스크톱 앱을 실행합니다. |
| `openwebui-qwen.exe` | Ollama와 Open WebUI. |

런처 창을 닫으면 Windows Job Object를 통해 모든 자식 프로세스가 종료되므로 모델이 VRAM을 계속 차지하지 않습니다.

## FAQ

**24 GB GPU에서 1250억 파라미터 모델을 실행할 수 있나요?**
네. 혼합 전문가(MoE) 모델과, 전문가를 VRAM·RAM·SSD에 나눠 두는 엔진이 있으면 가능합니다. Strata가 바로 그 일을 하고, 이 저장소는 그것을 여러분의 도구에 연결할 뿐입니다. IQ2_XS는 RAM과 VRAM을 합쳐 약 39 GB가 필요합니다.

**Claude Code를 로컬 모델과 함께 쓰려면?**
Strata는 Anthropic 호환 엔드포인트(`127.0.0.1:8090`의 `/v1/messages`)를 제공합니다. `claude-strata.exe`는 Strata를 실행하고, 그 주소를 가리키는 설정 오버레이와 함께 Claude Code를 실행합니다. 모델 이름은 `claude-sonnet-4-6` 같은 Claude 모델 이름을 별칭으로 씁니다. Claude Code 자체는 포함되어 있지 않으니 Anthropic에서 설치하세요.

**Claude 데스크톱 앱에서도 되나요?**
검증하지 못했습니다. 설정은 기록되지만, 앱이 일반 `http://127.0.0.1:8090` 게이트웨이를 받아들이는지는 모릅니다. `claude-desktop-strata.exe`는 실험적 기능으로 보세요.

**클라우드 모델과 품질이 같나요?**
아니요. 낮은 비트 양자화는 품질을 내주고 메모리를 얻는 방식입니다. 공개 리더보드에서 Qwen3.8-Flash-Next는 강한 오픈 모델이지만 최상위 폐쇄형 모델은 아니며, 점수 일부는 제공사 자체 보고입니다.

**Linux나 macOS에서도 되나요?**
지금은 안 됩니다. 런처가 Windows API를 사용합니다.

## 중요 안내

- Strata, OpenCode, Claude Code, Claude 데스크톱 앱, 모델 가중치는 **포함되어 있지 않습니다**. 직접 설치하세요.
- 독립적인 개인 프로젝트이며 Anthropic, Alibaba Cloud(Qwen), Strata, OpenCode와 관련이 없습니다.

## 소스에서 빌드

1. Go 1.25+, Node.js, PowerShell 5.1+를 설치합니다.
2. Strata를 `strata\Strata-main`에 두고(또는 `STRATA_DIR` 환경 변수 설정) `run-*.bat`를 생성합니다.
3. `opencode\`(OpenCode 런처)와 `claudeui\`(웹 UI)에서 `npm install`을 실행합니다.
4. `.\build.ps1`을 실행하거나, 하나만 빌드하려면 `.\build.ps1 -Only <이름>`(tui, web, desktop, openwebui, strata, claude, claudewrap, claudeui, claudedesktop)을 실행합니다. 결과물은 `dist\`에 생깁니다.

## 기여

Issue와 Pull Request를 환영합니다. 특히 Claude 데스크톱 경로의 테스트 결과, 다른 GPU의 속도 수치, Linux 포팅이 필요합니다. GPU, RAM, Strata 모델 크기를 함께 적어 주세요. 번역 개선도 환영합니다.

## 참고: WeChat 그룹

중국어권 사용자를 위한 WeChat 그룹 "AI时代"입니다(QR 코드는 10월 14일까지 유효하며, 만료되면 Issue로 알려 주세요).

<img src="docs/wechat-group.png" alt="WeChat 그룹 QR 코드" width="260">

## 크레딧 및 서드파티 라이선스

[Strata](https://github.com/Niko1221/Strata), [OpenCode](https://github.com/anomalyco/opencode)(MIT), [claude-code-webui](https://github.com/sugyan/claude-code-webui)(MIT), Qwen 모델은 각자의 라이선스를 따릅니다.
