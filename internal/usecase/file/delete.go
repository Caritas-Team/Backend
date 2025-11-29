package file

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Caritas-Team/reviewer/internal/logger"
	"github.com/Caritas-Team/reviewer/internal/memcached"
)

type Cleaner struct {
	cache    *memcached.Cache
	filesDir string
	log      *logger.Logger
}

type fileMetadata struct {
	UUID     string `json:"uuid"`
	Status   string `json:"status"`
	Filename string `json:"filename"`
}

func NewFileCleaner(log *logger.Logger, cache *memcached.Cache) *Cleaner {
	return &Cleaner{
		cache:    cache,
		filesDir: "./files",
		log:      log,
	}
}

func (fc *Cleaner) DeleteDownloadedFiles(ctx context.Context) error {
	files, err := os.ReadDir(fc.filesDir)
	if err != nil {
		return fmt.Errorf("error reading directory /files: %w", err)
	}

	fc.log.Info("Found files in directory", "count", len(files), "files", files)

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		filename := file.Name()
		uuid := filename[:len(filename)-len(filepath.Ext(filename))]

		status, err := fc.getFileStatus(ctx, uuid)
		if err != nil {
			fc.log.Warn("Cannot get file status", "uuid", uuid, "err", err)
			continue
		}

		if status == "DOWNLOADED" {
			filePath := filepath.Join(fc.filesDir, filename)
			if err = os.Remove(filePath); err != nil {
				fc.log.Error("Cannot remove downloaded file", "uuid", uuid, "path", filePath, "err", err)
				continue
			}

			if err = fc.cache.Delete(ctx, uuid); err != nil {
				fc.log.Error("Error removing data", "uuid", uuid, "error", err)
				continue
			}

			fc.log.Info("Removed file", "uuid", uuid)
		}
	}

	return nil
}

func (fc *Cleaner) getFileStatus(ctx context.Context, uuid string) (string, error) {
	data, err := fc.cache.Get(ctx, uuid)
	if err != nil {
		return "", fmt.Errorf("error getting file from cache: %w", err)
	}

	var metadata fileMetadata
	if err = json.Unmarshal(data, &metadata); err != nil {
		return "", fmt.Errorf("error unmarshalling file metadata: %w", err)
	}

	return metadata.Status, nil
}
