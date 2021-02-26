# [pathman](https://git.rootprojects.org/root/pathman)

Manage PATH on **Windows 10**, **Mac**, and **Linux** with various Shells

```bash
pathman list
pathman add ~/.local/bin
pathman remove ~/.local/bin
pathman version
pathman help
```

Where is the PATH managed?

-   **Windows 10**: stores `PATH` in the registry.
-   **Mac** & **Linux**: stores `PATH` in `~/.config/envman/PATH.env`

Note for **Windows 10** users: due to differences in how `cmd.exe`, PowerShell, and `pathman` use and interpret strings, spaces, paths, and variables, you'll get more consistent results if you:

-   Use `~` rather than `%USERPROFILE%` or `$Env:USERPROFILE`
-   Use `/` rather than `\` for delimiting paths

## Install

**Mac**, **Linux**:

```bash
curl -s https://webinstall.dev/pathman | bash
```

**Windows 10**:

This can be run from `cmd.exe` or PowerShell (`curl.exe` is a native part of Windows 10).

```bash
curl.exe -sA "MS" https://webinstall.dev/pathman | powershell
```

### Manual Install

1. [Download](#downloads)
2. Add to `PATH`

Or install via `npm`:

```bash
npm install -g pathman
```

#### Windows

```cmd
mkdir %userprofile%\bin
move pathman.exe %userprofile%\bin\pathman.exe
%userprofile%\bin\pathman.exe add ~/bin
```

#### Mac, Linux, etc

```bash
mkdir -p ~/.local/bin
mv ./pathman ~/.local/bin
pathman add ~/.local/bin
```

## Downloads

[Webi](https://webinstall.dev/pathman) (<https://webinstall.dev/pathman>) is the preferred install method,
but you can also download from [Git Releases](https://git.rootprojects.org/root/pathman/releases):
<https://git.rootprojects.org/root/pathman/releases>.

MacOS (including Apple Silicon M1), Linux, Raspberry Pi:

```bash
tar xvf pathman-v*.tar.gz
chmod a+x ./pathman
./pathman --help
```

Windows 10:

```bash
tar.exe xvf pathman-v*.zip
.\pathman.exe --help
```

### Supported Platforms

-   MacOS
    -   Apple Silicon M1
    -   Intel x86_64
-   Windows 10, 8, 7
-   Linux
    -   amd64 / x86_64
    -   386
-   Raspberry Pi (Linux ARM)
    -   RPi 4 (64-bit armv8)
    -   RPi 3 (armv7)
    -   ARMv6
    -   RPi Zero (armv5)

# CLI Help (API)

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
