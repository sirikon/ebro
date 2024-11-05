package cli

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"slices"
	"strings"

	"github.com/fatih/color"
)

type Command string
type Flag string

const (
	CommandHelp    Command = "help"
	CommandVersion Command = "version"
	CommandConfig  Command = "config"
	CommandCatalog Command = "catalog"
	CommandPlan    Command = "plan"
	CommandRun     Command = "run"
)

const (
	FlagFile Flag = "file"
)

var commands = []Command{CommandHelp, CommandVersion, CommandConfig, CommandCatalog, CommandPlan, CommandRun}
var flags = []Flag{FlagFile}

type Arguments struct {
	Command Command
	File    string
	Targets []string
}

var version = "dev"
var commandRe = regexp.MustCompile("^-([a-zA-Z0-9 ]+)$")
var flagRe = regexp.MustCompile("^--([a-zA-Z0-9 ]+)$")

func Parse() Arguments {
	result := Arguments{
		Command: CommandRun,
		File:    "Ebro.yaml",
		Targets: []string{":default"},
	}

	args := os.Args[1:]
	if len(args) == 0 {
		return result
	}

	if matches := commandRe.FindStringSubmatch(args[0]); matches != nil {
		args = args[1:]
		receivedCommand := Command(matches[1])

		if i := slices.Index(commands, receivedCommand); i == -1 {
			ExitWithError(errors.New("unknown command: " + string(receivedCommand)))
		}

		if receivedCommand == CommandHelp {
			printHelp()
			os.Exit(0)
		}

		if receivedCommand == CommandVersion {
			printVersion()
			os.Exit(0)
		}

		result.Command = receivedCommand
	}

	scanFlags := true
	for scanFlags {
		scanFlags = false
		if len(args) == 0 {
			continue
		}
		if matches := flagRe.FindStringSubmatch(args[0]); matches != nil {
			scanFlags = true
			args = args[1:]
			receivedFlag := Flag(matches[1])
			if i := slices.Index(flags, receivedFlag); i == -1 {
				ExitWithError(errors.New("unknown flag: " + string(receivedFlag)))
			}

			if receivedFlag == FlagFile {
				if len(args) >= 1 {
					result.File = args[0]
					args = args[1:]
				} else {
					ExitWithError(fmt.Errorf("expected value after --file flag"))
				}
			}
		}
	}

	if len(args) > 0 {
		result.Targets = []string{}
		for _, arg := range args {
			result.Targets = append(result.Targets, ":"+arg)
		}
	}

	return result
}

func printHelp() {
	fmt.Println(strings.Trim(`
Usage: ebro [-command?] [--flags?...] [targets?...]

Available commands:
`, " \n\t"))
	for _, command := range commands {
		if command == "run" {
			continue
		}
		fmt.Println("  -" + command)
	}
	fmt.Println()
	fmt.Println(strings.Trim(`
Available flags:
	`, " \n\t"))
	for _, flag := range flags {
		fmt.Println("  --" + flag)
	}
}

func printVersion() {
	fmt.Println(version)
}

func ExitWithError(err error) {
	if color.NoColor {
		fmt.Print("ERROR:")
	} else {
		color.New(color.BgRed).Add(color.FgWhite).Print(" ERROR ")
	}
	fmt.Print(" ")
	if strings.HasSuffix(err.Error(), "\n") {
		fmt.Print(err)
	} else {
		fmt.Println(err)
	}
	os.Exit(1)
}
