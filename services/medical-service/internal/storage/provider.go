package storage

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

// Provider is a generic storage interface.
// Implementations can be local filesystem, S3, GCS, etc.
type Provider interface {
	Upload(ctx context.Context, key string, data []byte, contentType string) (string, error)
	Delete(ctx context.Context, key string) error
}

// --- Local filesystem provider ---

type localStorageProvider struct {
	basePath string
	logger   *logrus.Logger
}

// NewLocalStorageProvider creates a Provider that writes files to the local filesystem
// under the given basePath directory.
func NewLocalStorageProvider(basePath string, logger *logrus.Logger) Provider {
	return &localStorageProvider{
		basePath: basePath,
		logger:   logger,
	}
}

func (p *localStorageProvider) Upload(ctx context.Context, key string, data []byte, contentType string) (string, error) {
	fullPath := filepath.Join(p.basePath, key)

	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	if err := os.WriteFile(fullPath, data, 0644); err != nil {
		return "", fmt.Errorf("failed to write file %s: %w", fullPath, err)
	}

	p.logger.Infof("LocalStorage: saved %s (%d bytes, %s)", fullPath, len(data), contentType)
	return fullPath, nil
}

func (p *localStorageProvider) Delete(ctx context.Context, key string) error {
	fullPath := filepath.Join(p.basePath, key)

	if err := os.Remove(fullPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete file %s: %w", fullPath, err)
	}

	p.logger.Infof("LocalStorage: deleted %s", fullPath)
	return nil
}
