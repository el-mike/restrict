package adapters

import (
	"encoding/json"
	"os"

	"gopkg.in/yaml.v3"
)

type FilePerm os.FileMode

// FileReadWriter - facade interface for os read/write file functions.
type FileReadWriter interface {
	ReadFile(name string) ([]byte, error)
	WriteFile(name string, data []byte, perm FilePerm) error
}

// defaultFileHandler - fileReadWriter implementation.
type defaultFileHandler struct{}

func newDefaultFileHandler() *defaultFileHandler {
	return &defaultFileHandler{}
}

func (dh *defaultFileHandler) ReadFile(name string) ([]byte, error) {
	return os.ReadFile(name)
}

func (dh *defaultFileHandler) WriteFile(name string, data []byte, perm FilePerm) error {
	return os.WriteFile(name, data, os.FileMode(perm))
}

// JSONMarshalUnmarshaler - facade interface for json operations.
type JSONMarshalUnmarshaler interface {
	Unmarshal(data []byte, v interface{}) error
	MarshalIndent(v interface{}, prefix, indent string) ([]byte, error)
}

type defaultJSONHandler struct{}

func newDefaultJSONHandler() *defaultJSONHandler {
	return &defaultJSONHandler{}
}

func (dh *defaultJSONHandler) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func (dh *defaultJSONHandler) MarshalIndent(v interface{}, prefix, indent string) ([]byte, error) {
	return json.MarshalIndent(v, prefix, indent)
}

// YAMLMarshalUnmarshaler - facade interface for yaml operations
type YAMLMarshalUnmarshaler interface {
	Unmarshal(in []byte, out interface{}) (err error)
	Marshal(in interface{}) (out []byte, err error)
}

type defaultYAMLHandler struct{}

func newDefaultYAMLHandler() *defaultYAMLHandler {
	return &defaultYAMLHandler{}
}

func (dh *defaultYAMLHandler) Unmarshal(in []byte, out interface{}) error {
	return yaml.Unmarshal(in, out)
}

func (dh *defaultYAMLHandler) Marshal(in interface{}) ([]byte, error) {
	return yaml.Marshal(in)
}
