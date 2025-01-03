package cli

import (
	"reflect"

	"github.com/sirikon/ebro/internal/constants"
)

var FlagFile = &Flag{
	Name:        "file",
	Description: "Specify the file that should be loaded as root module",
	Kind:        reflect.String,
	Default:     constants.DefaultFile,
}

var FlagForce = &Flag{
	Name:        "force",
	Description: "Ignore when.* conditionals and dont skip any task",
	Kind:        reflect.Bool,
	Default:     false,
}

var CommandRun = &Command{
	Name:           "",
	Description:    "Run everything",
	Flags:          []*Flag{FlagFile, FlagForce},
	AcceptsTargets: true,
}

var CommandInventory = &Command{
	Name:           "inventory",
	Description:    "Display complete inventory of tasks with their definitive configuration in YAML format",
	Flags:          []*Flag{FlagFile},
	AcceptsTargets: false,
}

var CommandPlan = &Command{
	Name:           "plan",
	Description:    "Display the execution plan",
	Flags:          []*Flag{FlagFile},
	AcceptsTargets: true,
}

var CommandVersion = &Command{
	Name:        "version",
	Description: "Display ebro's version information in YAML format",
}

var CommandHelp = &Command{
	Name:        "help",
	Description: "Display this help message",
}

var commands = []*Command{CommandRun, CommandInventory, CommandPlan, CommandVersion, CommandHelp}
