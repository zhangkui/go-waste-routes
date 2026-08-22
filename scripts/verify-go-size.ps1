$ErrorActionPreference = 'Stop'
$files = Get-ChildItem cmd,internal -Recurse -File -Filter *.go | Where-Object { $_.Name -notlike '*_test.go' }
$lines = ($files | ForEach-Object { (Get-Content $_.FullName).Count } | Measure-Object -Sum).Sum
[pscustomobject]@{
  files = $files.Count
  lines = $lines
  paths = $files.FullName
} | ConvertTo-Json -Depth 3
