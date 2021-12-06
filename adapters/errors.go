package adapters

import "fmt"

// FileTypeNotSupportedError - thrown when FileAdapter is used with inappropriate
// file type.
type FileTypeNotSupportedError struct {
	fileType string
}

// newFileTypeNotSupportedError - returns new FileTypeNotSupportedError instance.
func newFileTypeNotSupportedError(fileType string) *FileTypeNotSupportedError {
	return &FileTypeNotSupportedError{
		fileType: fileType,
	}
}

// Error - error interface implementation.
func (e *FileTypeNotSupportedError) Error() string {
	return fmt.Sprintf("File type: \"%s\" is not supported", e.fileType)
}
