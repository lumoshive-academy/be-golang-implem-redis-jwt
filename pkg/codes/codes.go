package codes

import (
	"go-42/pkg/utils"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

func UploadFile(file multipart.File, fileName string, logger *zap.Logger, config utils.Configuration) error {
	defer file.Close()

	// Ensure the destination folder exists
	destDir := config.PathUpload
	err := os.MkdirAll(destDir, os.ModePerm)
	if err != nil {
		logger.Error("failed to create upload directory", zap.Error(err))
		return err
	}

	// Full path for saving the uploaded file
	dstPath := filepath.Join(destDir, fileName)

	// Create the destination file
	dstFile, err := os.Create(dstPath)
	if err != nil {
		logger.Error("failed to create destination file", zap.Error(err))
		return err
	}
	defer dstFile.Close()

	// Copy the contents of the uploaded file to the destination file
	_, err = io.Copy(dstFile, file)
	if err != nil {
		logger.Error("failed to copy file content", zap.Error(err))
		return err
	}

	logger.Info("file uploaded successfully", zap.String("filePath", dstPath))
	return nil
}

func GeneratePassword(password string) (*string, error) {
	// hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// convert string to *string
	hashedStr := string(hashedPassword)
	return &hashedStr, nil
}
