package fmtx

import (
	"fmt"
	"os"

	"github.com/mattn/go-isatty"
)

var EnableColor = true
var ForceColor = false

func color(s string, start string, end string) string {
	envForceColor := os.Getenv("FORCE_COLOR") != "" && os.Getenv("FORCE_COLOR") != "0"
	if (ForceColor || envForceColor) ||
		((isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd())) && EnableColor) {
		return fmt.Sprintf("\x1b[%sm%s\x1b[%sm", start, s, end)
	}
	return s
}
