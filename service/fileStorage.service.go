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

type FileStorage interface {
	SaveFile(file multipart.File, filename string) error
}

type fileStorage struct{}

func DefaultFileStorage() FileStorage {
	return &fileStorage{}
}

func (fs *fileStorage) SaveFile(file multipart.File, filename string) error {
	defer file.Close()
	dst, err := os.Create(fmt.Sprintf("%s%s", RECEIPTS_PATH, filename))
	if err != nil {
		slog.Error("While creating file", err)
		return err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return err
	}
	return nil
}
