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

type JobTUIDetails struct {
	Product     string
	Workdir     string
	Domain      string
	Service     string
	TypeOrLayer string
	Environment string
	Region      string
	// Shows the task to run. E.g.: 'terraform apply'
	TaskName string
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

func (t *Title) ShowInitDetails(opt JobTUIDetails) {
	pterm.Println()
	pterm.DefaultBasicText.Println(pterm.LightWhite("--------------------------------------------------"))
	pterm.DefaultBasicText.Println("Product" + pterm.LightMagenta(fmt.Sprintf(" %s ",
		utils.NormaliseStringUpper(opt.Product))))

	pterm.DefaultBasicText.Println("Domain ->" + pterm.LightMagenta(fmt.Sprintf(" %s ",
		utils.NormaliseStringUpper(opt.Domain))))
	pterm.DefaultBasicText.Println("Service ->" + pterm.LightMagenta(fmt.Sprintf(" %s ",
		utils.NormaliseStringUpper(opt.Service))))
	pterm.DefaultBasicText.Println("Type/Layer ->" + pterm.LightMagenta(fmt.Sprintf(" %s ",
		utils.NormaliseStringUpper(opt.TypeOrLayer))))
	pterm.DefaultBasicText.Println("Environment ->" + pterm.LightMagenta(fmt.Sprintf(" %s ",
		utils.NormaliseStringUpper(opt.Environment))))
	pterm.DefaultBasicText.Println("Region ->" + pterm.LightMagenta(fmt.Sprintf(" %s ",
		utils.NormaliseStringUpper(opt.Region))))
	pterm.DefaultBasicText.Println("Workdir ->" + pterm.LightMagenta(fmt.Sprintf(" %s ",
		utils.NormaliseStringUpper(opt.Workdir))))
	pterm.DefaultBasicText.Println("Task ->" + pterm.LightMagenta(fmt.Sprintf(" %s ",
		utils.NormaliseStringUpper(opt.TaskName))))

	pterm.DefaultBasicText.Println(pterm.LightWhite("--------------------------------------------------"))
	pterm.Println()
	pterm.Println()
}

func NewTitle() UXTitleGenerator {
	return &Title{}
}
