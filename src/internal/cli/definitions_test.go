package cli

import (
	"reflect"
	"testing"
)

func TestFirstCommandIsDefault(t *testing.T) {
	if commands[0].Name != "" {
		t.Fatal("First command is not the default one")
	}
}

func TestOnlyOneCommandIsDefault(t *testing.T) {
	defaultCount := 0
	for _, command := range commands {
		if command.Name == "" {
			defaultCount += 1
		}
	}
	if defaultCount != 1 {
		t.Fatalf("Found %v default commands", defaultCount)
	}
}

func TestOnlyStringsAndBoolsAreAllowedInFlags(t *testing.T) {
	for _, command := range commands {
		for _, flag := range command.Flags {
			if flag.Kind != reflect.Bool && flag.Kind != reflect.String {
				t.Fatalf("Flag %v has invalid kind", flag.Name)
			}
		}
	}
}

func TestFlagDefaultValuesHaveCorrectType(t *testing.T) {
	for _, command := range commands {
		for _, flag := range command.Flags {
			if reflect.TypeOf(flag.Default).Kind() != flag.Kind {
				t.Fatalf("Flag %v has wrong default value type", flag.Name)
			}
		}
	}
}

func TestThereAreNoCommandOrFlagNameCollisions(t *testing.T) {
	commandNames := map[string]bool{}
	for _, command := range commands {
		if _, ok := commandNames[command.Name]; ok {
			t.Fatalf("Command name '%v' is repeated", command.Name)
		}
		commandNames[command.Name] = true

		if command.Name != "" {
			if _, ok := commandNames[command.ShortName]; ok {
				t.Fatalf("Command name '%v' is repeated", command.ShortName)
			}
			commandNames[command.ShortName] = true
		}

		flagNames := map[string]bool{}
		for _, flag := range command.Flags {
			if _, ok := flagNames[flag.Name]; ok {
				t.Fatalf("Flag name '%v' is repeated", flag.Name)
			}
			flagNames[flag.Name] = true
		}
	}
}
