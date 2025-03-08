package core

import (
	"github.com/goccy/go-yaml"
)

type Environment struct {
	Values []EnvironmentValue
}

type EnvironmentValue struct {
	Key   string
	Value string
}

func (env *Environment) Get(key string) *string {
	if env == nil || env.Values == nil {
		return nil
	}

	for i := range env.Values {
		if env.Values[i].Key == key {
			value := env.Values[i].Value
			return &value
		}
	}

	return nil
}

func (env *Environment) Set(key, value string) {
	if env.Values == nil {
		env.Values = []EnvironmentValue{}
	}
	existingPos := -1
	for i := range env.Values {
		if env.Values[i].Key == key {
			existingPos = i
		}
	}
	if existingPos >= 0 {
		env.Values = append(env.Values[:existingPos], env.Values[existingPos+1:]...)
	}
	env.Values = append(env.Values, EnvironmentValue{
		Key:   key,
		Value: value,
	})
}

func (env *Environment) Map() map[string]string {
	result := map[string]string{}
	if env == nil || env.Values == nil {
		return result
	}
	for _, envVal := range env.Values {
		result[envVal.Key] = envVal.Value
	}
	return result
}

func (env *Environment) YamlMapSlice() yaml.MapSlice {
	if env == nil {
		return nil
	}
	result := yaml.MapSlice{}
	if env.Values == nil {
		return result
	}
	for _, envVal := range env.Values {
		result = append(result, yaml.MapItem{Key: envVal.Key, Value: envVal.Value})
	}
	return result
}

func (env *Environment) Clone() *Environment {
	if env == nil {
		return nil
	}
	result := &Environment{}
	if env.Values == nil {
		return result
	}
	for _, envVal := range env.Values {
		result.Values = append(result.Values, EnvironmentValue{Key: envVal.Key, Value: envVal.Value})
	}
	return result
}
