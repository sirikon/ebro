package cli

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"slices"

	"github.com/sirikon/ebro/internal/logger"
)

var version = "dev"
var commandRe = regexp.MustCompile("^-([a-zA-Z0-9 ]+)$")
var flagRe = regexp.MustCompile("^--([a-zA-Z0-9 ]+)$")

func Parse() ExecutionArguments {
	result := ExecutionArguments{
		Command: CommandRun,
		Targets: []string{":default"},
	}

	args := os.Args[1:]
	if len(args) == 0 {
		return result
	}

	if matches := commandRe.FindStringSubmatch(args[0]); matches != nil {
		args = args[1:]
		receivedName := matches[1]
		i := slices.IndexFunc(commands, func(c *Command) bool { return c.Name == receivedName })
		if i == -1 {
			ExitWithError(errors.New("unknown command: " + receivedName))
		}
		result.Command = commands[i]
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
			receivedName := matches[1]

			i := slices.IndexFunc(result.Command.Flags, func(f *Flag) bool { return f.Name == receivedName })
			if i == -1 {
				ExitWithError(errors.New("unknown flag: " + receivedName))
			}
			flag := result.Command.Flags[i]

			if flag.Kind == reflect.String {
				if len(args) >= 1 {
					value := args[0]
					args = args[1:]
					result.Flags = append(result.Flags, FlagValue{Flag: flag, Value: value})
				} else {
					ExitWithError(fmt.Errorf("expected value after --file flag"))
				}
			} else if flag.Kind == reflect.Bool {
				result.Flags = append(result.Flags, FlagValue{Flag: flag, Value: true})
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

func ExitWithError(err error) {
	logger.Error(err.Error())
	os.Exit(1)
}
