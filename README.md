# [pathman](https://git.rootprojects.org/root/pathman)

Manage PATH on Windows, Mac, and Linux with various Shells

```bash
pathman list
pathman add ~/.local/bin
pathman remove ~/.local/bin
```

Windows: stores PATH in the registry.

Mac & Linux: stores PATH in `~/.config/envman/PATH.sh`

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
