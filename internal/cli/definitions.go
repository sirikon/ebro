package cli

import "reflect"

const DefaultTarget = "default"

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
	Description: "Display this help message",
}

var commands = []*Command{CommandRun, CommandCatalog, CommandPlan, CommandVersion, CommandHelp}
