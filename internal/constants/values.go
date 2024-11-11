package constants

var version = "dev"
var commit = "HEAD"
var timestamp = "0"

func GetVersion() string {
	return version
}

func GetCommit() string {
	return commit
}

func GetTimestamp() string {
	return timestamp
}

const DefaultTarget = "default"
const DefaultFile = "Ebro.yaml"
