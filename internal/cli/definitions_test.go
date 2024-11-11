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
