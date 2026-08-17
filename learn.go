package main

import (
	"errors"
	"fmt"
	interfaces "mywebdav/interface"
	"mywebdav/types"
	"mywebdav/utils"
	"time"
)

func main() {
	fmt.Println("Start Main!")

	// 1: Fundamentals
	//////////// START ////////////
	println("STAGE-1")
	value, err1 := utils.GetFileMimeType("a.txt")
	if err1 != nil {
		fmt.Printf("Error: %v", err1)
		return
	}
	fmt.Println(value)
	//////////// END ////////////

	// 2: Control Flow & Data Structures
	//////////// START ////////////
	println("STAGE-2")
	// Slices
	var directory []string
	var result string
	var err2 error
	result, err2 = utils.AddFileToDirectory(&directory, "file.txt")
	if err2 != nil {
		fmt.Printf("Error: %v", err2)
		return
	}
	result, err2 = utils.AddFileToDirectory(&directory, "drangon.png")
	if err2 != nil {
		fmt.Printf("Error: %v", err2)
		return
	}
	fmt.Printf("%v Directory: %v\n", result, directory)
	// Exercise
	// Q1
	files := []types.FileNode{
		{
			Name:     "document.pdf",
			Size:     2048,
			IsDir:    false,
			Modified: time.Now(),
		},
		{
			Name:     "images",
			Size:     4096,
			IsDir:    true,
			Modified: time.Now().Add(-24 * time.Hour),
		},
		{
			Name:     "dragon.png",
			Size:     1024,
			IsDir:    false,
			Modified: time.Now(),
		},
	}
	filesSum := utils.TotalSize(files)
	fmt.Printf("Total Files Size: %d\n", filesSum)
	// Q2
	dragonPicture := types.FileNode{
		Name:     "dragon.png",
		Size:     2048,
		IsDir:    false,
		Modified: time.Now(),
	}
	catPicture := types.FileNode{
		Name:     "cat.png",
		Size:     2048,
		IsDir:    false,
		Modified: time.Now(),
	}
	rootFolder := types.FileNode{
		Name:     "root",
		Size:     4096,
		IsDir:    true,
		Modified: time.Now(),
		Children: map[string]*types.FileNode{
			"dragon.png": &dragonPicture,
			"cat.png":    &catPicture,
		},
	}
	fileNode, status := utils.FindFile(rootFolder, "cat.png")
	if status {
		fmt.Println(fileNode)
	} else {
		fmt.Println("File not found.")
	}
	// Q3
	files2 := utils.ListDirectory(rootFolder)
	fmt.Println(files2)
	// Q4
	dirFileTypes := utils.CountFilesByType(rootFolder)
	fmt.Println(dirFileTypes)
	//////////// END ////////////

	// 3: Pointers & Methods
	//////////// START ////////////
	println("STAGE-3")
	// Q1
	utils.DoubleSize(&dragonPicture)
	fmt.Println(dragonPicture.Size)
	// Q2
	fmt.Println(rootFolder.RemoveChild("dragoncat.png"))
	// Q3
	fmt.Println(rootFolder.Describe())
	fmt.Println(dragonPicture.Describe())
	// Q4
	videoFile := types.FileNode{
		Name:     "video.mp4",
		Size:     6096,
		IsDir:    false,
		Modified: time.Now(),
	}
	textFile := types.FileNode{
		Name:     "file.txt",
		Size:     4096,
		IsDir:    false,
		Modified: time.Now(),
	}
	rootSubFolder1 := types.FileNode{
		Name:     "pictures",
		Size:     4096,
		IsDir:    true,
		Modified: time.Now(),
		Children: map[string]*types.FileNode{
			"dragon.png": &dragonPicture,
			"cat.png":    &catPicture,
		},
	}
	rootSubFolder2 := types.FileNode{
		Name:     "files",
		Size:     4096,
		IsDir:    true,
		Modified: time.Now(),
		Children: map[string]*types.FileNode{
			"video.mp4": &videoFile,
			"file.txt":  &textFile,
		},
	}
	rootFolderBonus := types.FileNode{
		Name:     "root",
		Size:     4096,
		IsDir:    true,
		Modified: time.Now(),
		Children: map[string]*types.FileNode{
			"pictures": &rootSubFolder1,
			"files":    &rootSubFolder2,
		},
	}
	fmt.Println(rootFolderBonus.TotalDescendants())
	fmt.Println(rootFolderBonus.FindDeep("video.mp4"))
	//////////// END ////////////

	// 4: Interfaces
	//////////// START ////////////
	println("STAGE-4")
	// Q1
	interfaces.PrintSize(&rootFolderBonus)
	// Q2
	logs5 := types.LogEntry{
		Message: "Interface log",
		Bytes:   1024,
	}
	interfaces.PrintSize(logs5)
	// Q3
	fmt.Println(rootFolderBonus.Stat("video.mp4"))
	// Q4
	utils.DescribeIfPossible(&rootFolderBonus)
	utils.DescribeIfPossible(logs5)
	utils.DescribeIfPossible("testing")
	//////////// END ////////////

	// 5: Error Handling
	//////////// START ////////////
	println("STAGE-5")
	// Q1
	utils.OpenAndCount(&dragonPicture)
	// Q2+Q3
	val, err5 := rootFolderBonus.StatErr("3video.mp4")
	if errors.Is(err5, types.ErrNotADirectory) {
		fmt.Println("This is not a directory.")
	} else if errors.Is(err5, types.ErrNotFound) {
		fmt.Println("File not found.")
	} else if err5 == nil {
		fmt.Println("File found", val)
	}
	// Q4
	var emptyFile *types.FileNode
	val5_4, err5_4 := types.SafeStat(emptyFile, "3video.mp4")
	if err5_4 != nil {
		fmt.Println("failed:", err5_4) // will print: "failed: recovered from panic: triggering a nil pointer dereference"
	} else {
		println("found safely", val5_4)
	}
	//////////// END ////////////

	// 6: Concurrency
	//////////// START ////////////
	println("STAGE-6")
	// Example
	ch1 := make(chan int)
	var msg int
	go func() {
		time.Sleep(1000 * time.Millisecond)
		ch1 <- 5
		fmt.Println("check-4")
	}()
	fmt.Println("check-1")
	time.Sleep(time.Second)
	msg = <-ch1
	fmt.Println("check-2", msg)
	fmt.Println("check-3")
	ch2 := make(chan int)
	select {
	case msg = <-ch2:
		fmt.Println("from 2nd channel", msg)
	case <-time.After(1000 * time.Millisecond): //check-4 gets printed because of time here (go routine gets extra time to execute)
		fmt.Println("waiting timeout")
	}
	// 2 Important concepts for concurrency (DEADLOCKs / RACE CONDITION)
	// there are 2 types of channels (buffered/unbuffered) if we define the channel size when
	// making them then they are buffered channels other wise unbuffered channel.
	// Note: Unbuffered channel has no storage so they require both sender/receiver at same time
	// if you define one of either of them then you will get DEADLOCK error
	// for buffered channels if the storage limit hits then you must define receiver before
	// adding more to it else you will get DEADLOCK
	// RACE CONDITION occurs when you try to change the same variable at same time. Like multiple
	// go routines updating the same variable at same time. To solve this problem use
	// sync.Mutex (Lock and Unlock) and if we want a go routine to wait for sometime before moving
	// to new go routine we can use sync.WaitGroup (wg.Add/wg.Done/wg.Wait)
	// More detail written in go_concurrency_notes file.

	// Q1
	//////////// END ////////////

	fmt.Println("End Main!")
}
