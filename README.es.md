<p align="center"><img src="docs/banner.png" alt="open-lowmem: ejecuta LLM enormes con muy poca memoria" width="720"></p>

# open-lowmem: ejecuta un LLM de 125B en una GPU de 24 GB y 32 GB de RAM (Windows)

[English](README.md) | [简体中文](README.zh-CN.md) | [繁體中文](README.zh-TW.md) | [日本語](README.ja.md) | [한국어](README.ko.md) | [Español](README.es.md) | [Français](README.fr.md) | [Deutsch](README.de.md) | [Русский](README.ru.md) | [Português](README.pt-BR.md)

**Ejecuta LLM enormes con muy poca memoria.** Lanzadores de Windows con un clic que ejecutan en local un **modelo MoE de 125 000 millones de parámetros (Qwen3.8-Flash-Next)** en un PC doméstico con **una tarjeta gráfica de 24 GB y 32 GB de RAM**, y lo conectan con **Claude Code, OpenCode, la aplicación de escritorio de Claude y Open WebUI**. Sin conexión, privado y gratuito.

Autor: **廖亦辰 (Liao Yichen)** · Xiamen, China · GitHub [@bathfire-dev](https://github.com/bathfire-dev) · Licencia: MIT

> Si esto te ahorra un fin de semana de configuración, una estrella ayuda a que otras personas con GPU pequeñas lo encuentren.

## ¿Qué es esto?

`open-lowmem` es un conjunto de pequeños lanzadores en Go para Windows. Cada uno arranca primero el motor de inferencia local [Strata](https://github.com/Niko1221/Strata), que hace que un modelo de mezcla de expertos de 125B parámetros quepa en un PC gamer normal repartiendo el trabajo entre la VRAM de la GPU, la RAM del sistema y el SSD. Después lanza el cliente que elijas y lo limpia todo al cerrar la ventana.

Palabras clave: LLM local, poca VRAM, poca memoria, ejecutar modelos de lenguaje grandes en una GPU de 24 GB, MoE, cuantización (IQ2_XS), Qwen3.8-Flash-Next, Strata, Claude Code con un modelo local, OpenCode, asistente de programación con IA sin conexión, RTX 4090, Windows.

## ¿Cuánta memoria necesita?

Las cifras proceden de la [documentación de Strata](https://github.com/Niko1221/Strata), no de mediciones de este repositorio:

| Tamaño (mismo modelo) | RAM + VRAM necesarias | Notas |
| --- | --- | --- |
| Q2_0 | unos 37,6 GB | el más rápido |
| IQ2_XS | unos 39,2 GB | recomendado; configuración por defecto de este repositorio |
| IQ3_XXS / IQ3_S | unos 47 / 55 GB | mejor calidad, necesita más RAM |

Strata necesita una tarjeta NVIDIA o AMD de 12 GB o más, 32 GB o más de RAM y unos 80 GB de disco libre. Sus autores midieron 79 tokens/s en una RTX 5070 (12 GB) con 64 GB de RAM para IQ2_XS. Con una tarjeta de 24 GB y 32 GB de RAM, Q2_0 e IQ2_XS funcionan en el modo de poca RAM de Strata. Este repositorio se desarrolla en una RTX 4090D de 24 GB con 32 GB de RAM.

## Lanzadores

| Lanzador | Qué hace |
| --- | --- |
| `claude-strata.exe` | Arranca Strata y ejecuta tu Claude Code instalado con `--settings`. No toca tu `~/.claude`. |
| `claude-strata-ui.exe` | Selector de carpeta y, después, Strata y la interfaz web `claude-code-webui`. |
| `claude-desktop-strata.exe` | Arranca Strata, escribe un perfil de pasarela de terceros para la app de escritorio de Claude y la inicia. |
| `opencode-qwen.exe`, `opencode-qwen-web.exe` | OpenCode en terminal y web, con la configuración incluida. |
| `opencode-qwen-strata-desktop.exe` | Arranca Strata y luego la app de escritorio de OpenCode. |
| `openwebui-qwen.exe` | Ollama más Open WebUI. |

Al cerrar la ventana del lanzador se terminan todos los procesos hijos mediante un Job Object de Windows, así ningún modelo sigue ocupando tu VRAM.

## Preguntas frecuentes

**¿Puedo ejecutar un modelo de 125B parámetros en una GPU de 24 GB?**
Sí, con un modelo de mezcla de expertos y un motor que reparta los expertos entre VRAM, RAM y SSD. Eso es lo que hace Strata; este repositorio solo lo conecta con tus herramientas. Cuenta con unos 39 GB entre RAM y VRAM para IQ2_XS.

**¿Cómo uso Claude Code con un modelo local?**
Strata expone un endpoint compatible con Anthropic (`/v1/messages` en `127.0.0.1:8090`). `claude-strata.exe` arranca Strata y lanza Claude Code con una capa de configuración que apunta a él, usando como alias un nombre de modelo de Claude como `claude-sonnet-4-6`. Claude Code no está incluido; instálalo desde Anthropic.

**¿Funciona con la app de escritorio de Claude?**
No está verificado. El perfil se escribe, pero no se sabe si la app acepta una pasarela `http://127.0.0.1:8090` sin HTTPS. Considera `claude-desktop-strata.exe` experimental.

**¿Tiene la misma calidad que los modelos en la nube?**
No. La cuantización de pocos bits cambia calidad por memoria. En las clasificaciones públicas Qwen3.8-Flash-Next es un buen modelo abierto, pero no un modelo cerrado de primera línea, y parte de esas puntuaciones las comunica el propio proveedor.

**¿Funciona en Linux o macOS?**
Ahora no. Los lanzadores usan APIs de Windows.

## Notas importantes

- Strata, OpenCode, Claude Code, la app de escritorio de Claude y los pesos del modelo **no están incluidos**. Instálalos tú mismo.
- Proyecto personal independiente, sin relación con Anthropic, Alibaba Cloud (Qwen), Strata ni OpenCode.

## Compilar desde el código fuente

1. Instala Go 1.25+, Node.js y PowerShell 5.1+.
2. Coloca Strata en `strata\Strata-main` (o define `STRATA_DIR`) y genera su `run-*.bat`.
3. Ejecuta `npm install` en `opencode\` (lanzadores de OpenCode) y en `claudeui\` (interfaz web).
4. Ejecuta `.\build.ps1`, o `.\build.ps1 -Only <nombre>` con uno de: tui, web, desktop, openwebui, strata, claude, claudewrap, claudeui, claudedesktop. La salida va a `dist\`.

## Contribuir

Se agradecen issues y pull requests, sobre todo: resultados de prueba de la ruta de escritorio de Claude, cifras de velocidad en otras GPU y un port a Linux. Indica tu GPU, tu RAM y el tamaño de modelo de Strata. Las traducciones también se pueden mejorar.

## Grupo de chat (opcional)

Grupo de WeChat «AI时代», principalmente en chino (el código QR caduca el 14 de octubre; abre un issue si ha caducado):

<img src="docs/wechat-group.png" alt="Código QR del grupo de WeChat" width="260">

## Créditos y licencias de terceros

[Strata](https://github.com/Niko1221/Strata), [OpenCode](https://github.com/anomalyco/opencode) (MIT), [claude-code-webui](https://github.com/sugyan/claude-code-webui) (MIT) y el modelo Qwen conservan cada uno su propia licencia.
