package adapters

import (
	"github.com/el-mike/restrict"
)

// AllowedFileType - alias type for describing allowed file types.
type AllowedFileType string

const (
	// JSONFile - JSON file type token.
	JSONFile AllowedFileType = "JSONFile"
	// YAMLFile - YAML file type token.
	YAMLFile AllowedFileType = "YAMLFile"
)

// defaultJSONIndent - default JSON file indentation.
const defaultJSONIndent = "\t"

// defaultFilePerm - default file's perm.
const defaultFilePerm FilePerm = 0644

// FileAdapter - StorageAdapter implementation, providing file-based persistence.
// It can be configured to use JSON or YAML format.
type FileAdapter struct {
	fileHandler FileReadWriter
	jsonHandler JSONMarshalUnmarshaler
	yamlHandler YAMLMarshalUnmarshaler

	fileName   string
	fileType   AllowedFileType
	filePerm   FilePerm
	jsonIndent string
}

// NewFileAdapter - returns new FileAdapter instance.
func NewFileAdapter(fileName string, fileType AllowedFileType) *FileAdapter {
	return &FileAdapter{
		fileHandler: newDefaultFileHandler(),
		jsonHandler: newDefaultJSONHandler(),
		yamlHandler: newDefaultYAMLHandler(),

		fileName:   fileName,
		fileType:   fileType,
		filePerm:   defaultFilePerm,
		jsonIndent: defaultJSONIndent,
	}
}

// SetJsonIndent - allows to set indentation used when marshaling the Policy into JSON.
func (fa *FileAdapter) SetJSONIndent(indent string) {
	fa.jsonIndent = indent
}

// SetFilePerm - allows to set perm of the file the policy is written into.
func (fa *FileAdapter) SetFilePerm(perm FilePerm) {
	fa.filePerm = perm
}

// LoadPolicy - loads and returns policy from file specified when creating FileAdapter.
func (fa *FileAdapter) LoadPolicy() (*restrict.PolicyDefinition, error) {
	data, err := fa.fileHandler.ReadFile(fa.fileName)
	if err != nil {
		return nil, err
	}

	if fa.fileType == JSONFile {
		return fa.createFromJSON(data)
	}

	if fa.fileType == YAMLFile {
		return fa.createFromYAML(data)
	}

	return nil, newFileTypeNotSupportedError(string(fa.fileType))
}

// createFromJSON - helper function for creating the policy from JSON data.
func (fa *FileAdapter) createFromJSON(data []byte) (*restrict.PolicyDefinition, error) {
	var policy *restrict.PolicyDefinition

	if err := fa.jsonHandler.Unmarshal(data, &policy); err != nil {
		return nil, err
	}

	return policy, nil
}

// createFromYAML - helper function for creating the policy from YAML data.
func (fa *FileAdapter) createFromYAML(data []byte) (*restrict.PolicyDefinition, error) {
	var policy *restrict.PolicyDefinition

	if err := fa.yamlHandler.Unmarshal(data, &policy); err != nil {
		return nil, err
	}

	return policy, nil
}

// SavePolicy - saves given policy in file specified when creating FileAdapter.
func (fa *FileAdapter) SavePolicy(policy *restrict.PolicyDefinition) error {
	if fa.fileType == JSONFile {
		return fa.saveJSON(policy)
	}

	if fa.fileType == YAMLFile {
		return fa.saveYAML(policy)
	}

	return newFileTypeNotSupportedError(string(fa.fileType))
}

// saveJSON - helper function for saving policy in JSON format.
func (fa *FileAdapter) saveJSON(policy *restrict.PolicyDefinition) error {
	json, err := fa.jsonHandler.MarshalIndent(policy, "", fa.jsonIndent)
	if err != nil {
		return err
	}

	if err := fa.saveFile(json); err != nil {
		return err
	}

	return nil
}

// saveYAML - helper function for saving policy in YAML format.
func (fa *FileAdapter) saveYAML(policy *restrict.PolicyDefinition) error {
	yaml, err := fa.yamlHandler.Marshal(policy)
	if err != nil {
		return err
	}

	if err := fa.saveFile(yaml); err != nil {
		return err
	}

	return nil
}

// saveFile - saves content to file.
func (fa *FileAdapter) saveFile(content []byte) error {
	return fa.fileHandler.WriteFile(fa.fileName, content, fa.filePerm)
}
