package fmtx

import (
	"os"
	"strings"

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

func color(s string, start string, end string) string {
	if !enableColorFromCache() {
		return s
	}
	arr := []string{"\x1b[", start, "m", s, "\x1b[", end, "m"}
	if start == "" {
		arr[0] = ""
		arr[2] = ""
	}
	if end == "" {
		arr[4] = ""
		arr[6] = ""
	}
	return strings.Join(arr, "")
}
