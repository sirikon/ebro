package runner

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strings"
)

func checkSourcesChanged(task_name string, sources []string) (bool, error) {
	sourceTrackerPath := getSourceTrackerPath(task_name)
	_, err := os.Stat(sourceTrackerPath)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return false, fmt.Errorf("checking if source tracker exists for task %v: %w", task_name, err)
	}

	sourceTrackerExists := err == nil
	if len(sources) == 0 && !sourceTrackerExists {
		return false, nil
	}
	if !sourceTrackerExists {
		return true, nil
	}

	data, err := os.ReadFile(sourceTrackerPath)
	if err != nil {
		return false, fmt.Errorf("reading source tracker for task %v: %w", task_name, err)
	}
	hashes := strings.Split(strings.Trim(string(data), " \n"), "\n")
	if len(hashes) != len(sources) {
		return true, nil
	}

	for i, source := range sources {
		hash, err := hashFile(source)
		if err != nil {
			return false, fmt.Errorf("checking if source %v changed: %w", source, err)
		}
		if hash != hashes[i] {
			return true, nil
		}
	}
	return false, nil
}

func hashSources(task_name string, sources []string) error {
	hashes := []string{}
	for _, source := range sources {
		hash, err := hashFile(source)
		if err != nil {
			return fmt.Errorf("hashing source %v for task %v: %w", source, task_name, err)
		}
		hashes = append(hashes, hash)
	}
	sourceTrackerPath := getSourceTrackerPath(task_name)
	err := os.MkdirAll(path.Dir(sourceTrackerPath), os.ModePerm)
	if err != nil {
		return fmt.Errorf("creating source tracker directory for task %v: %w", task_name, err)
	}
	err = os.WriteFile(sourceTrackerPath, []byte(strings.Join(hashes, "\n")), os.ModePerm)
	if err != nil {
		return fmt.Errorf("writing source tracker for task %v: %w", task_name, err)
	}
	return nil
}

func getSourceTrackerPath(task_name string) string {
	return path.Join(".ebro", "source_trackers", task_name)
}

func hashFile(filePath string) (string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("hashing file %v: %w", filePath, err)
	}
	defer f.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, f); err != nil {
		log.Fatal(err)
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}
