package util

import (
	"bytes"
	"encoding/gob"
	"os"
	"path/filepath"
)

type FileDB struct {
	filename string
}

func NewFileDB(filename string) *FileDB {
	return &FileDB{filename: filename}
}

func (f *FileDB) Store(v any) {
	buf := new(bytes.Buffer)
	gob.NewEncoder(buf).Encode(v)
	os.MkdirAll(filepath.Dir(f.filename), os.ModePerm)
	err := os.WriteFile(f.filename, buf.Bytes(), os.ModePerm)
	if err != nil {
		panic(err)
	}
}

func (f *FileDB) Load(v any) {
	content, err := os.ReadFile(f.filename)
	if err != nil {
		panic(err)
	}
	buf := new(bytes.Buffer)
	buf.Write(content)
	gob.NewDecoder(buf).Decode(v)
}

func (f *FileDB) Remove() {
	err := os.Remove(f.filename)
	if err != nil {
		panic(err)
	}
}
