package utils

import (
	"reflect"
	"testing"
)

func TestSetWorksWithSimpleScenario(t *testing.T) {
	set := NewSet[string]()
	set.Add("a", "b", "b")
	set.Add("c", "c", "d")
	if set.Length() != 4 {
		t.Fatalf("wrong length")
	}
	if !reflect.DeepEqual(set.List(), []string{"a", "b", "c", "d"}) {
		t.Fatalf("wrong list")
	}
	set.Delete("a", "d")
	if set.Length() != 2 {
		t.Fatalf("wrong length")
	}
	if !reflect.DeepEqual(set.List(), []string{"b", "c"}) {
		t.Fatalf("wrong list")
	}
}
