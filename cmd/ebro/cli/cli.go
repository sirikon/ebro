package cli

import (
	"fmt"
	"os"
	"reflect"
	"regexp"
)

type Arguments struct {
	Flags   Flags
	Targets []string
}

type Flags struct {
	Config  bool `flag:"config"`
	Catalog bool `flag:"catalog"`
	Plan    bool `flag:"plan"`
}

var flagRe = regexp.MustCompile("^-{1,2}([a-zA-Z0-9 ]+)$")

func Parse() Arguments {
	result := Arguments{
		Flags: Flags{
			Config:  false,
			Catalog: false,
			Plan:    false,
		},
		Targets: []string{":default"},
	}
	args := os.Args[1:]
	if len(args) == 0 {
		return result
	}

	if matches := flagRe.FindStringSubmatch(args[0]); matches != nil {
		args = args[1:]
		receivedFlag := matches[1]

		if receivedFlag == "help" {
			printHelp()
			os.Exit(0)
		}

		flagsType := reflect.TypeOf(result.Flags)
		found := false
		for i := 0; i < flagsType.NumField(); i++ {
			field := flagsType.Field(i)
			if receivedFlag == field.Tag.Get("flag") {
				reflect.ValueOf(&result.Flags).Elem().FieldByName(field.Name).SetBool(true)
				found = true
				break
			}
		}
		if !found {
			fmt.Println("Unknown flag " + receivedFlag)
			os.Exit(1)
		}
	}

	if len(args) > 0 {
		result.Targets = args
	}

	return result
}

func printHelp() {
	fmt.Println("bruh help")
}
