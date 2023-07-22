package tui

import (
	"fmt"
	"github.com/pterm/pterm"
	"strings"
)

type Message struct {
}

func (t *Message) ShowError(title, msg string, err error) {
	if title != "" {
		pterm.Error.Prefix = pterm.Prefix{
			Text:  strings.ToUpper(title),
			Style: pterm.NewStyle(pterm.BgCyan, pterm.FgRed),
		}
	}

	if err == nil {
		pterm.Error.Println(msg)
		return
	}

	var errMsg string
	if msg != "" {
		errMsg = fmt.Sprintf("%s: %s", msg, err)
	} else {
		errMsg = err.Error()
	}

	pterm.Error.Println(errMsg)
}

func (t *Message) ShowInfo(title, msg string) {
	if title != "" {
		pterm.Info.Prefix = pterm.Prefix{
			Text:  strings.ToUpper(title),
			Style: pterm.NewStyle(pterm.BgCyan, pterm.FgBlack),
		}
	}

	pterm.Info.Println(msg)
}

func (t *Message) ShowSuccess(title, msg string) {
	if title != "" {
		pterm.Success.Prefix = pterm.Prefix{
			Text:  strings.ToUpper(title),
			Style: pterm.NewStyle(pterm.BgCyan, pterm.FgBlack),
		}
	}

	pterm.Success.Println(msg)
}

func (t *Message) ShowWarning(title, msg string) {
	if title != "" {
		pterm.Warning.Prefix = pterm.Prefix{
			Text:  strings.ToUpper(title),
			Style: pterm.NewStyle(pterm.BgCyan, pterm.FgBlack),
		}
	}

	pterm.Warning.Println(msg)
}

func NewTUIMessage() UXMessenger {
	return &Message{}
}
