package service

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/rs/zerolog"
)

//go:generate mockery --name=IUploadService
type IUploadService interface {
	SaveLocal(path string, imageFile *multipart.FileHeader, unique string) (string, error)
	RemoveLocal(baseDir, imagePath string) error
}

type uploadService struct {
	log zerolog.Logger
}

func NewUploadService(log zerolog.Logger) IUploadService {
	return &uploadService{log}
}

func (s *uploadService) SaveLocal(path string, imageFile *multipart.FileHeader, unique string) (string, error) {
	timeUnix := time.Now().Unix()
	fileName := fmt.Sprintf("%s-%d-%s", unique, timeUnix, imageFile.Filename)
	filePath := filepath.Join(path, fileName)

	// Open the image file
	src, err := imageFile.Open()
	if err != nil {
		s.log.Error().
			Str("filename", imageFile.Filename).
			Err(err).
			Msg("open image file failed")
		return "", errors.New("unable to open the uploaded image file")
	}
	defer src.Close()

	// Create directory if it doesn't exist
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		s.log.Error().
			Str("path", path).
			Err(err).
			Msg("create directory failed")
		return "", errors.New("unable to create target directory")
	}

	// Create file at the target path
	dest, err := os.Create(filePath)
	if err != nil {
		s.log.Error().
			Str("path", filePath).
			Err(err).
			Msg("create file failed")
		return "", errors.New("unable to create file at specified path")
	}
	defer dest.Close()

	// Copy content from source to destination file
	if _, err := io.Copy(dest, src); err != nil {
		s.log.Error().
			Str("destination", filePath).
			Err(err).
			Msg("copy file failed")
		return "", errors.New("unable to copy file content")
	}

	return filePath, nil
}

func (s *uploadService) RemoveLocal(baseDir, imagePath string) error {
	fullPath := filepath.Join(baseDir, imagePath)

	// Remove the file
	if err := os.Remove(fullPath); err != nil {
		s.log.Error().
			Str("path", fullPath).
			Err(err).
			Msg("remove file failed")
		return errors.New("unable to remove file at specified path")
	}

	return nil
}
