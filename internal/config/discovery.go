package config

func DiscoverConfig() (*Module, error) {
	return parseModuleFromFile("Ebro.yaml")
}
