package cli

import (
	"fmt"
	"strings"
)

func PrintVersion() {
	fmt.Println(version)
}

func PrintHelp() {
	fmt.Println(strings.Trim(`
Usage: ebro [-command?] [--flags?...] [targets?...]

Available commands:
`, " \n\t"))
	printAvailableCommands()
}

func printAvailableCommands() {
	padding := 0
	for _, command := range commands {
		if len(command.Name) > padding {
			padding = len(command.Name)
		}
	}

	for _, command := range commands {
		if command.Name == "" {
			continue
		}
		fmt.Println("  -" + padRight(command.Name, padding) + "  " + command.Description)
	}
}

func padRight(text string, size int) string {
	for i := len(text); i < size; i++ {
		text = text + " "
	}
	return text
}
