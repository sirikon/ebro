package utils

// func TestExpandMergeEnvsWorksWithSimpleScenario(t *testing.T) {
// 	childEnv := map[string]string{
// 		"A": "1",
// 		"B": "22",
// 	}
// 	parentEnv := map[string]string{
// 		"B": "2",
// 	}
// 	grandparentEnv := map[string]string{
// 		"C": "3",
// 	}

// 	expectedResult := map[string]string{
// 		"A": "1",
// 		"B": "22",
// 		"C": "3",
// 	}

// 	result, err := ExpandMergeEnvs(childEnv, parentEnv, grandparentEnv)
// 	if err != nil {
// 		t.Fatalf("Error during execution: %v", err)
// 	}
// 	if !reflect.DeepEqual(result, expectedResult) {
// 		t.Fatalf("Not deeply equal: \n%v\n%v", result, expectedResult)
// 	}
// }

// func TestExpandMergeEnvsWorksWithComplexScenario(t *testing.T) {
// 	childEnv := map[string]string{
// 		"DOCKER_VERSION":             "3.0.0",
// 		"DOCKER_VERSION_FOR_INSTALL": "${DOCKER_VERSION}-${ACTIVE_PACKAGE_MANAGER}",
// 	}
// 	parentEnv := map[string]string{
// 		"DOCKER_VERSION": "2.0.0",
// 	}
// 	grandparentEnv := map[string]string{
// 		"ACTIVE_PACKAGE_MANAGER": "apt",
// 	}

// 	// This is the correct behavior. DOCKER_VERSION is redefined at the same
// 	// environment it's being used, so it will end up in the final environment,
// 	// but the interpolation values in DOCKER_VERSION_FOR_INSTALL will come
// 	// from the parent environment exclusively
// 	expectedResult := map[string]string{
// 		"DOCKER_VERSION":             "3.0.0",
// 		"ACTIVE_PACKAGE_MANAGER":     "apt",
// 		"DOCKER_VERSION_FOR_INSTALL": "2.0.0-apt",
// 	}

// 	result, err := ExpandMergeEnvs(childEnv, parentEnv, grandparentEnv)
// 	if err != nil {
// 		t.Fatalf("Error during execution: %v", err)
// 	}
// 	if !reflect.DeepEqual(result, expectedResult) {
// 		t.Fatalf("Not deeply equal: \n%v\n%v", result, expectedResult)
// 	}
// }
