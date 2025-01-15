package config

import (
	"reflect"
	"testing"
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
		err := ValidateTaskReference(testCase.input)
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
		expected TaskReference
	}{
		{"default", TaskReference{
			Path:       []string{"default"},
			IsRelative: true,
			IsOptional: false,
		}},
		{"docker:bin", TaskReference{
			Path:       []string{"docker", "bin"},
			IsRelative: true,
			IsOptional: false,
		}},
		{"docker:package:default?", TaskReference{
			Path:       []string{"docker", "package", "default"},
			IsRelative: true,
			IsOptional: true,
		}},
		{":docker:package?", TaskReference{
			Path:       []string{"docker", "package"},
			IsRelative: false,
			IsOptional: true,
		}},
	}

	for _, testCase := range testCases {
		result := MustParseTaskReference(testCase.input)
		if !reflect.DeepEqual(testCase.expected, result) {
			t.Fatalf("Not deeply equal: \n%v\n%v", testCase.expected, result)
		}
	}
}

func TesTaskReferenceAbsoluteWorks(t *testing.T) {
	testCases := []struct {
		TaskReference TaskReference
		path          []string
		expected      TaskReference
	}{
		{TaskReference{
			Path:       []string{"default"},
			IsRelative: true,
			IsOptional: true,
		}, []string{}, TaskReference{
			Path:       []string{"default"},
			IsRelative: false,
			IsOptional: true,
		}},
		{TaskReference{
			Path:       []string{"default"},
			IsRelative: true,
			IsOptional: true,
		}, []string{"docker"}, TaskReference{
			Path:       []string{"docker", "default"},
			IsRelative: false,
			IsOptional: true,
		}},
		{TaskReference{
			Path:       []string{"default"},
			IsRelative: false,
			IsOptional: true,
		}, []string{"docker"}, TaskReference{
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
