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

	if storageType == "" {
		storageType = "s3"
	}

	switch storageType {
	case "local":
		return &StorageProvider{
			storageService: storage.NewLocalStorageService("uploads"),
		}
	case "s3",
		"minio":
		bucket := viper.GetString("MEDIA_BUCKET_NAME")

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
