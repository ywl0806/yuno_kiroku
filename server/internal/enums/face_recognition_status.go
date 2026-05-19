package enums

type FaceRecognitionStatus string

const (
	// pending SQS 발행 전 (기본값)
	FaceRecognitionStatusPending FaceRecognitionStatus = "pending"
	// dispatched SQS 발행 완료
	FaceRecognitionStatusDispatched FaceRecognitionStatus = "dispatched"
	// completed 얼굴인식 처리 완료
	FaceRecognitionStatusCompleted FaceRecognitionStatus = "completed"
)
