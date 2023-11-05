#!/bin/sh
set -e
set -u

# NOTE
#     building for macOS on macOS seems to yield smaller sizes

my_version="$(
    git describe --tags
)"
export GOOS="darwin"

fn_package() { (
    if test -n "${GOARM:-}"; then
        my_arch="${GOARCH}v${GOARM}"
    elif test -n "${GOAMD64:-}"; then
        my_arch="${GOARCH}_${GOAMD64}"
    else
        my_arch="${GOARCH}"
    fi

    my_bin="pathman-${my_version}-${GOOS}-${my_arch}"
    tinygo build -no-debug -o "${my_bin}"
    strip "${my_bin}" || true

    tar cvf "$my_bin.tar" "$my_bin"
    gzip --keep "$my_bin.tar"
    xz --keep "$my_bin.tar"

    echo "$my_bin.tar.xz"
); }

go generate ./...

export GOAMD64=v2
export GOARCH=amd64
fn_package
export GOAMD64=''

export GOARCH=arm64
fn_package
