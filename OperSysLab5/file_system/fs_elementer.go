package file_system

type fsElementer interface {
	IsDir() bool
	IsFile() bool
	GetName() string
}
