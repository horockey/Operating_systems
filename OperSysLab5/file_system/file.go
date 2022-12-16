package file_system

import "strings"

type File struct {
	Name   string
	Blocks []*MemoryBlock[byte]
}

func (f *File) IsDir() bool  { return false }
func (f *File) IsFile() bool { return true }
func (f *File) GetName() string {
	splittedName := strings.Split(f.Name, "/")
	return splittedName[len(splittedName)-1]
}
