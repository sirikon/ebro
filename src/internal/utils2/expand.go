package utils2

import (
	"fmt"
	"os"

	"github.com/sirikon/ebro/internal/core2"
	"mvdan.cc/sh/v3/shell"
)

func ExpandMergeEnvs(envs ...*core2.Environment) (*core2.Environment, error) {
	result := &core2.Environment{}
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

func ExpandString(s string, env *core2.Environment) (string, error) {
	return shell.Expand(s, func(s string) string {
		if val := env.Get(s); val != nil {
			return *val
		}
		return os.Getenv(s)
	})
}
