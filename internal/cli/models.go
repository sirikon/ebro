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
