# [pathman](https://git.rootprojects.org/root/pathman)

Manage PATH on Windows, Mac, and Linux with various Shells

```bash
pathman list
pathman add ~/.local/bin
pathman remove ~/.local/bin
pathman version
pathman help
```

Windows: stores PATH in the registry.

Mac & Linux: stores PATH in `~/.config/envman/PATH.sh`

## Downloads

### MacOS

MacOS (darwin): [64-bit Download ](https://rootprojects.org/pathman/dist/darwin/amd64/pathman)

```
curl https://rootprojects.org/pathman/dist/darwin/amd64/pathman -o pathman
chmod +x ./pathman
```

### Windows

<details>
<summary>See download options</summary>
Windows 10: [64-bit Download](https://rootprojects.org/pathman/dist/windows/amd64/pathman.exe)

```
powershell.exe $ProgressPreference = 'SilentlyContinue'; Invoke-WebRequest https://rootprojects.org/pathman/dist/windows/amd64/pathman.exe -OutFile pathman.exe
```

Windows 7: [32-bit Download](https://rootprojects.org/pathman/dist/windows/386/pathman.exe)

```
powershell.exe "(New-Object Net.WebClient).DownloadFile('https://rootprojects.org/pathman/dist/windows/386/pathman.exe', 'pathman.exe')"
```

</details>

### Linux

<details>
<summary>See download options</summary>

Linux (64-bit): [Download](https://rootprojects.org/pathman/dist/linux/amd64/pathman)

```
curl https://rootprojects.org/pathman/dist/linux/amd64/pathman -o pathman
chmod +x ./pathman
```

Linux (32-bit): [Download](https://rootprojects.org/pathman/dist/linux/386/pathman)

```
curl https://rootprojects.org/pathman/dist/linux/386/pathman -o pathman
chmod +x ./pathman
```

</details>

### Raspberry Pi (Linux ARM)

<details>
<summary>See download options</summary>

RPi 4 (64-bit armv8): [Download](https://rootprojects.org/pathman/dist/linux/armv8/pathman)

```
curl https://rootprojects.org/pathman/dist/linux/armv8/pathman -o pathman`
chmod +x ./pathman
```

RPi 3 (armv7): [Download](https://rootprojects.org/pathman/dist/linux/armv7/pathman)

```
curl https://rootprojects.org/pathman/dist/linux/armv7/pathman -o pathman
chmod +x ./pathman
```

ARMv6: [Download](https://rootprojects.org/pathman/dist/linux/armv6/pathman)

```
curl https://rootprojects.org/pathman/dist/linux/armv6/pathman -o pathman
chmod +x ./pathman
```

RPi Zero (armv5): [Download](https://rootprojects.org/pathman/dist/linux/armv5/pathman)

```
curl https://rootprojects.org/pathman/dist/linux/armv5/pathman -o pathman
chmod +x ./pathman
```

</details>

## Install

1. Download (see below)
2. Add to `PATH`

**Windows**

```cmd
mkdir %userprofile%\bin
move pathman.exe %userprofile%\bin\pathman.exe
%userprofile%\bin\pathman.exe add ~\bin
```

**Mac, Linux, etc**

```bash
mkdir -p ~/.local/bin
mv ./pathman ~/.local/bin
pathman add ~/.local/bin
```

# CLI / API

# add

```bash
pathman add ~/.local/bin
```

```txt
Saved PATH changes. To set the PATH immediately, update the current session:

	export PATH="/Users/me/.local/bin:$PATH"
```

# remove

```bash
pathman remove ~/.local/bin
```

```txt
Saved PATH changes. To set the PATH immediately, update the current session:

	export PATH="/usr/local/bin:/usr/bin:/bin:/usr/sbin:/sbin"
```

# list

```bash
pathman list
```

```txt
pathman-managed PATH entries:

	$HOME/.local/bin

other PATH entries:

	/usr/local/bin
	/usr/bin
	/bin
	/usr/sbin
	/sbin

```

# Windows

You can use `~` as a shortcut for `%USERPROFILE%`.

```bash
pathman add ~\.local\bin
```

The registry will be used, even when your using Node Bash, Git Bash, or MINGW.

# build

```bash
git clone https://git.rootprojects.org/root/pathman.git
```

```bash
go mod tidy
go mod vendor
go generate -mod=vendor ./...
go build -mod=vendor
./pathman list
```
