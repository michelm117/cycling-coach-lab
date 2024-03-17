package utils

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
)

func GetProjectRoot() string {
	_, currentFilePath, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatalf("failed to get path")
	}

	if os.Getenv("ENVIRONMENT") != "production" {
		return filepath.Join(filepath.Dir(currentFilePath), "../")
	}
	return filepath.Join(filepath.Dir(currentFilePath), ".")
}
