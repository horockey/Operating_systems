package file_system

import "strings"

type Directory struct {
	Name     string
	Content  map[string]fsElementer
	MetaInfo map[string]*MetaInfo
}

func (d *Directory) IsDir() bool  { return true }
func (d *Directory) IsFile() bool { return false }
func (d *Directory) GetName() string {
	splittedName := strings.Split(d.Name, "/")
	return splittedName[len(splittedName)-1]
}
