package file_system

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Record struct {
	StartIndex int
	Length     int
}

type Memory[T any] struct {
	storageFileName string
	blockSize       int
	blockCount      int

	freeBlocks     []*Record
	occupiedBlocks []*Record

	memory []*MemoryBlock[T]
}

type MemoryParams struct {
	StorageFileName string
	BlockSize       int
	BlockCount      int
}

func NewMemory[T any](params *MemoryParams) *Memory[T] {
	memory := make([]*MemoryBlock[T], params.BlockCount)
	for i := params.BlockCount - 1; i >= 0; i-- {
		var next *MemoryBlock[T]
		if i < params.BlockCount-1 {
			next = memory[i+1]
		}
		memory[i] = &MemoryBlock[T]{IsFree: true, next: next}
	}
	return &Memory[T]{
		storageFileName: params.StorageFileName,
		blockSize:       params.BlockSize,
		blockCount:      params.BlockCount,
		freeBlocks:      []*Record{{StartIndex: 0, Length: params.BlockCount}},
		occupiedBlocks:  []*Record{},
		memory:          memory,
	}
}

func (fs *Memory[T]) Allocate(blocks int) ([]int, error) {
	blocksLeft := blocks
	var last *MemoryBlock[T]
	res := []int{}
	for idx := range fs.memory {
		if !fs.memory[idx].IsFree {
			continue
		}
		res = append(res, idx)
		blocksLeft--
		if last != nil {
			last.next = fs.memory[idx]
		}
		last = fs.memory[idx]

		if blocksLeft == 0 {
			break
		}
	}
	if blocksLeft > 0 {
		return nil, fmt.Errorf("not enought memory")
	}
	return res, nil
}

func (fs *Memory[T]) Free(blocks []int) error {
	for _, block := range blocks {
		if block < 0 || block > len(fs.memory) {
			return fmt.Errorf("invalid block number")
		}
		fs.memory[block].IsFree = true
	}
	return nil
}

func (fs *Memory[T]) Write(data []T, blockNumbers []int) error {
	if len(data) != len(blockNumbers) {
		return fmt.Errorf("wrong input")
	}
	for idx, block := range blockNumbers {
		if block < 0 || block > len(fs.memory) {
			return fmt.Errorf("invalid block number")
		}
		fs.memory[block].IsFree = false
		fs.memory[block].Data = data[idx]
	}
	return nil
}

func (fs *Memory[T]) Read(blockNumbers []int, buf []T) error {
	if len(blockNumbers) != len(buf) {
		return fmt.Errorf("wrong input")
	}
	for idx, block := range blockNumbers {
		if block < 0 || block > len(fs.memory) {
			return fmt.Errorf("invalid block number")
		}
		buf[idx] = fs.memory[block].Data
	}
	return nil
}

func (fs *Memory[T]) Save() error {
	obj := struct {
		Memory []*MemoryBlock[T] `json:"memory"`
	}{
		Memory: fs.memory,
	}
	f, err := os.Create(fs.storageFileName)
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(obj, "", "\t")
	if err != nil {
		return err
	}
	_, err = f.Write(data)
	return err
}

func (fs *Memory[T]) Load() error {
	f, err := os.Open(fs.storageFileName)
	if err != nil {
		return err
	}
	var data []byte
	data, err = ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	obj := struct {
		Memory []*MemoryBlock[T] `json:"memory"`
	}{}
	err = json.Unmarshal(data, &obj)
	fs.memory = obj.Memory
	return err
}
