package cli

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

var version = "dev"

type Arguments struct {
	Flags   Flags
	Targets []string
}

type Flags struct {
	Config  bool `flag:"config" doc:"0|Display all imported configuration files merged into one"`
	Catalog bool `flag:"catalog" doc:"1|Display complete catalog of tasks with their definitive configuration"`
	Plan    bool `flag:"plan" doc:"2|Display the execution plan"`
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
	flagsWithDoc := make(map[int][]string)
	flagsOrders := []int{}
	flagLength := 0
	for i := 0; i < flagsType.NumField(); i++ {
		field := flagsType.Field(i)
		flag := field.Tag.Get("flag")
		doc_parts := strings.Split(field.Tag.Get("doc"), "|")
		order, err := strconv.Atoi(doc_parts[0])
		if err != nil {
			ExitWithError(err)
		}
		doc := doc_parts[1]
		flagsWithDoc[order] = []string{flag, doc}
		flagsOrders = append(flagsOrders, order)
		if len(flag) > flagLength {
			flagLength = len(flag)
		}
	}

	slices.Sort(flagsOrders)

	for _, order := range flagsOrders {
		data := flagsWithDoc[order]
		flag := data[0]
		doc := data[1]
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
	if color.NoColor {
		fmt.Print("ERROR:")
	} else {
		color.New(color.BgRed).Add(color.FgWhite).Print(" ERROR ")
	}
	fmt.Print(" ")
	fmt.Println(err)
	os.Exit(1)
}
