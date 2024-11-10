package cli

import (
	"fmt"
	"reflect"

	"github.com/fatih/color"
	"github.com/sirikon/ebro/internal/build"
)

func PrintVersion() {
	fmt.Println(build.Version)
}

func PrintHelp() {
	printCommands()
}

func printCommands() {
	for _, command := range commands {
		fmt.Println()
		fmt.Print(color.CyanString("  ebro"))
		if command.Name != "" {
			fmt.Print(" -" + color.GreenString(command.Name))
		}
		if len(command.Flags) > 0 {
			fmt.Print(" [--" + color.YellowString("flags") + "...]")
		}
		if command.AcceptsTargets {
			fmt.Print(" [" + color.MagentaString("targets") + "...]")
		}
		fmt.Println()
		fmt.Print(color.HiBlackString("    # "))
		fmt.Print(command.Description)
		fmt.Println()
		if len(command.Flags) > 0 {
			fmt.Println("    flags:")
			printFlags(command.Flags)
		}
		if command.AcceptsTargets {
			fmt.Println("    targets:")
			fmt.Println("      defaults to [" + color.MagentaString(DefaultTarget) + "]")
		}
		fmt.Println()
	}
}

func printFlags(flags []*Flag) {
	padding := 0
	for _, flag := range flags {
		l := len(flag.Name)
		if flag.Kind == reflect.String {
			l = l + len(" value")
		}
		if l > padding {
			padding = l
		}
	}

	for _, flag := range flags {
		fmt.Print("      --")
		if flag.Kind == reflect.String {
			fmt.Print(color.YellowString(padRight(flag.Name+" value", padding)))
		} else {
			fmt.Print(color.YellowString(padRight(flag.Name, padding)))
		}
		fmt.Print("  ")
		fmt.Print(flag.Description)
		fmt.Print(". default: ")
		fmt.Println(flag.Default)
	}
}

func padRight(text string, size int) string {
	for i := len(text); i < size; i++ {
		text = text + " "
	}
	return text
}
