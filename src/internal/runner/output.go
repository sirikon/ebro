package runner

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/sirikon/ebro/internal/core"
)

func storeTaskOutputAndCheckIfChanged(taskId core.TaskId, output []byte) (bool, error) {
	changed := false

	newHash, err := hashBytes(output)
	if err != nil {
		return changed, fmt.Errorf("hashing output: %w", err)
	}

	outputPath := path.Join(".ebro", "output_tracking", string(taskId))
	err = os.MkdirAll(path.Dir(outputPath), os.ModePerm)
	if err != nil {
		return changed, fmt.Errorf("creating directory for output tracking: %w", err)
	}

	currentHashBytes, err := os.ReadFile(outputPath)
	if errors.Is(err, os.ErrNotExist) {
		changed = true
	} else {
		if err != nil {
			return changed, fmt.Errorf("reading output tracker: %w", err)
		}
		currentHash := string(currentHashBytes)

		changed = currentHash != newHash
	}

	if changed {
		err = os.WriteFile(outputPath, []byte(newHash), os.ModePerm)
		if err != nil {
			return changed, fmt.Errorf("writing output tracker: %w", err)
		}
	}

	return changed, nil
}

func removeTaskOutput(taskId core.TaskId) error {
	outputPath := path.Join(".ebro", "output_tracking", string(taskId))
	err := os.Remove(outputPath)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("removing task output: %w", err)
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
