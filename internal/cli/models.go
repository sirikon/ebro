package cli

import "reflect"

type Flag struct {
	Name        string
	Description string
	Kind        reflect.Kind
	Default     any
}

type FlagValue struct {
	Flag  *Flag
	Value any
}

type Command struct {
	Name           string
	Description    string
	Flags          []*Flag
	AcceptsTargets bool
}

type ExecutionArguments struct {
	Command *Command
	Flags   []FlagValue
	Targets []string
}

func (ea ExecutionArguments) GetFlagString(flag *Flag) *string {
	value := ea.getFlag(flag).(string)
	return &value
}

func (ea ExecutionArguments) GetFlagBool(flag *Flag) *bool {
	value := ea.getFlag(flag).(bool)
	return &value
}

func (ea ExecutionArguments) getFlag(flag *Flag) any {
	for _, flagValue := range ea.Flags {
		if flagValue.Flag == flag {
			return flagValue.Value
		}
	}
	return flag.Default
}

// type Command string
// type FlagKey string
// type FlagKind int
// type AcceptsTargets bool

// type Arguments struct {
// 	Command Command
// 	Flags   []Flag
// 	Targets []string
// }

// type commandInfo struct {
// 	Command        Command
// 	Description    string
// 	Flags          []Flag
// 	AcceptsTargets AcceptsTargets
// }

// type flagInfo struct {
// 	Flag        Flag
// 	Description string
// 	Kind        FlagKind
// }

// const (
// 	FlagKindBool   FlagKind = iota
// 	FlagKindString FlagKind = iota
// )

// const (
// 	AcceptsTargetsYes AcceptsTargets = true
// 	AcceptsTargetsNo  AcceptsTargets = false
// )

// const (
// 	CommandRun     Command = "run"
// 	CommandConfig  Command = "config"
// 	CommandCatalog Command = "catalog"
// 	CommandPlan    Command = "plan"
// 	CommandVersion Command = "version"
// 	CommandHelp    Command = "help"
// )

// var commandList = []commandInfo{
// 	{CommandRun, "Run everything", []Flag{FlagFile, FlagForce}, AcceptsTargetsYes},
// 	{CommandConfig, "Display all imported configuration files merged into one", []Flag{FlagFile}, AcceptsTargetsNo},
// 	{CommandCatalog, "Display complete catalog of tasks with their definitive configuration", []Flag{FlagFile}, AcceptsTargetsNo},
// 	{CommandPlan, "Display the execution plan", []Flag{FlagFile}, AcceptsTargetsYes},
// 	{CommandVersion, "Display ebro's version", []Flag{}, AcceptsTargetsNo},
// 	{CommandHelp, "Displays this help message", []Flag{}, AcceptsTargetsNo},
// }

// var flagList = []flagInfo{
// 	{FlagFile, "Specify the file that should be loaded as root module. default: Ebro.yaml", FlagKindString},
// 	{FlagForce, "Ignore when.* conditionals and dont skip any task", FlagKindBool},
// }
