package models

type FileMeta struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Size int64  `json:"size"`
	Link string `json:"link"`
}

func NewFileMeta(name string, size int64, link string) *FileMeta {
	return &FileMeta{Name: name, Size: size, Link: link}
}
