package config

import "testing"

func TestIndexedModuleFindTaskWorks(t *testing.T) {
	indexedModule := NewIndexedModule(&Module{
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
		{":default", indexedModule.Module.Tasks["default"]},
		{":docker:bin", indexedModule.Module.Modules["docker"].Tasks["bin"]},
		{":nonexistent", nil},
	}

	for _, testCase := range testCases {
		taskReference := mustParseTaskReference(testCase.input)
		_, result := FindTask(indexedModule, taskReference)
		if result != testCase.expected {
			t.Fatalf("Not the same object: \n%v\n%v", testCase.expected, result)
		}
	}
}
