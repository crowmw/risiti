package service

import (
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"os"
)

const (
	RECEIPTS_PATH = "./data/"
)

type IFileStorage interface {
	SaveFile(file multipart.File, filename string) error
}

type FileStorage struct{}

func NewFileStorage() *FileStorage {
	return &FileStorage{}
}

func (fs *FileStorage) SaveFile(file multipart.File, filename string) error {
	defer file.Close()
	dst, err := os.Create(fmt.Sprintf("%s%s", RECEIPTS_PATH, filename))
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
