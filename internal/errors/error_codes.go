package errors

type ErrorCode struct{}

const (
	CodePdfGenerationErr     = "PDF_GENERATION_ERROR"
	CodePdfMergeErr          = "PDF_MERGE_ERROR"
	CodePdfEncryptionErr     = "PDF_ENCRYPTION_ERROR"
	CodeUnreadlableFileErr   = "UNREADABLE_FILE"
	CodeValidationErr        = "VALIDATION_ERROR"
	CodeUnknownErr           = "UNKNOWN_ERROR"
	CodeUndecodableBase64Err = "UNDECODABLE_BASE64"
)
