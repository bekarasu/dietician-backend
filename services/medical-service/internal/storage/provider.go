package storage

import (
	"context"

	"github.com/sirupsen/logrus"
)

type Provider interface {
	Upload(ctx context.Context, key string, data []byte, contentType string) (string, error)
	Delete(ctx context.Context, key string) error
}

type noOpStorageProvider struct {
	logger *logrus.Logger
}

func NewNoOpStorageProvider(logger *logrus.Logger) Provider {
	return &noOpStorageProvider{
		logger: logger,
	}
}

func (p *noOpStorageProvider) Upload(ctx context.Context, key string, data []byte, contentType string) (string, error) {
	p.logger.Infof("NoOpStorageProvider: Uploading %s (size: %d bytes, type: %s)", key, len(data), contentType)
	// Return a dummy URL
	return "https://dummy-storage.local/" + key, nil
}

func (p *noOpStorageProvider) Delete(ctx context.Context, key string) error {
	p.logger.Infof("NoOpStorageProvider: Deleting %s", key)
	return nil
}
