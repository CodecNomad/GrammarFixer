// Package clipboard provides functions to read from and write to the system clipboard on wayland systems.
package clipboard

import (
	"bytes"
	"os/exec"
)

func ReadClipboardText() (string, error) {
	var buf bytes.Buffer
	cmd := exec.Command("wl-paste")
	cmd.Stdout = &buf
	err := cmd.Run()
	return buf.String(), err
}

func WriteToClipboard(text string) error {
	cmd := exec.Command("wl-copy")
	cmd.Stdin = bytes.NewBufferString(text)
	return cmd.Run()
}
