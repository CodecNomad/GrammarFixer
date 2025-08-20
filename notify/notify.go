// Package notify provides a simple interface to send desktop notifications on Linux systems using `notify-send`.
package notify

import (
	"os/exec"
	"strconv"
)

type NotificationLevel int

const (
	NotificationLevelError = iota
	NotificationLevelWarn
	NotificationLevelInfo
)

func SendNotification(text string, level NotificationLevel) error {
	return SendTimedNotification(text, level, 1500)
}

func SendTimedNotification(text string, level NotificationLevel, timeOutMs int) error {
	var title string
	switch level {
	case NotificationLevelError:
		title = "[ERROR]"
	case NotificationLevelWarn:
		title = "[WARNING]"
	case NotificationLevelInfo:
		title = "[INFO]"
	default:
		title = "[UNKNOWN]"
	}
	title += " GrammarFixer"

	cmd := exec.Command("notify-send", "-t", strconv.Itoa(timeOutMs), title, text)
	err := cmd.Run()
	return err
}
