package core2

type Environment struct {
	Values []EnvironmentValue
}

type EnvironmentValue struct {
	Key   string
	Value string
}

func (env *Environment) Get(key string) *string {
	if env.Values == nil {
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
