// 1. Fundamentals
package utils

import "fmt"

func GetFileMimeType(name string) (string, error) {
	if name == "" {
		return "", fmt.Errorf("FileName is empty")
	}

	for i := len(name) - 1; i >= 0; i-- {
		if name[i] == '.' {
			return name[i:], nil
		}
	}

	return "", fmt.Errorf("No extension found in file name.")
}
