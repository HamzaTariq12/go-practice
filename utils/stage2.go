// 1. Control Flow & Data Structures
package utils

import (
	"fmt"
	"mywebdav/types"
)

const DirectorySize = 2

func AddFileToDirectory(directory *[]string, fileName string) (string, error) {
	if len(*directory) >= DirectorySize {
		return "", fmt.Errorf("Directory is full.")
	}
	*directory = append(*directory, fileName)
	return "File added to directory successfully.", nil
}

func TotalSize(files []types.FileNode) int64 {
	var sum int64
	for _, file := range files {
		sum += file.Size
	}
	return sum
}

func FindFile(root types.FileNode, name string) (*types.FileNode, bool) {
	root_files := root.Children
	val, ok := root_files[name]
	if ok {
		return val, true
	}
	return nil, false
}

func ListDirectory(dir types.FileNode) []string {
	dir_files := dir.Children
	var files []string
	for k := range dir_files {
		if k[0] == '.' {
			continue
		}
		files = append(files, k)
	}
	return files
}

func CountFilesByType(dir types.FileNode) map[string]int {
	dir_files := dir.Children
	var fileTypeCounts = map[string]int{
		"image": 0,
		"text":  0,
	}
	for k := range dir_files {
		for i := len(k) - 1; i >= 0; i-- {
			if k[i] == '.' {
				mimeType := k[i+1:]
				switch mimeType {
				case "txt":
					fileTypeCounts["text"] += 1
				case "png":
					fileTypeCounts["image"] += 1
				}
			}
		}
	}
	return fileTypeCounts
}
