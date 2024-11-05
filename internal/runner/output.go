package runner

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path"
)

func storeTaskOutputAndCheckIfChanged(task_name string, output []byte) (bool, error) {
	changed := false

	newHash, err := hashBytes(output)
	if err != nil {
		return changed, fmt.Errorf("hashing output: %w", err)
	}

	outputPath := path.Join(".ebro", "output_tracking", task_name)
	err = os.MkdirAll(path.Dir(outputPath), os.ModePerm)
	if err != nil {
		return changed, fmt.Errorf("creating directory for output tracking: %w", err)
	}

	currentHashBytes, err := os.ReadFile(outputPath)
	if errors.Is(err, os.ErrNotExist) {
		changed = true
	} else {
		if err != nil {
			return changed, fmt.Errorf("reading output tracker for task %v: %w", task_name, err)
		}
		currentHash := string(currentHashBytes)

		changed = currentHash != newHash
	}

	if changed {
		err = os.WriteFile(outputPath, []byte(newHash), os.ModePerm)
		if err != nil {
			return changed, fmt.Errorf("writing output tracker for task %v: %w", task_name, err)
		}
	}

	return changed, nil
}

func hashBytes(data []byte) (string, error) {
	hasher := sha256.New()
	_, err := hasher.Write(data)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(hasher.Sum(nil)), nil
}
