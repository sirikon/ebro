package logger

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

func Info(text string) {
	Line(color.FgGreen, text)
}

func Notice(text string) {
	Line(color.FgYellow, text)
}

func Error(text string) {
	color.New(color.FgHiRed).Add(color.Bold).Print("███ ERROR: ")
	fmt.Println(normalizeFinalNewLine(text))
}

func Line(colorAttr color.Attribute, text string) {
	color.New(colorAttr).Println("███ " + normalizeFinalNewLine(text))
}

func normalizeFinalNewLine(text string) string {
	result, _ := strings.CutSuffix(text, "\n")
	return result
}
