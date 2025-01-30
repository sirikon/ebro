package config

import (
	"os"
	"path"
	"reflect"
	"testing"

	"github.com/goccy/go-yaml"

	"github.com/sirikon/ebro/internal/core"
)

func TestParseTaskReferenceChecksRegex(t *testing.T) {
	testCases := []struct {
		input       string
		should_work bool
	}{
		{"default", true},
		{"docker:bin", true},
		{"docker:package:default?", true},
		{":docker:package?", true},
		{"", false},
		{":", false},
		{"::docker", false},
		{"docker??", false},
		{"docker!", false},
		{"DOCKER.thing", true},
		{":DOCKER.thing:other.thing?", true},
	}

	for _, testCase := range testCases {
		err := validateTaskReference(testCase.input)
		if testCase.should_work {
			if err != nil {
				t.Fatal(err)
			}
		} else {
			if err == nil || err.Error() != "task reference is invalid" {
				t.Fatal(err)
			}
		}
	}
}

func TestParseTaskReferenceWorks(t *testing.T) {
	testCases := []struct {
		input    string
		expected taskReference
	}{
		{"default", taskReference{
			Path:       []string{"default"},
			IsRelative: true,
			IsOptional: false,
		}},
		{"docker:bin", taskReference{
			Path:       []string{"docker", "bin"},
			IsRelative: true,
			IsOptional: false,
		}},
		{"docker:package:default?", taskReference{
			Path:       []string{"docker", "package", "default"},
			IsRelative: true,
			IsOptional: true,
		}},
		{":docker:package?", taskReference{
			Path:       []string{"docker", "package"},
			IsRelative: false,
			IsOptional: true,
		}},
	}

	for _, testCase := range testCases {
		result := mustParseTaskReference(testCase.input)
		if !reflect.DeepEqual(testCase.expected, result) {
			t.Fatalf("Not deeply equal: \n%v\n%v", testCase.expected, result)
		}
	}
}

func TesTaskReferenceAbsoluteWorks(t *testing.T) {
	testCases := []struct {
		TaskReference taskReference
		path          []string
		expected      taskReference
	}{
		{taskReference{
			Path:       []string{"default"},
			IsRelative: true,
			IsOptional: true,
		}, []string{}, taskReference{
			Path:       []string{"default"},
			IsRelative: false,
			IsOptional: true,
		}},
		{taskReference{
			Path:       []string{"default"},
			IsRelative: true,
			IsOptional: true,
		}, []string{"docker"}, taskReference{
			Path:       []string{"docker", "default"},
			IsRelative: false,
			IsOptional: true,
		}},
		{taskReference{
			Path:       []string{"default"},
			IsRelative: false,
			IsOptional: true,
		}, []string{"docker"}, taskReference{
			Path:       []string{"default"},
			IsRelative: false,
			IsOptional: true,
		}},
	}

	for _, testCase := range testCases {
		result := testCase.TaskReference.Absolute(testCase.path)
		if !reflect.DeepEqual(testCase.expected, result) {
			t.Fatalf("Not deeply equal: \n%v\n%v", testCase.expected, result)
		}
	}
}

func TestModuleMappingWorks(t *testing.T) {
	workingDirectory, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	baseDir := path.Join(workingDirectory, "../../../playground")

	rootModule, err := parseModule(path.Join(baseDir, "Ebro.yaml"))
	if err != nil {
		t.Fatal(err)
	}

	err = ValidateModule(rootModule)
	if err != nil {
		t.Fatal(err)
	}

	indexedModule := NewIndexedModule(rootModule)
	PurgeModule(indexedModule)
	ctx := ctxNormalizeModule{
		indexedModule: indexedModule,
	}
	err = ctx.normalizeModule(ctx.indexedModule.Module, []string{})
	if err != nil {
		t.Fatal(err)
	}

	// convert using yaml
	yamlResult := &core.Module{}
	data, err := yaml.Marshal(indexedModule.Module)
	if err != nil {
		t.Fatal(err)
	}
	yaml.Unmarshal(data, yamlResult)

	// convert using mapper
	mapperResult := MapToCoreModule(indexedModule.Module)

	// convert both to yaml to compare
	yamlResultData, err := yaml.Marshal(yamlResult)
	if err != nil {
		t.Fatal(err)
	}

	mapperResultData, err := yaml.Marshal(mapperResult)
	if err != nil {
		t.Fatal(err)
	}

	if string(yamlResultData) != string(mapperResultData) {
		t.Fatalf("Not deeply equal: \n%v\n%v", string(yamlResultData), string(mapperResultData))
	}
}
