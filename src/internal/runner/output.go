package runner

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path"
)

func storeTaskOutputAndCheckIfChanged(taskName string, output []byte) (bool, error) {
	changed := false

	newHash, err := hashBytes(output)
	if err != nil {
		return changed, fmt.Errorf("hashing output: %w", err)
	}

	outputPath := path.Join(".ebro", "output_tracking", taskName)
	err = os.MkdirAll(path.Dir(outputPath), os.ModePerm)
	if err != nil {
		return changed, fmt.Errorf("creating directory for output tracking: %w", err)
	}

	currentHashBytes, err := os.ReadFile(outputPath)
	if errors.Is(err, os.ErrNotExist) {
		changed = true
	} else {
		if err != nil {
			return changed, fmt.Errorf("reading output tracker for task %v: %w", taskName, err)
		}
		currentHash := string(currentHashBytes)

		changed = currentHash != newHash
	}

	if changed {
		err = os.WriteFile(outputPath, []byte(newHash), os.ModePerm)
		if err != nil {
			return changed, fmt.Errorf("writing output tracker for task %v: %w", taskName, err)
		}
	}

	return changed, nil
}

func removeTaskOutput(taskName string) error {
	outputPath := path.Join(".ebro", "output_tracking", taskName)
	err := os.Remove(outputPath)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("removing %v task output: %w", taskName, err)
	}
	return nil
}

func hashBytes(data []byte) (string, error) {
	hasher := sha256.New()
	_, err := hasher.Write(data)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(hasher.Sum(nil)), nil
}