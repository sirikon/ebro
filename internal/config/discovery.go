package config

func DiscoverModule() (*Module, error) {
	return parseModuleFromFile("Ebro.yaml")
}
