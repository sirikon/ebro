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
		_, err := ParseTaskReference(testCase.input)
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
			Parts:      []string{"default"},
			IsRelative: true,
			IsOptional: false,
		}},
		{"docker:bin", TaskReference{
			Parts:      []string{"docker", "bin"},
			IsRelative: true,
			IsOptional: false,
		}},
		{"docker:package:default?", TaskReference{
			Parts:      []string{"docker", "package", "default"},
			IsRelative: true,
			IsOptional: true,
		}},
		{":docker:package?", TaskReference{
			Parts:      []string{"docker", "package"},
			IsRelative: false,
			IsOptional: true,
		}},
	}

	for _, testCase := range testCases {
		result, err := ParseTaskReference(testCase.input)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(testCase.expected, result) {
			t.Fatalf("Not deeply equal: \n%v\n%v", testCase.expected, result)
		}
	}
}

func TesTaskReferenceAbsoluteWorks(t *testing.T) {
	testCases := []struct {
		TaskReference TaskReference
		parts         []string
		expected      TaskReference
	}{
		{TaskReference{
			Parts:      []string{"default"},
			IsRelative: true,
			IsOptional: true,
		}, []string{}, TaskReference{
			Parts:      []string{"default"},
			IsRelative: false,
			IsOptional: true,
		}},
		{TaskReference{
			Parts:      []string{"default"},
			IsRelative: true,
			IsOptional: true,
		}, []string{"docker"}, TaskReference{
			Parts:      []string{"docker", "default"},
			IsRelative: false,
			IsOptional: true,
		}},
		{TaskReference{
			Parts:      []string{"default"},
			IsRelative: false,
			IsOptional: true,
		}, []string{"docker"}, TaskReference{
			Parts:      []string{"default"},
			IsRelative: false,
			IsOptional: true,
		}},
	}

	for _, testCase := range testCases {
		result := testCase.TaskReference.Absolute(testCase.parts)
		if !reflect.DeepEqual(testCase.expected, result) {
			t.Fatalf("Not deeply equal: \n%v\n%v", testCase.expected, result)
		}
	}
}
