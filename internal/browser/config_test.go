package browser_test

import (
	"testing"

	"go-rod-testing-browser-restrict/internal/browser"
)

// Test untuk memastikan config default memiliki URL yang valid
func TestDefaultConfig(t *testing.T) {
	config := browser.DefaultConfig()

	if config.DownloadURL == "" {
		t.Error("DownloadURL should not be empty")
	}

	if config.InstallDirName == "" {
		t.Error("InstallDirName should not be empty")
	}

	if config.Version == "" {
		t.Error("Version should not be empty")
	}

	// Pastikan URL menggunakan v142.7444.229 bukan v142.0.7444.229
	expectedTag := "v142.7444.229"
	if !contains(config.DownloadURL, expectedTag) {
		t.Errorf("DownloadURL should contain %s, got: %s", expectedTag, config.DownloadURL)
	}
}

func contains(str, substr string) bool {
	return len(str) >= len(substr) && (str == substr ||
		(len(str) > len(substr) && containsHelper(str, substr)))
}

func containsHelper(str, substr string) bool {
	for i := 0; i <= len(str)-len(substr); i++ {
		if str[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
