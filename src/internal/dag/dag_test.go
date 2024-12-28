package dag

import (
	"reflect"
	"testing"
)

func TestWorksWithSimpleScenario(t *testing.T) {
	testDag := NewDag()
	testDag.Link("A", "B")
	testDag.Link("B", "C")
	testDag.Link("Y", "Z")
	target := []string{"A"}

	expectedResult := []string{"C", "B", "A"}

	result, remains := testDag.Resolve(target)
	if remains != nil {
		t.Fatalf("There are remains: %v", remains)
	}
	if !reflect.DeepEqual(result, expectedResult) {
		t.Fatalf("Not deep equal: \n%v\n%v", result, expectedResult)
	}
}
