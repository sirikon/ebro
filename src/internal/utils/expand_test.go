package utils

import (
	"reflect"
	"testing"

	"github.com/sirikon/ebro/internal/core"
)

func TestExpandMergeEnvsWorksWithSimpleScenario(t *testing.T) {
	childEnv := core.NewEnvironment(
		core.EnvironmentValue{Key: "A", Value: "1"},
		core.EnvironmentValue{Key: "B", Value: "22"},
	)
	parentEnv := core.NewEnvironment(
		core.EnvironmentValue{Key: "B", Value: "2"},
	)
	grandparentEnv := core.NewEnvironment(
		core.EnvironmentValue{Key: "C", Value: "3"},
	)

	expectedResult := core.NewEnvironment(
		core.EnvironmentValue{Key: "C", Value: "3"},
		core.EnvironmentValue{Key: "A", Value: "1"},
		core.EnvironmentValue{Key: "B", Value: "22"},
	)

	result, err := ExpandMergeEnvs(childEnv, parentEnv, grandparentEnv)
	if err != nil {
		t.Fatalf("Error during execution: %v", err)
	}
	if !reflect.DeepEqual(result, expectedResult) {
		t.Fatalf("Not deeply equal: \n%v\n%v", result, expectedResult)
	}
}

func TestExpandMergeEnvsWorksWithComplexScenario(t *testing.T) {
	childEnv := core.NewEnvironment(
		core.EnvironmentValue{Key: "DOCKER_VERSION", Value: "3.0.0"},
		core.EnvironmentValue{Key: "DOCKER_VERSION_FOR_INSTALL", Value: "${DOCKER_VERSION}-${ACTIVE_PACKAGE_MANAGER}"},
	)
	parentEnv := core.NewEnvironment(
		core.EnvironmentValue{Key: "DOCKER_VERSION", Value: "2.0.0"},
	)
	grandparentEnv := core.NewEnvironment(
		core.EnvironmentValue{Key: "ACTIVE_PACKAGE_MANAGER", Value: "apt"},
	)

	// This is the correct behavior. DOCKER_VERSION is redefined at the same
	// environment it's being used, so it will end up in the final environment,
	// but the interpolation values in DOCKER_VERSION_FOR_INSTALL will come
	// from the parent environment exclusively
	expectedResult := core.NewEnvironment(
		core.EnvironmentValue{Key: "ACTIVE_PACKAGE_MANAGER", Value: "apt"},
		core.EnvironmentValue{Key: "DOCKER_VERSION", Value: "3.0.0"},
		core.EnvironmentValue{Key: "DOCKER_VERSION_FOR_INSTALL", Value: "2.0.0-apt"},
	)

	result, err := ExpandMergeEnvs(childEnv, parentEnv, grandparentEnv)
	if err != nil {
		t.Fatalf("Error during execution: %v", err)
	}
	if !reflect.DeepEqual(result, expectedResult) {
		t.Fatalf("Not deeply equal: \n%v\n%v", result, expectedResult)
	}
}
