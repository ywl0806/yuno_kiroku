package enums

type FailureReason string

const (
	FailureReasonUnsupportedFormat FailureReason = "unsupported_format"
	FailureReasonCorruptedFile     FailureReason = "corrupted_file"
	FailureReasonFileTooLarge      FailureReason = "file_too_large"
)
