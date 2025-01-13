package local

import "path"

type LocalSource struct{}

func (ls LocalSource) Match(from string) (bool, error) {
	return true, nil
}

func (ls LocalSource) Resolve(base string, from string) (string, error) {
	if path.IsAbs(from) {
		return from, nil
	}
	return path.Join(base, from), nil
}
