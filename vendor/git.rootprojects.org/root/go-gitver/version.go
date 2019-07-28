package main

// use recently generated version info as a fallback
// for when git isn't present (i.e. go run <url>)
func init() {
	GitRev = "0921ed1e6007493c886c87ee9a15f2cceecb1f9f"
	GitVersion = "v1.1.2"
	GitTimestamp = "2019-07-01T02:32:58-06:00"
}
