package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"example.com.backend/pkg/logger"
	"go.uber.org/zap"
)

type IUploadService interface {
	SaveLocal(ctx context.Context, path string, imageFile *multipart.FileHeader, unique string) (string, error)
	RemoveLocal(ctx context.Context, baseDir, imagePath string) error
}

type uploadService struct{}

func NewUploadService() IUploadService {
	return &uploadService{}
}

func (s *uploadService) SaveLocal(ctx context.Context, path string, imageFile *multipart.FileHeader, unique string) (string, error) {
	log := logger.GetLoggerFromContext(ctx)

	timeUnix := time.Now().Unix()
	fileName := fmt.Sprintf("%s-%d-%s", unique, timeUnix, imageFile.Filename)
	filePath := filepath.Join(path, fileName)

	// Open the uploaded image file
	src, err := imageFile.Open()
	if err != nil {
		log.Error("open image file failed",
			zap.String("upload_filename", imageFile.Filename),
			zap.Error(err),
		)
		return "", errors.New("unable to open the uploaded image file")
	}
	defer src.Close()

	// Ensure directory exists
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		log.Error("create directory failed",
			zap.String("upload_path", path),
			zap.Error(err),
		)
		return "", errors.New("unable to create target directory")
	}

	// Create destination file
	dest, err := os.Create(filePath)
	if err != nil {
		log.Error("create file failed",
			zap.String("path", filePath),
			zap.Error(err),
		)
		return "", errors.New("unable to create file at specified path")
	}
	defer dest.Close()

	// Copy content
	if _, err := io.Copy(dest, src); err != nil {
		log.Error("copy file failed",
			zap.String("upload_destination", filePath),
			zap.Error(err),
		)
		return "", errors.New("unable to copy file content")
	}

	log.Info("successfully saved local file",
		zap.String("upload_path", filePath),
		zap.String("upload_original_filename", imageFile.Filename),
	)
	return filePath, nil
}

func (s *uploadService) RemoveLocal(ctx context.Context, baseDir, imagePath string) error {
	log := logger.GetLoggerFromContext(ctx)

	fullPath := filepath.Join(baseDir, imagePath)

	// Remove the file
	if err := os.Remove(fullPath); err != nil {
		log.Error("remove file failed",
			zap.String("upload_path", fullPath),
			zap.Error(err),
		)
		return errors.New("unable to remove file at specified path")
	}

	log.Info("successfully removed local file", zap.String("upload_path", fullPath))
	return nil
}
