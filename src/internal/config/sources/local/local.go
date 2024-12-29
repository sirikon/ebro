package local

import "path"

type LocalSource struct{}

func (ls LocalSource) Match(from string) (bool, error) {
	return true, nil
}

func (ls LocalSource) Resolve(base string, from string) (string, error) {
	return path.Join(base, from), nil
}
