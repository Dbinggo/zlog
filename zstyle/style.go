package zstyle

import (
	"fmt"
)

type Style string

const (
	Reset Style = "\033[0m"

	Red     Style = "\033[31m"
	Green   Style = "\033[32m"
	Yellow  Style = "\033[33m"
	Blue    Style = "\033[34m"
	Magenta Style = "\033[35m"
	Cyan    Style = "\033[36m"
	White   Style = "\033[37m"

	RedBold     Style = "\033[31;1m"
	GreenBold   Style = "\033[32;1m"
	YellowBold  Style = "\033[33;1m"
	BlueBold    Style = "\033[34;1m"
	MagentaBold Style = "\033[35;1m"
	CyanBold    Style = "\033[36;1m"
	WhiteBold   Style = "\033[36;1m"

	BackRed     Style = "\033[41m"
	BackGreen   Style = "\033[42m"
	BackYellow  Style = "\033[43m"
	BackBlue    Style = "\033[44m"
	BackMagenta Style = "\033[45m"
	BackCyan    Style = "\033[46m"
	BackWhite   Style = "\033[47m"
)

func SetStylef(style Style, format string, v ...any) string {

	return fmt.Sprintf("%s%s%s", style, fmt.Sprintf(format, v...), Reset)
}
