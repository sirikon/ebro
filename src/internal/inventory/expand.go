package inventory

import (
	"fmt"
	"maps"
	"os"
	"slices"

	"mvdan.cc/sh/v3/shell"
)

func expandMergeEnvs(envs ...map[string]string) (map[string]string, error) {
	result := map[string]string{}
	for i := (len(envs) - 1); i >= 0; i-- {
		parentEnv := maps.Clone(result)
		env := envs[i]
		// We want to iterate through keys in a repeatable and predictable way.
		// The order in which we process each key SHOULD NOT BE IMPORTANT, but
		// in the scenario of a bug in here, we want the behavior to be
		// consistent.
		//
		// That's why we're sorting the keys and iterating over them
		// instead of `range`ing the map directly.
		envKeys := slices.Sorted(maps.Keys(env))
		for _, key := range envKeys {
			expandedValue, err := expandString(env[key], parentEnv)
			if err != nil {
				return nil, fmt.Errorf("expanding %v: %w", env[key], err)
			}
			result[key] = expandedValue
		}
	}
	return result, nil
}

func expandString(s string, env map[string]string) (string, error) {
	return shell.Expand(s, func(s string) string {
		if val, ok := env[s]; ok {
			return val
		}
		return os.Getenv(s)
	})
}
