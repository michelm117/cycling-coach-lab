package utils_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/michelm117/cycling-coach-lab/utils"
)

func TestGetProjectRootForDevVsProd(t *testing.T) {
	os.Setenv("ENVIRONMENT", "development")
	devPath := utils.GetProjectRoot()
	os.Setenv("ENVIRONMENT", "production")
	prodPath := utils.GetProjectRoot()

	relativePath, err := filepath.Rel(devPath, prodPath)
	if err != nil {
		t.Errorf("Error comparing paths: %s", err)
	}

	if relativePath != "src/utils" &&
		len(strings.Split(devPath, "/")) == len(strings.Split(prodPath, "/"))+2 {
		t.Errorf("Expected %s to be a subdirectory of %s", prodPath, devPath)
	}
}
