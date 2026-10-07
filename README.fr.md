<p align="center"><img src="docs/banner.png" alt="open-lowmem : faire tourner d'énormes LLM avec très peu de mémoire" width="720"></p>

# open-lowmem : faire tourner un LLM de 125 Md de paramètres à ~75 tokens/s sur un GPU de 24 Go et 32 Go de RAM (Windows)

[English](README.md) | [简体中文](README.zh-CN.md) | [繁體中文](README.zh-TW.md) | [日本語](README.ja.md) | [한국어](README.ko.md) | [Español](README.es.md) | [Français](README.fr.md) | [Deutsch](README.de.md) | [Русский](README.ru.md) | [Português](README.pt-BR.md)

**Faites tourner d'énormes LLM avec très peu de mémoire.** Des lanceurs Windows en un clic qui exécutent en local un **modèle MoE de 125 milliards de paramètres (Qwen3.8-Flash-Next)** sur un PC grand public équipé d'**une carte graphique de 24 Go et de 32 Go de RAM**, puis le branchent à **Claude Code, OpenCode, l'application de bureau Claude et Open WebUI**. Hors ligne, privé, gratuit.

Auteur : **廖亦辰 (Liao Yichen)** · Xiamen, Chine · GitHub [@bathfire-dev](https://github.com/bathfire-dev) · Licence : MIT

> Si ce projet vous épargne un week-end de configuration, une étoile aide d'autres personnes avec de petits GPU à le trouver.

## De quoi s'agit-il ?

`open-lowmem` est un ensemble de petits lanceurs Go pour Windows. Chacun démarre d'abord le moteur d'inférence local [Strata](https://github.com/Niko1221/Strata), qui permet à un modèle à mélange d'experts de 125 Md de paramètres de tenir dans un PC de jeu ordinaire en répartissant le travail entre la VRAM du GPU, la RAM du système et le SSD. Le lanceur démarre ensuite le client choisi et nettoie tout à la fermeture de la fenêtre.

Mots-clés : LLM local, peu de VRAM, faible mémoire, exécuter de grands modèles de langage sur un GPU de 24 Go, MoE, quantification (IQ2_XS), Qwen3.8-Flash-Next, Strata, Claude Code avec un modèle local, OpenCode, assistant de programmation IA hors ligne, RTX 4090, Windows.

## Vitesse mesurée

Mesurée sur la machine de développement de ce dépôt (RTX 4090D 24 Go, 32 Go de RAM, IQ2_XS, Windows, 7 octobre 2026) : 8 exécutions de 512 tokens générés ont donné **entre 62 et 94 tokens/s, avec une médiane d'environ 77**. Une requête à la fois, prompt court. Votre vitesse variera selon la longueur du contexte, la charge en arrière-plan et la taille du modèle.

## De combien de mémoire a-t-on besoin ?

Les chiffres proviennent de la [documentation de Strata](https://github.com/Niko1221/Strata), pas de mesures faites dans ce dépôt :

| Taille (même modèle) | RAM + VRAM requises | Remarques |
| --- | --- | --- |
| Q2_0 | environ 37,6 Go | le plus rapide |
| IQ2_XS | environ 39,2 Go | recommandé ; configuration par défaut de ce dépôt |
| IQ3_XXS / IQ3_S | environ 47 / 55 Go | meilleure qualité, demande plus de RAM |

Strata demande une carte NVIDIA ou AMD de 12 Go ou plus, 32 Go de RAM ou plus et environ 80 Go d'espace disque libre. Ses auteurs ont mesuré 79 tokens/s sur une RTX 5070 (12 Go) avec 64 Go de RAM pour IQ2_XS. Avec une carte de 24 Go et 32 Go de RAM, Q2_0 et IQ2_XS fonctionnent dans le mode « peu de RAM » de Strata. Ce dépôt est développé sur une RTX 4090D de 24 Go avec 32 Go de RAM.

## Lanceurs

| Lanceur | Rôle |
| --- | --- |
| `claude-strata.exe` | Démarre Strata, puis exécute votre Claude Code installé avec `--settings`. Votre `~/.claude` n'est pas modifié. |
| `claude-strata-ui.exe` | Sélecteur de dossier, puis Strata et l'interface web `claude-code-webui`. |
| `claude-desktop-strata.exe` | Démarre Strata, écrit un profil de passerelle tierce pour l'application de bureau Claude, puis la lance. |
| `opencode-qwen.exe`, `opencode-qwen-web.exe` | OpenCode en terminal et en web, configuration intégrée. |
| `opencode-qwen-strata-desktop.exe` | Démarre Strata, puis l'application de bureau OpenCode. |
| `openwebui-qwen.exe` | Ollama plus Open WebUI. |

Fermer la fenêtre du lanceur termine tous les processus enfants via un Job Object Windows : aucun modèle ne continue à occuper votre VRAM.

## FAQ

**Puis-je faire tourner un modèle de 125 Md de paramètres sur un GPU de 24 Go ?**
Oui, avec un modèle à mélange d'experts et un moteur qui répartit les experts entre VRAM, RAM et SSD. C'est ce que fait Strata ; ce dépôt ne fait que le relier à vos outils. Comptez environ 39 Go de RAM et VRAM cumulées pour IQ2_XS.

**Comment utiliser Claude Code avec un modèle local ?**
Strata expose un point d'accès compatible Anthropic (`/v1/messages` sur `127.0.0.1:8090`). `claude-strata.exe` démarre Strata et lance Claude Code avec une surcouche de réglages qui pointe vers lui, en utilisant comme alias un nom de modèle Claude tel que `claude-sonnet-4-6`. Claude Code n'est pas fourni ; installez-le depuis Anthropic.

**L'application de bureau Claude fonctionne-t-elle ?**
Non vérifié. Le profil est écrit, mais on ignore si l'application accepte une passerelle `http://127.0.0.1:8090` sans HTTPS. Considérez `claude-desktop-strata.exe` comme expérimental.

**La qualité est-elle équivalente à celle des modèles cloud ?**
Non. La quantification à faible nombre de bits échange de la qualité contre de la mémoire. Dans les classements publics, Qwen3.8-Flash-Next est un bon modèle ouvert, mais pas un modèle fermé de pointe, et une partie de ces scores est déclarée par l'éditeur.

**Cela fonctionne-t-il sous Linux ou macOS ?**
Pas pour l'instant. Les lanceurs utilisent des API Windows.

## Remarques importantes

- Strata, OpenCode, Claude Code, l'application de bureau Claude et les poids du modèle **ne sont pas inclus**. Installez-les vous-même.
- Projet personnel indépendant, sans lien avec Anthropic, Alibaba Cloud (Qwen), Strata ou OpenCode.

## Compiler depuis les sources

1. Installez Go 1.25+, Node.js et PowerShell 5.1+.
2. Placez Strata dans `strata\Strata-main` (ou définissez `STRATA_DIR`) et générez son `run-*.bat`.
3. Lancez `npm install` dans `opencode\` (lanceurs OpenCode) et dans `claudeui\` (interface web).
4. Lancez `.\build.ps1`, ou `.\build.ps1 -Only <nom>` avec l'un de : tui, web, desktop, openwebui, strata, claude, claudewrap, claudeui, claudedesktop. La sortie est dans `dist\`.

## Contribuer

Les issues et pull requests sont les bienvenues, en particulier : un résultat de test du chemin application de bureau Claude, des mesures de vitesse sur d'autres GPU et un portage Linux. Indiquez votre GPU, votre RAM et la taille de modèle Strata. Les traductions peuvent aussi être améliorées.

## Groupe de discussion (facultatif)

Groupe WeChat « AI时代 », surtout en chinois (le QR code expire le 14 octobre ; ouvrez une issue s'il a expiré) :

<img src="docs/wechat-group.png" alt="QR code du groupe WeChat" width="260">

## Crédits et licences tierces

[Strata](https://github.com/Niko1221/Strata), [OpenCode](https://github.com/anomalyco/opencode) (MIT), [claude-code-webui](https://github.com/sugyan/claude-code-webui) (MIT) et le modèle Qwen conservent chacun leur propre licence.
