package adapters

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/el-Mike/restrict"
	"gopkg.in/yaml.v3"
)

// AllowedFileType - alias type for describing allowed file types.
type AllowedFileType string

const (
	// JSONFile - JSON file type.
	JSONFile AllowedFileType = "JSONFile"
	// YAMLFile - YAML file type.
	YAMLFile AllowedFileType = "YAMLFile"
)

// FileAdapter - policy storage adapter, for handling file storage.
type FileAdapter struct {
	FileName string
	FileType AllowedFileType
}

// NewFileAdapter - returns new FileAdapter instance.
func NewFileAdapter(fileName string, fileType AllowedFileType) *FileAdapter {
	return &FileAdapter{
		FileName: fileName,
		FileType: fileType,
	}
}

// LoadPolicy - loads and returns policy from file specified when creating FileAdapter.
func (fa *FileAdapter) LoadPolicy() (*restrict.PolicyDefinition, error) {
	data, err := os.ReadFile(fa.FileName)
	if err != nil {
		return nil, err
	}

	if fa.FileType == JSONFile {
		return fa.createFromJSON(data)
	}

	if fa.FileType == YAMLFile {
		return fa.createFromYAML(data)
	}

	return nil, errors.New("File type not supported!")
}

// createFromJSON - helper function for creating the policy from JSON data.
func (fa *FileAdapter) createFromJSON(data []byte) (*restrict.PolicyDefinition, error) {
	var policy *restrict.PolicyDefinition

	if err := json.Unmarshal(data, &policy); err != nil {
		return nil, err
	}

	return policy, nil
}

// createFromYAML - helper function for creating the policy from YAML data.
func (fa *FileAdapter) createFromYAML(data []byte) (*restrict.PolicyDefinition, error) {
	var policy *restrict.PolicyDefinition

	if err := yaml.Unmarshal(data, &policy); err != nil {
		return nil, err
	}

	return policy, nil
}

// SavePolicy - saves given policy in file specified when creating FileAdapter.
func (fa *FileAdapter) SavePolicy(policy *restrict.PolicyDefinition) error {
	if fa.FileType == JSONFile {
		return fa.saveJSON(policy)
	}

	if fa.FileType == YAMLFile {
		return fa.saveYAML(policy)
	}

	return errors.New("File type not supported!")
}

// saveJSON - helper function for saving policy in JSON format.
func (fa *FileAdapter) saveJSON(policy *restrict.PolicyDefinition) error {
	json, err := json.MarshalIndent(policy, "", "\t")
	if err != nil {
		return err
	}

	if err := os.WriteFile(fa.FileName, json, 0644); err != nil {
		return err
	}

	return nil
}

// saveYAML - helper function for saving policy in YAML format.
func (fa *FileAdapter) saveYAML(policy *restrict.PolicyDefinition) error {
	yaml, err := yaml.Marshal(policy)
	if err != nil {
		return err
	}

	if err := os.WriteFile(fa.FileName, yaml, 0644); err != nil {
		return err
	}

	return nil
}
