package tui

type UXMessenger interface {
	ShowError(title, msg string, err error)
	ShowInfo(title, msg string)
	ShowSuccess(title, msg string)
	ShowWarning(title, msg string)
}

type UXTitleGenerator interface {
	ShowTitleAndDescription(title, description string)
	ShowTitle(title string)
	ShowSubTitle(mainTitle, subtitle string)
	ShowInitDetails(opt JobTUIDetails)
}
