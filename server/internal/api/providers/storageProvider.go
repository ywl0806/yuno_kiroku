package providers

import (
	"github.com/spf13/viper"
	"github.com/ywl0806/yuno_kiroku/pkg/storage"
)

type StorageProvider struct {
	thumbnailStorage storage.StorageService
	originalStorage  storage.StorageService
}

func NewStorageProvider() *StorageProvider {

	storageType := viper.GetString("STORAGE_TYPE")

	switch storageType {
	case "local":
		thumbnailRootDir := viper.GetString("THUMBNAIL_ROOT_DIR")
		originalRootDir := viper.GetString("ORIGINAL_ROOT_DIR")

		if thumbnailRootDir == "" || originalRootDir == "" {
			panic("thumbnail root dir or original root dir is not set")
		}

		return &StorageProvider{
			thumbnailStorage: storage.NewLocalStorageService(thumbnailRootDir),
			originalStorage:  storage.NewLocalStorageService(originalRootDir),
		}
	case "s3",
		"minio":
		thumbnailBucket := viper.GetString("THUMBNAIL_BUCKET")
		originalBucket := viper.GetString("ORIGINAL_BUCKET")

		if thumbnailBucket == "" || originalBucket == "" {
			panic("thumbnail bucket or original bucket is not set")
		}

		return &StorageProvider{
			thumbnailStorage: storage.NewS3StorageService(thumbnailBucket),
			originalStorage:  storage.NewS3StorageService(originalBucket),
		}
	default:
		panic("invalid storage type")
	}
}

func (s *StorageProvider) ThumbnailStorage() storage.StorageService {
	return s.thumbnailStorage
}

func (s *StorageProvider) OriginalStorage() storage.StorageService {
	return s.originalStorage
}
