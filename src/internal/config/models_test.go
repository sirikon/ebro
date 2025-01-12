package config

import (
	"testing"
)

func TestModuleGetTaskWorks(t *testing.T) {
	module := Module{
		Tasks: map[string]*Task{
			"default": {Script: "default"},
		},
		Modules: map[string]*Module{
			"docker": {
				Tasks: map[string]*Task{
					"bin": {Script: "docker:bin"},
				},
			},
		},
	}

	testCases := []struct {
		input    string
		expected *Task
	}{
		{":default", module.Tasks["default"]},
		{":docker:bin", module.Modules["docker"].Tasks["bin"]},
		{":nonexistent", nil},
	}

	for _, testCase := range testCases {
		taskReference, err := ParseTaskReference(testCase.input)
		if err != nil {
			t.Fatal(err)
		}
		_, result := module.GetTask(taskReference)
		if result != testCase.expected {
			t.Fatalf("Not the same object: \n%v\n%v", testCase.expected, result)
		}
	}
}
