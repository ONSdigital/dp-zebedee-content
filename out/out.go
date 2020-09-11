package out

import (
	"fmt"

	"github.com/fatih/color"
)

var (
	infoBoldC    = color.New(color.Bold, color.FgHiGreen)
	infoC        = color.New(color.FgCyan)
	outPrefix    = "[generate]"
)

type Level int

const (
	INFO Level = iota + 1
	WARN
	ERROR
)

func cliPrefix(c *color.Color) {
	c.Printf("%s ", outPrefix)
}

func Info(msg string) {
	cliPrefix(infoBoldC)
	fmt.Printf("%s\n", msg)
}

func InfoAppend(msg string) {
	infoC.Print(msg)
}

func InfoAppendF(msg string, args ...interface{}) {
	infoC.Printf(msg, args...)
}

func InfoF(msg string, args ...interface{}) {
	cliPrefix(infoBoldC)
	fmt.Printf(msg, args...)
}

func InfoFHighlight(msg string, args ...interface{}) {
	cliPrefix(infoBoldC)
	highlight(infoC, msg, args...)
}

func highlight(c *color.Color, formattedMsg string, args ...interface{}) {
	var highlighted []interface{}

	for _, val := range args {
		highlighted = append(highlighted, c.SprintFunc()(val))
	}

	formattedMsg = fmt.Sprintf(formattedMsg, highlighted...)
	fmt.Printf("%s\n", formattedMsg)
}
