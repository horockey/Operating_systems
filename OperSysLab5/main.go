package main

import (
	"fmt"
	"log"
	"main/file_system"
)

func main() {
	// lab5ShowOff()

	fs := file_system.NewFileSystem()
	err := fs.NewFile("/1.txt")
	fatalOnErr(err)
	fs.WriteFile("/1.txt", []byte("Hello world"))
	data, err := fs.ReadFile("/1.txt")
	fatalOnErr(err)
	fmt.Printf("%s\n", string(data))
	fs.NewDir("/my_folder")
	fatalOnErr(fs.NewFile("/my_folder/2.txt"))
	fatalOnErr(fs.WriteFile("/my_folder/2.txt", []byte("file 2 content")))
	data, err = fs.ReadFile("/my_folder/2.txt")
	fatalOnErr(err)
	fmt.Printf("%s\n", string(data))
	// fatalOnErr(fs.DelDir("/my_folder"))
	fs.Memory.Save()
	fs.GetDirInfo("/")
}

func lab5ShowOff() {
	fs := file_system.NewMemory[int](&file_system.MemoryParams{
		StorageFileName: `D:\Repos\Operating_systems\OperSysLab5\assets\fs0`,
		BlockSize:       8,
		BlockCount:      16,
	})
	blks, err := fs.Allocate(3)
	fatalOnErr(err)
	fatalOnErr(fs.Write([]int{7, 4, 61}, blks))
	fatalOnErr(fs.Save())
	fatalOnErr(fs.Free(blks))
	fatalOnErr(fs.Load())
	buf := make([]int, 3)
	fatalOnErr(fs.Read(blks, buf))
	for _, el := range buf {
		fmt.Printf("%d ", el)
	}
	fmt.Println()
}

func fatalOnErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
