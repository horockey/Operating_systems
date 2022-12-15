package file_system

type MemoryBlock[T any] struct {
	Data   T
	IsFree bool

	next *MemoryBlock[T]
}
