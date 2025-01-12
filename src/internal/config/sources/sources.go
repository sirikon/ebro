package sources

import (
	// "github.com/sirikon/ebro/internal/config/sources/git"
	"github.com/sirikon/ebro/internal/config/sources/local"
)

type Source interface {
	Match(from string) (bool, error)
	Resolve(base string, from string) (string, error)
}

var Sources = []Source{
	// git.GitSource{},
	local.LocalSource{},
}
