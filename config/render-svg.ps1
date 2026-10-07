# Render an SVG to a PNG (same folder, same name) with headless Edge, so it can be viewed.
# Usage: render-svg.ps1 -Svg path\to\file.svg [-Size 800]
param(
    [Parameter(Mandatory = $true)][string]$Svg,
    [int]$Size = 800
)
$ErrorActionPreference = 'Stop'
$full = (Resolve-Path -LiteralPath $Svg).Path
$png  = [IO.Path]::ChangeExtension($full, '.png')
$html = Join-Path ([IO.Path]::GetTempPath()) ("svgview_" + [guid]::NewGuid().ToString('N') + ".html")
$edge = @("${env:ProgramFiles(x86)}\Microsoft\Edge\Application\msedge.exe", "$env:ProgramFiles\Microsoft\Edge\Application\msedge.exe") |
        Where-Object { Test-Path $_ } | Select-Object -First 1
if (-not $edge) { throw 'Microsoft Edge not found' }
$uri = 'file:///' + ($full -replace '\\', '/')
$page = '<!doctype html><html><body style="margin:0;background:#fff"><img src="' + $uri + '" style="display:block;width:100vw;height:100vh;object-fit:contain"></body></html>'
[IO.File]::WriteAllText($html, $page, (New-Object Text.UTF8Encoding($false)))
if (Test-Path $png) { Remove-Item $png -Force }
$ErrorActionPreference = 'Continue'   # Edge prints a status line to stderr
& $edge --headless=new --disable-gpu --hide-scrollbars --window-size=$Size,$Size --screenshot=$png ('file:///' + ($html -replace '\\', '/')) 2>&1 | Out-Null
$ErrorActionPreference = 'Stop'
for ($i = 0; $i -lt 25 -and -not (Test-Path $png); $i++) { Start-Sleep -Milliseconds 400 }
Remove-Item $html -Force -ErrorAction SilentlyContinue
if (Test-Path $png) { Write-Output "Rendered: $png" } else { Write-Error "Render failed for $full"; exit 1 }
