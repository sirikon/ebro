package core2

import (
	"fmt"
	"regexp"
)

var nameValidCharsRe = `[a-zA-Z0-9-_\.]`
var NameRe = regexp.MustCompile("^" + nameValidCharsRe + "+$")

func ValidateName(name string) error {
	if !NameRe.MatchString(name) {
		return fmt.Errorf("name does not match the following regex: %v", NameRe.String())
	}
	return nil
}
