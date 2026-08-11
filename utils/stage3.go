package utils

import (
	"mywebdav/types"
)

func DoubleSize(f *types.FileNode) {
	f.Size *= 2
}
