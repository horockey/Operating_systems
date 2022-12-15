package file_system

import (
	"fmt"
	"math"
	"strings"
	"time"
)

type FileSystem struct {
	Memory  *Memory[byte]
	RootDir *Directory
}

func NewFileSystem() *FileSystem {
	return &FileSystem{
		Memory: NewMemory[byte](&MemoryParams{
			StorageFileName: `D:\Repos\Operating_systems\OperSysLab5\assets\fs6`,
			BlockSize:       1,
			BlockCount:      16,
		}),
		RootDir: NewDir("/"),
	}
}

type File struct {
	Name   string
	Blocks []*MemoryBlock[byte]
}

func (fs *FileSystem) NewFile(name string) error {
	dirs := strings.Split(name, "/")
	curDir := fs.RootDir
	for _, dir := range dirs[1 : len(dirs)-1] {
		haveElem := false
		for _, contentElem := range curDir.Content {
			if contentElem.IsDir() && contentElem.GetName() == dir {
				curDir = contentElem.(*Directory)
				haveElem = true
				break
			}
		}
		if !haveElem {
			return fmt.Errorf("wrong file name: %s", name)
		}
	}
	shortName := dirs[len(dirs)-1]
	curDir.Content[shortName] = &File{
		Name:   name,
		Blocks: []*MemoryBlock[byte]{},
	}
	return nil
}
func (fs *FileSystem) WriteToFile(fileName string, data []byte) error {
	blocksCount := math.Ceil(float64(len(data)) / float64(fs.Memory.blockSize))
	blocks, err := fs.Memory.Allocate(int(blocksCount))
	if err != nil {
		return err
	}
	fs.Memory.Write(data, blocks)
	dirs := strings.Split(fileName, "/")
	curDir := fs.RootDir
	for _, dir := range dirs[1 : len(dirs)-1] {
		haveElem := false
		for _, contentElem := range curDir.Content {
			if contentElem.IsDir() && contentElem.GetName() == dir {
				curDir = contentElem.(*Directory)
				haveElem = true
				break
			}
		}
		if !haveElem {
			return fmt.Errorf("unknown file name: %s", fileName)
		}
	}
	shortName := dirs[len(dirs)-1]
	if val, ok := curDir.Content[shortName]; ok {
		for _, blockNumber := range blocks {
			val.(*File).Blocks = append(val.(*File).Blocks, fs.Memory.memory[blockNumber])
		}
	} else {
		return fmt.Errorf("wrong file name: %s", fileName)
	}
	return nil
}

type Directory struct {
	Name     string
	Content  map[string]fsElementer
	MetaInfo map[string]*MetaInfo
}

func NewDir(name string) *Directory {
	return &Directory{
		Name:     name,
		Content:  make(map[string]fsElementer),
		MetaInfo: make(map[string]*MetaInfo),
	}
}

type fsElementer interface {
	IsDir() bool
	IsFile() bool
	GetName() string
}

func (f *File) IsDir() bool          { return false }
func (f *File) IsFile() bool         { return true }
func (f *File) GetName() string      { return f.Name }
func (d *Directory) IsDir() bool     { return true }
func (d *Directory) IsFile() bool    { return false }
func (d *Directory) GetName() string { return d.Name }

type MetaInfo struct {
	Size         int
	LastModified time.Time
	Type         string
}
