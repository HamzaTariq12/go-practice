package utils

import (
	"errors"
	"io"
	"mywebdav/types"
	"os"
	"path/filepath"
	"strings"
)

func SafeJoin(root, userPath string) (string, error) {
	cleanRoot := filepath.Clean(root)
	fullPath := filepath.Clean(filepath.Join(cleanRoot, userPath))

	if fullPath != cleanRoot && !strings.HasPrefix(fullPath, cleanRoot+string(filepath.Separator)) {
		return "", errors.New("path traversal attempt detected")
	}
	return fullPath, nil
}

// func CheckAndCreateUserPath(root string, userPath string) (bool, error) {
// 	cleanRoot := filepath.Clean(root)
// 	fullPath := filepath.Clean(filepath.Join(cleanRoot, userPath))

// 	// Security check: prevent path traversal
// 	if !strings.HasPrefix(fullPath, cleanRoot) {
// 		return false, errors.New("path traversal attempt detected")
// 	}

// 	// Check if directory exists first
// 	info, err := os.Stat(fullPath)
// 	fmt.Println(info)
// 	if err == nil {
// 		// Path exists, check if it's a directory
// 		if info.IsDir() {
// 			return true, nil // Already exists, nothing to do
// 		}
// 		return false, errors.New("path exists but is not a directory")
// 	}

// 	// If error is not "not exists", return the error
// 	if !os.IsNotExist(err) {
// 		return false, err
// 	}

// 	// Directory doesn't exist, create it
// 	err = os.Mkdir(fullPath, 0755)
// 	if err != nil {
// 		return false, err
// 	}

// 	return true, nil
// }

func SaveNode(n *types.FileNode6, dir string) error {
	// Your specific path
	filePath := filepath.Join(dir, n.Name)

	// Ensure the directory exists
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}

	// Write "Hello World" to the file
	err = os.WriteFile(filePath, []byte("Hello World"), 0644)
	if err != nil {
		return err
	}

	return nil
}

func LoadNode(path string) (*types.FileNode6, error) {
	// Step 1: Open the file
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close() // Important: Always close the file

	// Step 2: Get file info using Stat
	info, err := file.Stat()
	if err != nil {
		return nil, err
	}

	// Step 3: Read all content using io.ReadAll
	content, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	// Step 4: Create and return FileNode
	node := &types.FileNode6{
		Name:     info.Name(),
		Size:     info.Size(),
		Modified: info.ModTime(),
		Data:     content,
	}

	return node, nil
}

func UserListDirectory(root string) ([]*types.FileNode6, error) {
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}

	var files []*types.FileNode6
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			return nil, err
		}
		files = append(files, &types.FileNode6{
			Name:     info.Name(),
			Size:     info.Size(),
			IsDir:    info.IsDir(), // now correctly flagged
			Modified: info.ModTime(),
			// Data intentionally omitted — a listing shouldn't load file contents
		})
	}
	return files, nil
}
