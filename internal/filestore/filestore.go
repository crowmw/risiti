package filestore

import (
	"mime/multipart"
)

type FileStore interface {
	SaveFile(file multipart.File, filename string) error
}
