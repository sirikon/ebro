package utils

import (
	"reflect"
	"testing"
)

func TestWorksWithSimpleScenario(t *testing.T) {
	data := []string{"a", "b", "c", "b", "d"}
	expectedResult := []string{"a", "b", "c", "d"}
	result := Dedupe(data)
	if !reflect.DeepEqual(result, expectedResult) {
		t.Fatalf("Not deeply equal: \n%v\n%v", result, expectedResult)
	}
}
