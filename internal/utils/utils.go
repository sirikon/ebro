package utils

import "maps"

func MergeEnv(envs ...map[string]string) map[string]string {
	result := make(map[string]string)
	for _, env := range envs {
		maps.Copy(result, env)
	}
	return result
}
