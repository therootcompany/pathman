#!/bin/sh
set -e
set -u

# NOTE
#     building for linux on linux seems to yield smaller sizes

# TODO use the git log describe thing
my_version="$(
    git describe --tags
)"
export GOOS="linux"

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

export GOAMD64=v1
export GOARCH=amd64
fn_package
export GOAMD64=''

export GOARCH=386
fn_package

export GOARCH=arm64
fn_package

export GOARCH=arm
export GOARM=7
fn_package

export GOARCH=arm
export GOARM=6
fn_package

# no longer supported
# export GOARCH=arm
# export GOARM=5
# fn_package
