package utils

import "os"

// LoadSetting looks up a value from an environment variable.
// If it is not set, the defaultValue is returned.
func LoadSetting(name string, defaultValue string) string {
	if v, ok := os.LookupEnv(name); ok {
		return v
	}
	return defaultValue
}
