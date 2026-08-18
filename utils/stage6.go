package utils

import (
	"fmt"
	"mywebdav/types"
	"sync"
)

func ProcessFiles(names []string) {
	var wg sync.WaitGroup
	for _, name := range names {
		wg.Add(1)
		go func(fileName string) {
			defer wg.Done()
			fmt.Println("processing " + fileName)
		}(name)
	}
	wg.Wait()
	fmt.Println("All files processed!")
}

func CollectSizes(files []*types.FileNode6) int64 {
	var wg sync.WaitGroup
	var sizeColl int64
	chan1 := make(chan int64)
	for _, file := range files {
		wg.Add(1)
		go func(file *types.FileNode6) {
			defer wg.Done()
			chan1 <- file.Size
		}(file)
	}

	// this is required otherwise the loop below that is reading from channel will give deadlock as
	// it will continue to read until channel is closed
	// Note: We used go routine as wg.Wait will block the thread and we wont be able to use for loop to read from channels
	go func() {
		wg.Wait()    // counts all goroutines
		close(chan1) // below wg.Wait() so we can close channel and for loop reading once go routines finishes
	}()

	// for loop here works like sizeColl += <- chan1 and its continuously reading from channel until closed
	for value := range chan1 {
		sizeColl += value
	}
	return sizeColl
}
