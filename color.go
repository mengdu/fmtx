package fmtx

import (
	"os"

	"github.com/mattn/go-isatty"
)

var enableColor = true
var forceColor = false
var cacheEnableColor = 0

func SetEnableColor(enable bool) {
	enableColor = enable
	cacheEnableColor = 0
}

func SetForceColor(force bool) {
	forceColor = force
	cacheEnableColor = 0
}

// reset enable color cache
func ResetColorCache() {
	cacheEnableColor = 0
}

func enableColorFromCache() bool {
	if cacheEnableColor != 0 {
		return cacheEnableColor == 1
	}

	if getEnableColor() {
		cacheEnableColor = 1
		return true
	} else {
		cacheEnableColor = 2
		return false
	}
}

func getEnableColor() bool {
	envVal := os.Getenv("FORCE_COLOR")
	envForceColor := envVal != "" && envVal != "0"
	return (forceColor || envForceColor) || ((isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd())) && enableColor)
}

func Color(s string, start string, end string) string {
	if !enableColorFromCache() {
		return s
	}
	// buf := make([]byte, 0, len(start)+len(s)+len(end))
	buf := []byte{}
	if start != "" {
		buf = append(buf, "\x1b["...)
		buf = append(buf, start...)
		buf = append(buf, 'm')
	}
	buf = append(buf, s...)
	if end != "" {
		buf = append(buf, "\x1b["...)
		buf = append(buf, end...)
		buf = append(buf, 'm')
	}
	return string(buf)
}
