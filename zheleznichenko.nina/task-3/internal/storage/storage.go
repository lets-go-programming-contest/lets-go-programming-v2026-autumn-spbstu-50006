package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const (
	outputDirPermissions  = 0o750
	outputFilePermissions = 0o600
)

func WriteResults(path string, results any) error {
	output, err := json.MarshalIndent(results, "", "    ")
	if err != nil {
		return fmt.Errorf("failed to marshal json: %w", err)
	}

	outputDir := filepath.Dir(path)
	if err := os.MkdirAll(outputDir, outputDirPermissions); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	if err := os.WriteFile(path, output, outputFilePermissions); err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}

	return nil
}