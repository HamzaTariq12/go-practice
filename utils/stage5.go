package utils

import (
	"fmt"
	"mywebdav/types"
)

func closeFile(f *types.FileNode) {
	fmt.Println("closing " + f.Name)
}

func OpenAndCount(f *types.FileNode) int {
	fmt.Println("opening " + f.Name)
	defer closeFile(f)
	return int(f.Size)
}
