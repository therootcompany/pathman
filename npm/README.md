# pathman

A cross-platform PATH manager

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

## Meta Package

This is a meta-package to fetch and install the correction version of
[go-pathman](https://git.rootprojects.org/root/pathman)
for your architecture and platform.

```bash
npm install -g @root/pathman
```

# Supported Shells

In theory, anything with bourne-compatible exports. Specifically:

-   bash
-   zsh
-   fish

On Windows, all shells inherit from the registry.

-   cmd.exe
-   PowerShell
-   Git Bash
