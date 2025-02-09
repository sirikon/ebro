package utils

import (
	"fmt"
	"os"
	"strings"

	"github.com/sirikon/ebro/internal/core"
	"mvdan.cc/sh/v3/expand"
	"mvdan.cc/sh/v3/syntax"
)

func ExpandMergeEnvs(envs ...*core.Environment) (*core.Environment, error) {
	result := &core.Environment{}
	for i := (len(envs) - 1); i >= 0; i-- {
		env := envs[i]
		if env == nil {
			continue
		}

		for _, envValue := range env.Values {
			expandedValue, err := ExpandString(envValue.Value, result)
			if err != nil {
				return nil, fmt.Errorf("expanding %v: %w", envValue.Value, err)
			}
			result.Set(envValue.Key, expandedValue)
		}
	}
	return result, nil
}

func ExpandString(s string, env *core.Environment) (string, error) {
	// This function is a re-implementation of shell.Expand just to be able to
	// change the default configuration (NoUnset).

	p := syntax.NewParser()
	word, err := p.Document(strings.NewReader(s))
	if err != nil {
		return "", err
	}

	cfg := &expand.Config{
		NoUnset: true,
		Env: expand.FuncEnviron(func(s string) string {
			if val := env.Get(s); val != nil {
				return *val
			}
			return os.Getenv(s)
		}),
	}

	return expand.Document(cfg, word)
}
