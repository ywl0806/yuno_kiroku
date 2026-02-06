package providers

import (
	"github.com/spf13/viper"
	"github.com/ywl0806/yuno_kiroku/pkg/storage"
)

type StorageProvider struct {
	storageService storage.StorageService
}

func NewStorageProvider() *StorageProvider {

	storageType := viper.GetString("STORAGE_TYPE")

	switch storageType {
	case "local":
		rootDir := viper.GetString("STORAGE_ROOT_DIR")

		if rootDir == "" {
			panic("storage root dir is not set")
		}

		return &StorageProvider{
			storageService: storage.NewLocalStorageService(rootDir),
		}
	case "s3",
		"minio":
		bucket := viper.GetString("STORAGE_BUCKET")

		if bucket == "" {
			panic("storage bucket is not set")
		}

		return &StorageProvider{
			storageService: storage.NewS3StorageService(bucket),
		}
	default:
		panic("invalid storage type")
	}
}

func (s *StorageProvider) StorageService() storage.StorageService {
	return s.storageService
}
