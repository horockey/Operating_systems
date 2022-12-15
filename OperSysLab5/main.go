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
	fs.WriteToFile("/1.txt", []byte("Hello world"))
	fmt.Println()
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
