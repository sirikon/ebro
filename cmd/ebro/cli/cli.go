package cli

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"slices"
	"strings"

	"github.com/sirikon/ebro/internal/logger"
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
	FlagFile  Flag = "file"
	FlagForce Flag = "force"
)

type commandInfo struct {
	Command     Command
	Description string
}

type flagInfo struct {
	Flag        Flag
	Description string
}

var commandList = []commandInfo{
	{CommandHelp, "Displays this help message"},
	{CommandVersion, "Display ebro's version"},
	{CommandConfig, "Display all imported configuration files merged into one"},
	{CommandCatalog, "Display complete catalog of tasks with their definitive configuration"},
	{CommandPlan, "Display the execution plan"},
	{CommandRun, ""},
}

var flagList = []flagInfo{
	{FlagFile, "Specify the file that should be loaded as root module. default: Ebro.yaml"},
	{FlagForce, "Ignore when.* conditionals and dont skip any task"},
}

type Arguments struct {
	Command Command
	File    string
	Targets []string
	Force   bool
}

var version = "dev"
var commandRe = regexp.MustCompile("^-([a-zA-Z0-9 ]+)$")
var flagRe = regexp.MustCompile("^--([a-zA-Z0-9 ]+)$")

func Parse() Arguments {
	result := Arguments{
		Command: CommandRun,
		File:    "Ebro.yaml",
		Targets: []string{":default"},
		Force:   false,
	}

	args := os.Args[1:]
	if len(args) == 0 {
		return result
	}

	if matches := commandRe.FindStringSubmatch(args[0]); matches != nil {
		args = args[1:]
		receivedCommand := Command(matches[1])

		if i := slices.IndexFunc(commandList, func(ci commandInfo) bool { return ci.Command == receivedCommand }); i == -1 {
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
			if i := slices.IndexFunc(flagList, func(ci flagInfo) bool { return ci.Flag == receivedFlag }); i == -1 {
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

			if receivedFlag == FlagForce {
				result.Force = true
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
	printAvailableCommands()
	fmt.Println()
	fmt.Println(strings.Trim(`
Available flags:
	`, " \n\t"))
	printAvailableFlags()
}

func printAvailableCommands() {
	padding := 0
	for _, info := range commandList {
		if len(string(info.Command)) > padding {
			padding = len(string(info.Command))
		}
	}

	for _, info := range commandList {
		if info.Description == "" {
			continue
		}
		fmt.Println("  -" + padRight(string(info.Command), padding) + "  " + info.Description)
	}
}

func printAvailableFlags() {
	padding := 0
	for _, info := range flagList {
		if len(string(info.Flag)) > padding {
			padding = len(string(info.Flag))
		}
	}

	for _, info := range flagList {
		if info.Description == "" {
			continue
		}
		fmt.Println("  --" + padRight(string(info.Flag), padding) + "  " + info.Description)
	}
}

func printVersion() {
	fmt.Println(version)
}

func ExitWithError(err error) {
	logger.Error(err.Error())
	os.Exit(1)
}

func padRight(text string, size int) string {
	for i := len(text); i < size; i++ {
		text = text + " "
	}
	return text
}
