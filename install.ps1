$X86_ID="84521240"
$A64_ID="84521239"
$ARCH_STRING=(Get-CimInstance Win32_operatingsystem).OSArchitecture
$IS_ARM=$ARCH_STRING -like "ARM"

Add-Type -AssemblyName System.IO.Compression.FileSystem
function Unzip
{
    param([string]$zipfile, [string]$outpath)
    [System.IO.Compression.ZipFile]::ExtractToDirectory($zipfile, $outpath)
}

if ($IS_ARM) {
    $downloadUri = "https://api.github.com/repos/ovrclk/akcmd/actions/artifacts/$A64_ID/zip"
}
else {
    $downloadUri = "https://api.github.com/repos/ovrclk/akcmd/actions/artifacts/$X86_ID/zip"
}

wget $downloadUri -outfile "$Env:USERPROFILE\file.zip" -Headers @{"Authorization"="Basic ZG1pa2V5OmdocF9kMXphdEs1VUtzamdaRWNUMmt5VW9HVk45dU5wa0YzaWdabFo="}

Unzip "$Env:USERPROFILE\file.zip" $Env:USERPROFILE