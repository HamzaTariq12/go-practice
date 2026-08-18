package types

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

var ErrNotFound = errors.New("file not found")
var ErrNotADirectory = errors.New("not a directory")

// Note Receiver funcions should be in same struct package
type FileNode struct {
	Name     string
	IsDir    bool
	Size     int64
	Modified time.Time
	Data     []byte               // file contents (nil if it's a directory)
	Children map[string]*FileNode // only used if IsDir == true
}

func (f *FileNode) AddChild(child *FileNode) error {
	if !f.IsDir {
		return fmt.Errorf("%s is not a directory", f.Name)
	}
	if f.Children == nil {
		f.Children = make(map[string]*FileNode) // maps must be initialized before use!
	}
	f.Children[child.Name] = child
	child.Modified = time.Now()
	return nil
}

func (f *FileNode) RemoveChild(name string) bool {
	_, ok := f.Children[name]
	if ok {
		delete(f.Children, name)
		return true
	}
	return false
}

func (f *FileNode) Describe() string {
	if f.IsDir {
		return fmt.Sprintf("%s (directory, %d items)", f.Name, len(f.Children))
	}
	return fmt.Sprintf("%s (file, %d bytes)", f.Name, f.Size)
}

func (f *FileNode) TotalDescendants() int {
	if !f.IsDir {
		return 0
	}

	count := 0

	for _, child := range f.Children {
		count += 1
		count += child.TotalDescendants()
	}

	return count
}

func (f *FileNode) FindDeep(name string) (*FileNode, bool) {
	if !f.IsDir {
		return nil, false
	}

	node, ok := f.Children[name]

	if ok {
		return node, true
	}

	for _, child := range f.Children {
		node, status := child.FindDeep(name)
		if status {
			return node, status
		}
	}

	return nil, false
}

// Stage 4
func (f *FileNode) SizeInBytes() int64 {
	return f.Size
}

// Stat acts as an interface adapter/wrapper method that delegates search requests
// to internal lookup methods (like FindDeep).
//
// WHY THIS EXISTS (Interface Adapters & Multi-Backend Support):
//  1. Interface Compliance: Go interfaces require EXACT method names (e.g., Stat).
//     This wrapper exposes your internal FindDeep logic under the standard Stat name
//     so FileNode can be passed into external frameworks (like WebDAV servers).
//
// 2. Multi-Backend Polymorphism: Different storage backends implement Stat differently:
//   - MemoryFS (this struct) -> calls f.FindDeep(name) (RAM tree traversal)
//   - LocalDiskFS           -> calls os.Stat(path) (System OS call)
//   - S3CloudFS             -> calls s3Client.HeadObject(...) (Network API call)
//
// Because all backends implement Stat(name string) (*FileNode, bool), web servers
// can process files seamlessly across Memory, Disk, or S3 using a single shared
// interface type without caring where or how the data is stored.
func (f *FileNode) Stat(name string) (*FileNode, bool) {
	return f.FindDeep(name)
}

// Stage 5
func (f *FileNode) StatErr(name string) (*FileNode, error) {
	if !f.IsDir {
		err := &PathError{
			Op:   "when opening",
			Path: name,
			Err:  ErrNotADirectory,
		}
		return nil, err
	}
	val, ok := f.FindDeep(name)
	if !ok {
		err := &PathError{
			Op:   "when finding",
			Path: name,
			Err:  ErrNotFound,
		}
		return nil, err
	}
	return val, nil
}

func SafeStat(f *FileNode, name string) (node *FileNode, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("recovered from panic: %v", r) // turn panic into a normal error
		}
	}()

	if f == nil {
		panic("triggering a nil pointer dereference")
	}
	if !f.IsDir {
		return nil, ErrNotADirectory
	}
	val, ok := f.FindDeep(name)
	if !ok {
		return nil, ErrNotFound
	}
	return val, nil
}

// Stage 6
type FileNode6 struct {
	Name     string
	IsDir    bool
	Size     int64
	Modified time.Time
	Data     []byte                // file contents (nil if it's a directory)
	Children map[string]*FileNode6 // only used if IsDir == true

	mu sync.Mutex
}

func (f *FileNode6) AddChild6(child *FileNode6) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if !f.IsDir {
		return fmt.Errorf("%s is not a directory", f.Name)
	}
	if f.Children == nil {
		f.Children = make(map[string]*FileNode6) // maps must be initialized before use!
	}
	f.Children[child.Name] = child
	child.Modified = time.Now()
	return nil
}

func (f *FileNode6) RemoveChild6(name string) bool {
	f.mu.Lock()
	defer f.mu.Unlock()
	_, ok := f.Children[name]
	if ok {
		delete(f.Children, name)
		return true
	}
	return false
}
