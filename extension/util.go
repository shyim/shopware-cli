package extension

import (
	"os"
	"path/filepath"
	"strings"
)

func PlatformPath(projectRoot, component, path string) string {
	if _, err := os.Stat(filepath.Join(projectRoot, "src", "Core", "composer.json")); err == nil {
		return filepath.Join(projectRoot, "src", component, path)
	} else if _, err := os.Stat(filepath.Join(projectRoot, "vendor", "shopware", "platform")); err == nil {
		return filepath.Join(projectRoot, "vendor", "shopware", "platform", "src", component, path)
	}

	return filepath.Join(projectRoot, "vendor", "shopware", strings.ToLower(component), path)
}

// IsContributeProject checks if the project is a contribution project aka shopware/shopware
func IsContributeProject(projectRoot string) bool {
	if _, err := os.Stat(filepath.Join(projectRoot, "src", "Core", "composer.json")); err == nil {
		return true
	}

	return false
}
