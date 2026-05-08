package enums

type UploadStatus string

const (
	// 01: pending 업로드 대기
	UploadStatusPending UploadStatus = "01"
	// 02: processing 업로드 처리중
	UploadStatusProcessing UploadStatus = "02"
	// 03: completed 업로드 완료
	UploadStatusCompleted UploadStatus = "03"
	// 04: failed 업로드 실패
	UploadStatusFailed UploadStatus = "04"
	// 05: duplicate 업로드 중복
	UploadStatusDuplicate UploadStatus = "05"
)
