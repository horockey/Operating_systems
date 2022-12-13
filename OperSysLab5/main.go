package main

import (
	"log"
	"main/file_system"
)

func main() {
	fs := file_system.NewFileSystem[int](&file_system.FileSystemParams{
		StorageFileName: `D:\Repos\Operating_systems\OperSysLab5\assets\fs0`,
		BlockSize:       8,
		BlockCount:      16,
	})
	blks, err := fs.Allocate(3)
	fatalOnErr(err)
	err = fs.Write([]int{7, 4, 61}, blks)
	fatalOnErr(err)
	fs.Save()
}

func fatalOnErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
