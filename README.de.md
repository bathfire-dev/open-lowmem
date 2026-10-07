<p align="center"><img src="docs/banner.png" alt="open-lowmem: riesige LLMs mit winzigem Speicher ausführen" width="720"></p>

# open-lowmem: ein 125-Mrd.-LLM mit ~75 Tokens/s auf einer 24-GB-GPU mit 32 GB RAM ausführen (Windows)

[English](README.md) | [简体中文](README.zh-CN.md) | [繁體中文](README.zh-TW.md) | [日本語](README.ja.md) | [한국어](README.ko.md) | [Español](README.es.md) | [Français](README.fr.md) | [Deutsch](README.de.md) | [Русский](README.ru.md) | [Português](README.pt-BR.md)

**Riesige LLMs mit winzigem Speicher.** Windows-Starter mit einem Klick, die ein **MoE-Modell mit 125 Milliarden Parametern (Qwen3.8-Flash-Next)** lokal auf einem normalen PC mit **24-GB-Grafikkarte und 32 GB RAM** ausführen und mit **Claude Code, OpenCode, der Claude-Desktop-App und Open WebUI** verbinden. Offline, privat, kostenlos.

Autor: **廖亦辰 (Liao Yichen)** · Xiamen, China · GitHub [@bathfire-dev](https://github.com/bathfire-dev) · Lizenz: MIT

> Wenn dir das ein Wochenende Einrichtung erspart, hilft ein Stern anderen mit kleinen GPUs, das Projekt zu finden.

## Was ist das?

`open-lowmem` ist eine Sammlung kleiner Go-Starter für Windows. Jeder startet zuerst die lokale Inferenz-Engine [Strata](https://github.com/Niko1221/Strata). Sie verteilt die Arbeit auf GPU-VRAM, Arbeitsspeicher und SSD, sodass ein Mixture-of-Experts-Modell mit 125 Mrd. Parametern in einen normalen Gaming-PC passt. Danach startet der Starter den gewählten Client und räumt beim Schließen des Fensters alles auf.

Schlüsselwörter: lokales LLM, wenig VRAM, geringer Speicherbedarf, große Sprachmodelle auf einer 24-GB-GPU, MoE, Quantisierung (IQ2_XS), Qwen3.8-Flash-Next, Strata, Claude Code mit lokalem Modell, OpenCode, Offline-KI-Coding-Assistent, RTX 4090, Windows.

## Gemessene Geschwindigkeit

Gemessen auf dem Entwicklungsrechner dieses Repositorys (RTX 4090D mit 24 GB, 32 GB RAM, IQ2_XS, Windows, 7. Oktober 2026): 8 Läufe mit je 512 erzeugten Tokens ergaben **62 bis 94 Tokens/s, Median etwa 77**. Eine Anfrage nach der anderen, kurzer Prompt. Deine Geschwindigkeit hängt von Kontextlänge, Hintergrundlast und Modellgröße ab.

## Wie viel Speicher wird benötigt?

Die Zahlen stammen aus der [Strata-Dokumentation](https://github.com/Niko1221/Strata), nicht aus Messungen in diesem Repository:

| Größe (gleiches Modell) | Benötigter RAM + VRAM | Hinweise |
| --- | --- | --- |
| Q2_0 | etwa 37,6 GB | am schnellsten |
| IQ2_XS | etwa 39,2 GB | empfohlen; Standardkonfiguration dieses Repositorys |
| IQ3_XXS / IQ3_S | etwa 47 / 55 GB | bessere Qualität, braucht mehr RAM |

Strata benötigt eine NVIDIA- oder AMD-Karte mit mindestens 12 GB, mindestens 32 GB RAM und etwa 80 GB freien Speicherplatz. Die Autoren haben für IQ2_XS 79 Tokens/s auf einer RTX 5070 (12 GB) mit 64 GB RAM gemessen. Mit einer 24-GB-Karte und 32 GB RAM laufen Q2_0 und IQ2_XS in Stratas RAM-sparendem Modus. Dieses Repository wird auf einer RTX 4090D mit 24 GB und 32 GB RAM entwickelt.

## Starter

| Starter | Was er tut |
| --- | --- |
| `claude-strata.exe` | Startet Strata und führt dein installiertes Claude Code mit `--settings` aus. Dein `~/.claude` bleibt unberührt. |
| `claude-strata-ui.exe` | Ordnerauswahl, dann Strata und die Browser-Oberfläche `claude-code-webui`. |
| `claude-desktop-strata.exe` | Startet Strata, schreibt ein Drittanbieter-Gateway-Profil für die Claude-Desktop-App und startet sie. |
| `opencode-qwen.exe`, `opencode-qwen-web.exe` | OpenCode im Terminal und im Web, Konfiguration eingebettet. |
| `opencode-qwen-strata-desktop.exe` | Startet Strata, dann die OpenCode-Desktop-App. |
| `openwebui-qwen.exe` | Ollama plus Open WebUI. |

Beim Schließen des Starter-Fensters werden alle Kindprozesse über ein Windows Job Object beendet, sodass kein Modell weiter deinen VRAM belegt.

## FAQ

**Kann ich ein Modell mit 125 Mrd. Parametern auf einer 24-GB-GPU ausführen?**
Ja, mit einem Mixture-of-Experts-Modell und einer Engine, die die Experten auf VRAM, RAM und SSD aufteilt. Das macht Strata; dieses Repository verbindet es nur mit deinen Werkzeugen. Rechne bei IQ2_XS mit etwa 39 GB RAM und VRAM zusammen.

**Wie nutze ich Claude Code mit einem lokalen Modell?**
Strata stellt einen Anthropic-kompatiblen Endpunkt bereit (`/v1/messages` auf `127.0.0.1:8090`). `claude-strata.exe` startet Strata und startet Claude Code mit einer Einstellungs-Überlagerung, die darauf zeigt. Als Alias dient ein Claude-Modellname wie `claude-sonnet-4-6`. Claude Code selbst ist nicht enthalten; installiere es von Anthropic.

**Funktioniert die Claude-Desktop-App damit?**
Nicht geprüft. Das Profil wird geschrieben, aber es ist unbekannt, ob die App ein einfaches `http://127.0.0.1:8090`-Gateway akzeptiert. Behandle `claude-desktop-strata.exe` als experimentell.

**Ist die Qualität gleich wie bei Cloud-Modellen?**
Nein. Quantisierung mit wenigen Bits tauscht Qualität gegen Speicher. In öffentlichen Ranglisten ist Qwen3.8-Flash-Next ein starkes offenes Modell, aber kein geschlossenes Spitzenmodell, und die Werte stammen teils vom Hersteller selbst.

**Läuft es unter Linux oder macOS?**
Derzeit nicht. Die Starter nutzen Windows-APIs.

## Wichtige Hinweise

- Strata, OpenCode, Claude Code, die Claude-Desktop-App und die Modellgewichte sind **nicht enthalten**. Installiere sie selbst.
- Unabhängiges privates Projekt, nicht verbunden mit Anthropic, Alibaba Cloud (Qwen), Strata oder OpenCode.

## Aus dem Quellcode bauen

1. Installiere Go 1.25+, Node.js und PowerShell 5.1+.
2. Lege Strata in `strata\Strata-main` ab (oder setze `STRATA_DIR`) und erzeuge dessen `run-*.bat`.
3. Führe `npm install` in `opencode\` (OpenCode-Starter) und in `claudeui\` (Web-Oberfläche) aus.
4. Führe `.\build.ps1` aus, oder `.\build.ps1 -Only <Name>` mit einem von: tui, web, desktop, openwebui, strata, claude, claudewrap, claudeui, claudedesktop. Die Ausgabe landet in `dist\`.

## Mitmachen

Issues und Pull Requests sind willkommen, besonders: Testergebnisse für den Claude-Desktop-Weg, Geschwindigkeitswerte auf anderen GPUs und ein Linux-Port. Bitte gib deine GPU, deinen RAM und die Strata-Modellgröße an. Auch Übersetzungen lassen sich verbessern.

## Chatgruppe (optional)

WeChat-Gruppe „AI时代“, überwiegend auf Chinesisch (der QR-Code läuft am 14. Oktober ab; öffne ein Issue, falls er abgelaufen ist):

<img src="docs/wechat-group.png" alt="QR-Code der WeChat-Gruppe" width="260">

## Danksagung und Drittanbieter-Lizenzen

[Strata](https://github.com/Niko1221/Strata), [OpenCode](https://github.com/anomalyco/opencode) (MIT), [claude-code-webui](https://github.com/sugyan/claude-code-webui) (MIT) und das Qwen-Modell behalten jeweils ihre eigene Lizenz.
