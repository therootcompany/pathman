# winpath

An example of getting, setting, and broadcasting PATHs on Windows.

This requires the `unsafe` package to use a syscall with special message poitners to update `PATH` without a reboot.
It will also build without `unsafe`.

```bash
go build -tags unsafe -o winpath.exe
```

```bash
winpath show

        %USERPROFILE%\AppData\Local\Microsoft\WindowsApps
        C:\Users\me\AppData\Local\Programs\Microsoft VS Code\bin
        %USERPROFILE%\go\bin
        C:\Users\me\AppData\Roaming\npm
        C:\Users\me\AppData\Local\Keybase\
```

```bash
winpath append C:\someplace\special

	Run the following for changes to take affect immediately:
	PATH %PATH%;C:\someplace\special
```

```bash
winpath prepend C:\someplace\special

	Run the following for changes to take affect immediately:
	PATH C:\someplace\special;%PATH%
```

```bash
winpath remove C:\someplace\special
```

# Special Considerations

Giving away the secret sauce right here:

* `HWND_BROADCAST`
* `WM_SETTINGCHANGE`

This is essentially the snippet you need to have the HKCU and HKLM Environment registry keys propagated without rebooting:

```go
	HWND_BROADCAST   := uintptr(0xffff)
	WM_SETTINGCHANGE := uintptr(0x001A)
	_, _, err := syscall.
		NewLazyDLL("user32.dll").
		NewProc("SendMessageW").
		Call(HWND_BROADCAST, WM_SETTINGCHANGE, 0, uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr("ENVIRONMENT"))))

```

* `os.Getenv("COMSPEC")`
* `os.Getenv("SHELL")`

If you check `SHELL` and it isn't empty, then you're probably in MINGW or some such.
If that's empty but `COMSPEC` isn't, you can be reasonably sure that you're in cmd.exe or Powershell.
