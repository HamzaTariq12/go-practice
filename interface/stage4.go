package interfaces

import (
	"fmt"
	"mywebdav/types"
)

type Sizeable interface {
	SizeInBytes() int64
}

func PrintSize(s Sizeable) {
	fmt.Println("Size in bytes:", s.SizeInBytes())
}

type Stater interface {
	Stat(name string) (*types.FileNode, bool)
}
