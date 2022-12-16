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
			BlockCount:      64,
		}),
		RootDir: &Directory{
			Name:     "/",
			Content:  make(map[string]fsElementer),
			MetaInfo: make(map[string]*MetaInfo),
		},
	}
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
	curDir.MetaInfo[shortName] = &MetaInfo{
		Size:         0,
		LastModified: time.Now(),
		Type:         "File",
	}
	return nil
}
func (fs *FileSystem) WriteFile(fileName string, data []byte) error {
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
		curDir.MetaInfo[shortName] = &MetaInfo{
			Size:         len(blocks) * fs.Memory.blockSize,
			LastModified: time.Now(),
		}
	} else {
		return fmt.Errorf("wrong file name: %s", fileName)
	}
	return nil
}
func (fs *FileSystem) ReadFile(fileName string) ([]byte, error) {
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
			return nil, fmt.Errorf("wrong file name: %s", fileName)
		}
	}
	shortName := dirs[len(dirs)-1]
	if val, ok := curDir.Content[shortName]; !ok || !val.IsFile() {
		return nil, fmt.Errorf("wrong file name: %s", fileName)
	} else {
		res := []byte{}
		for _, bl := range curDir.Content[shortName].(*File).Blocks {
			res = append(res, bl.Data)
		}
		return res, nil
	}
}
func (fs *FileSystem) DelFile(fileName string) error {
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
			return fmt.Errorf("wrong file name: %s", fileName)
		}
	}
	shortName := dirs[len(dirs)-1]
	if val, ok := curDir.Content[shortName]; !ok || !val.IsFile() {
		return fmt.Errorf("wrong file name: %s", fileName)
	} else {
		for _, bl := range curDir.Content[shortName].(*File).Blocks {
			bl.Data = 0
			bl.IsFree = true
		}
		delete(curDir.Content, shortName)
		delete(curDir.MetaInfo, shortName)
		return nil
	}
}

func (fs *FileSystem) NewDir(name string) error {
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
			return fmt.Errorf("wrong dir name: %s", name)
		}
	}
	shortName := dirs[len(dirs)-1]
	curDir.Content[shortName] = &Directory{
		Name:     name,
		Content:  make(map[string]fsElementer),
		MetaInfo: make(map[string]*MetaInfo),
	}
	curDir.MetaInfo[shortName] = &MetaInfo{
		Size:         0,
		LastModified: time.Now(),
		Type:         "Folder",
	}
	return nil
}
func (fs *FileSystem) DelDir(name string) error {
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
			return fmt.Errorf("wrong dir name: %s", name)
		}
	}
	shortName := dirs[len(dirs)-1]
	if val, ok := curDir.Content[shortName]; !ok || !curDir.Content[shortName].IsDir() {
		return fmt.Errorf("wrong path: %s", name)
	} else {
		nestedContent := val.(*Directory).Content
		for key, elem := range nestedContent {
			if elem.IsFile() {
				fs.DelFile(fmt.Sprintf("%s/%s", name, key))
			} else {
				fs.DelDir(fmt.Sprintf("%s/%s", name, key))
			}
		}
		delete(curDir.Content, shortName)
		return nil
	}
}
func (fs *FileSystem) GetDirInfo(name string) error {
	dirs := strings.Split(name, "/")
	if dirs[len(dirs)-1] == "" {
		dirs = dirs[:len(dirs)-1]
	}
	curDir := fs.RootDir
	for _, dir := range dirs[1:] {
		haveElem := false
		for _, contentElem := range curDir.Content {
			if contentElem.IsDir() && contentElem.GetName() == dir {
				curDir = contentElem.(*Directory)
				haveElem = true
				break
			}
		}
		if !haveElem {
			return fmt.Errorf("wrong dir name: %s", name)
		}
	}
	fmt.Printf("===%s info===\n", name)
	for k, v := range curDir.MetaInfo {
		fmt.Printf("%s:\nType:%s\nSize: %db\nLast modified: %s\n\n",
			k, v.Type, v.Size, v.LastModified.Format("2006-01-02 15:04:05"))
	}
	return nil
}

type MetaInfo struct {
	Size         int
	LastModified time.Time
	Type         string
}
