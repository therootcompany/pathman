#!/bin/bin/env pwsh

# TODO use the git log describe thing
$my_version = git describe --tags
$Env:GOOS = "windows"

function fn_package() {
    IF ($Env:GOARM.Length -gt 0) {
        $my_arch = "${Env:GOARCH}v${Env:GOARM}"
    } ELSEIF ($Env:GOAMD64.Length -gt 0) {
        $my_arch = "${Env:GOARCH}_${Env:GOAMD64}"
    } ELSE {
        $my_arch = "${Env:GOARCH}"
    }

    $my_bin = "pathman-${my_version}-${Env:GOOS}-${my_arch}"
    tinygo build -no-debug -o "${my_bin}"
    #strip "${my_bin}" || true

    tar cvzf "$my_bin.tar.gz" "$my_bin"
    Compress-Archive "$my_bin" "$my_bin.zip"

    Write-Output "$my_bin.zip"
}

go generate ./...

$Env:GOAMD64 = "v1"
$Env:GOARCH = "amd64"
fn_package
$Env:GOAMD64 = ""

# $Env:GOARCH = "386"
# fn_package

$Env:GOARCH = "arm64"
fn_package

# $Env:GOARCH = "arm"
# $Env:GOARM = "7"
# fn_package

# $Env:GOARCH = "arm"
# $Env:GOARM = "6"
# fn_package

# unset vars
$Env:GOOS = ""
$Env:GOARCH = ""
$Env:GOARM = ""
$Env:GOAMD64 = ""
