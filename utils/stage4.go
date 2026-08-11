// 1. Fundamentals
package utils

import (
	"fmt"
	"mywebdav/types"
)

func DescribeIfPossible(x any) {
	// nodePtr, ok := x.(*types.FileNode) (words same)
	switch v := x.(type) {
	case int:
		fmt.Println("This is an integer:", v)
	case string:
		fmt.Println("This is a string:", v)
	case *types.FileNode:
		fmt.Println("This is a filenode:", v.Describe())
	case types.LogEntry:
		fmt.Println("This is a logentry:", v)
	default:
		fmt.Println("Unknown type")
	}
}
