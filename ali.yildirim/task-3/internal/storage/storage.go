package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const (
	dirPerm  = 0o750
	filePerm = 0o600
)

func WriteJSON(path string, data any) error {
	encoded, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return fmt.Errorf("encode json: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(path), dirPerm); err != nil {
		return fmt.Errorf("create output dir: %w", err)
	}

	if err := os.WriteFile(path, encoded, filePerm); err != nil {
		return fmt.Errorf("write output file: %w", err)
	}

	return nil
}
