package infrastructure

import (
	"context"
	"os"
)

type FileRepositoryImpl struct{}

func NewFileRepository() *FileRepositoryImpl {
	return &FileRepositoryImpl{}
}

func (r *FileRepositoryImpl) Exists(ctx context.Context, filePath string) (bool, error) {
	_, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (r *FileRepositoryImpl) Write(ctx context.Context, filePath string, content string) error {
	return os.WriteFile(filePath, []byte(content), 0644)
}