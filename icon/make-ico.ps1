# Build a multi-size .ico from an .svg using headless Edge. Usage: make-ico.ps1 -Name openwebui-qwen
# Reads <Name>.svg in this folder and writes <Name>.ico (PNG-compressed, 16..256 px). ASCII only.
param([Parameter(Mandatory = $true)][string]$Name, [switch]$KeepPng)
$ErrorActionPreference = 'Stop'
$dir = $PSScriptRoot
$svg = Join-Path $dir "$Name.svg"
if (-not (Test-Path $svg)) { throw "missing $svg" }
$edge = @("${env:ProgramFiles(x86)}\Microsoft\Edge\Application\msedge.exe", "$env:ProgramFiles\Microsoft\Edge\Application\msedge.exe") |
        Where-Object { Test-Path $_ } | Select-Object -First 1
$sizes = 16, 24, 32, 48, 64, 128, 256
$uri = 'file:///' + ($svg -replace '\\', '/')
foreach ($s in $sizes) {
    $html = Join-Path $dir "_r$s.html"; $png = Join-Path $dir "_p$s.png"
    $page = '<!doctype html><html><body style="margin:0;background:transparent"><img src="' + $uri + '" width="' + $s + '" height="' + $s + '" style="display:block"></body></html>'
    [IO.File]::WriteAllText($html, $page, (New-Object Text.UTF8Encoding($false)))
    Remove-Item $png -ErrorAction SilentlyContinue
    $ErrorActionPreference = 'Continue'   # Edge prints a status line to stderr
    & $edge --headless=new --disable-gpu --hide-scrollbars --default-background-color=00000000 --window-size=$s,$s --screenshot=$png ('file:///' + ($html -replace '\\', '/')) 2>&1 | Out-Null
    $ErrorActionPreference = 'Stop'
    for ($i = 0; $i -lt 25 -and -not (Test-Path $png); $i++) { Start-Sleep -Milliseconds 400 }
    if (-not (Test-Path $png)) { throw "render failed at $s px" }
}
$ms = New-Object IO.MemoryStream; $bw = New-Object IO.BinaryWriter($ms)
$bw.Write([uint16]0); $bw.Write([uint16]1); $bw.Write([uint16]$sizes.Count)
$datas = $sizes | ForEach-Object { , [IO.File]::ReadAllBytes((Join-Path $dir "_p$_.png")) }
$offset = 6 + 16 * $sizes.Count
for ($i = 0; $i -lt $sizes.Count; $i++) {
    $s = $sizes[$i]; $d = $datas[$i]; $b = if ($s -ge 256) { 0 } else { $s }
    $bw.Write([byte]$b); $bw.Write([byte]$b); $bw.Write([byte]0); $bw.Write([byte]0)
    $bw.Write([uint16]1); $bw.Write([uint16]32); $bw.Write([uint32]$d.Length); $bw.Write([uint32]$offset); $offset += $d.Length
}
foreach ($d in $datas) { $bw.Write($d) }
[IO.File]::WriteAllBytes((Join-Path $dir "$Name.ico"), $ms.ToArray()); $bw.Close()
if ($KeepPng) { Copy-Item (Join-Path $dir '_p256.png') (Join-Path $dir "$Name.preview.png") -Force }
Get-ChildItem $dir -Filter '_*' | Remove-Item -Force
Get-Item (Join-Path $dir "$Name.ico") | Format-List Name, Length
