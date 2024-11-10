package cli

import "reflect"

var FlagFile = &Flag{
	Name:        "file",
	Description: "Specify the file that should be loaded as root module",
	Kind:        reflect.String,
	Default:     "Ebro.yaml",
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

var CommandConfig = &Command{
	Name:           "config",
	Description:    "Display all imported configuration files merged into one",
	Flags:          []*Flag{FlagFile},
	AcceptsTargets: false,
}

var CommandCatalog = &Command{
	Name:           "catalog",
	Description:    "Display complete catalog of tasks with their definitive configuration",
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
	Description: "Display ebro's version",
}

var CommandHelp = &Command{
	Name:        "help",
	Description: "Displays this help message",
}

var commands = []*Command{CommandRun, CommandConfig, CommandCatalog, CommandPlan, CommandVersion, CommandHelp}
