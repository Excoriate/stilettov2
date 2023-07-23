package tui

import (
	"fmt"
	"github.com/excoriate/stiletto/internal/utils"
	"github.com/pterm/pterm"
	"github.com/pterm/pterm/putils"
	"strings"
)

type Title struct {
}

type ExecutionDetails struct {
	Workdir  string
	BaseDir  string
	MountDir string
	TaskName string
	Id       string
}

func (t *Title) ShowTitleAndDescription(title, description string) {
	titleNormalised := strings.TrimSpace(strings.ToUpper(title))
	s, _ := pterm.DefaultBigText.WithLetters(pterm.NewLettersFromString(titleNormalised)).
		Srender()
	pterm.DefaultCenter.Println(s)

	pterm.DefaultCenter.WithCenterEachLineSeparately().Println(description)
}

func (t *Title) ShowSubTitle(mainTitle string, subTitle string) {
	_ = pterm.DefaultBigText.WithLetters(
		putils.LettersFromStringWithStyle(strings.ToUpper(mainTitle), pterm.NewStyle(pterm.FgCyan)),
		putils.LettersFromStringWithStyle(strings.ToUpper(subTitle), pterm.NewStyle(pterm.FgLightMagenta))).
		Render()
}

func (t *Title) ShowTitle(title string) {
	titleNormalised := strings.TrimSpace(strings.ToUpper(title))
	s, _ := pterm.DefaultBigText.WithLetters(pterm.NewLettersFromString(titleNormalised)).
		Srender()
	pterm.DefaultCenter.Println(s)
}

func (t *Title) ShowExecutionDetails(opt ExecutionDetails) {
	pterm.Println()
	pterm.DefaultBasicText.Println(pterm.LightWhite("--------------------------------------------------"))
	pterm.DefaultBasicText.Println("Id" + pterm.LightMagenta(fmt.Sprintf(" %s ",
		utils.NormaliseStringUpper(opt.Id))))
	pterm.DefaultBasicText.Println("Task" + pterm.LightMagenta(fmt.Sprintf(" %s ", opt.TaskName)))
	pterm.DefaultBasicText.Println("Workdir" + pterm.LightMagenta(fmt.Sprintf(" %s ", opt.Workdir)))
	pterm.DefaultBasicText.Println("BaseDir" + pterm.LightMagenta(fmt.Sprintf(" %s ", opt.BaseDir)))
	pterm.DefaultBasicText.Println("MountDir" + pterm.LightMagenta(fmt.Sprintf(" %s ", opt.MountDir)))

	pterm.DefaultBasicText.Println(pterm.LightWhite("--------------------------------------------------"))
	pterm.Println()
	pterm.Println()
}

func NewTitle() UXTitleGenerator {
	return &Title{}
}
