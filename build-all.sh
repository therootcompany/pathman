#GOOS=windows GOARCH=amd64 go install
#go tool dist list

# TODO move this into tools/build.go

export CGO_ENABLED=0
exe=pathman
gocmd=.

echo ""
go generate -mod=vendor ./...

echo ""
echo "Windows amd64"
#GOOS=windows GOARCH=amd64 go build -mod=vendor -o dist/windows/amd64/${exe}.exe -ldflags "-s -w -H=windowsgui" $gocmd
GOOS=windows GOARCH=amd64 go build -mod=vendor -o dist/windows/amd64/${exe}.exe -ldflags "-s -w" $gocmd
#GOOS=windows GOARCH=amd64 go build -mod=vendor -o dist/windows/amd64/${exe}.debug.exe
echo "Windows 386"
GOOS=windows GOARCH=386 go build -mod=vendor -o dist/windows/386/${exe}.exe -ldflags "-s -w" $gocmd

echo ""
echo "Darwin (macOS) amd64"
GOOS=darwin GOARCH=amd64 go build -mod=vendor -o dist/darwin/amd64/${exe} -ldflags "-s -w" $gocmd

echo ""
echo "Linux amd64"
GOOS=linux GOARCH=amd64 go build -mod=vendor -o dist/linux/amd64/${exe} -ldflags "-s -w" $gocmd
echo "Linux 386"
GOOS=linux GOARCH=386 go build -mod=vendor -o dist/linux/386/${exe} -ldflags "-s -w" $gocmd

echo ""
echo "RPi 4 (64-bit) ARMv8"
GOOS=linux GOARCH=arm64 go build -mod=vendor -o dist/linux/armv8/${exe} -ldflags "-s -w" $gocmd
echo "RPi 3 B+ ARMv7"
GOOS=linux GOARCH=arm GOARM=7 go build -mod=vendor -o dist/linux/armv7/${exe} -ldflags "-s -w" $gocmd
echo "ARMv6"
GOOS=linux GOARCH=arm GOARM=6 go build -mod=vendor -o dist/linux/armv6/${exe} -ldflags "-s -w" $gocmd
echo "RPi Zero ARMv5"
GOOS=linux GOARCH=arm GOARM=5 go build -mod=vendor -o dist/linux/armv5/${exe} -ldflags "-s -w" $gocmd

echo ""
#rsync -av ./dist/ ubuntu@rootprojects.org:/srv/www/rootprojects.org/pathman/dist/
# https://rootprojects.org/pathman/dist/windows/amd64/pathman.exe
