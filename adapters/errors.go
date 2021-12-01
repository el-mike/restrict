package adapters

import "fmt"

// FileTypeNotSupportedError - thrown when FileAdapter is used with inappropriate
// file type.
type FileTypeNotSupportedError struct {
	fileType string
}

// NewFileTypeNotSupportedError - returns new FileTypeNotSupportedError instance.
func NewFileTypeNotSupportedError(fileType string) *FileTypeNotSupportedError {
	return &FileTypeNotSupportedError{
		fileType: fileType,
	}
}

// Error - error interface implementation.
func (e *FileTypeNotSupportedError) Error() string {
	return fmt.Sprintf("File type: \"%s\" is not supported", e.fileType)
}
