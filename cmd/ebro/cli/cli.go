package cli

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strings"

	"github.com/fatih/color"
)

var version = "dev"

type Arguments struct {
	Flags   Flags
	Targets []string
}

type Flags struct {
	Config  bool `flag:"config" doc:"Display all imported configuration files merged into one"`
	Catalog bool `flag:"catalog" doc:"Display complete catalog of tasks with their definitive configuration"`
	Plan    bool `flag:"plan" doc:"Display the execution plan"`
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

		if receivedFlag == "version" {
			printVersion()
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
			ExitWithError(errors.New("unknown flag: " + receivedFlag))
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
Usage: ebro [flag?] [targets?...]

Available flags:
`, " \n\t"))
	flagsType := reflect.TypeOf(Flags{})
	flagsWithDoc := make(map[string]string)
	flagLength := 0
	for i := 0; i < flagsType.NumField(); i++ {
		field := flagsType.Field(i)
		flag := field.Tag.Get("flag")
		flagsWithDoc[flag] = field.Tag.Get("doc")
		if len(flag) > flagLength {
			flagLength = len(flag)
		}
	}

	for flag, doc := range flagsWithDoc {
		fmt.Print("  -" + flag)
		for i := len(flag); i < (flagLength + 2); i++ {
			fmt.Print(" ")
		}
		fmt.Println(doc)
	}
}

func printVersion() {
	fmt.Println(version)
}

func ExitWithError(err error) {
	color.New(color.BgRed).Add(color.FgWhite).Print(" ERROR ")
	fmt.Print(" ")
	fmt.Println(err)
	os.Exit(1)
}
