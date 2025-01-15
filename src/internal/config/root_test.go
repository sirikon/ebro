package config

import "testing"

func TestModuleGetTaskWorks(t *testing.T) {
	rootModule := NewRootModule(&Module{
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
	})

	testCases := []struct {
		input    string
		expected *Task
	}{
		{":default", rootModule.Module.Tasks["default"]},
		{":docker:bin", rootModule.Module.Modules["docker"].Tasks["bin"]},
		{":nonexistent", nil},
	}

	for _, testCase := range testCases {
		taskReference := MustParseTaskReference(testCase.input)
		_, result := FindTask(rootModule, taskReference)
		if result != testCase.expected {
			t.Fatalf("Not the same object: \n%v\n%v", testCase.expected, result)
		}
	}
}
