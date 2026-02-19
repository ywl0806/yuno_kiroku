package storage

// Storage서비스의 인터페이스
type StorageService interface {
	SaveFile(file []byte, filePath string, fileName string) (string, error)
}
