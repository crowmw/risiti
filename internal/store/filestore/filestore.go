package filestore

import (
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"os"
)

type IFileStore interface {
	SaveFile(file multipart.File, filename string) error
}

type FileStore struct{}

func NewFileStore() *FileStore {
	return &FileStore{}
}

func (fs *FileStore) SaveFile(file multipart.File, filename string) error {
	defer file.Close()
	dst, err := os.Create(fmt.Sprintf("./bin/recipes/%s", filename))
	if err != nil {
		slog.Error("Creating file!", err)
		return err
	}

	defer dst.Close()
	if _, err := io.Copy(dst, file); err != nil {
		return err
	}
	return nil
}
