package adapters

import (
	"errors"
	"testing"

	"github.com/el-mike/restrict"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type fileHandlerMock struct {
	mock.Mock
}

func (m *fileHandlerMock) ReadFile(name string) ([]byte, error) {
	args := m.Called(name)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]byte), args.Error(1)
}

func (m *fileHandlerMock) WriteFile(name string, data []byte, perm FilePerm) error {
	args := m.Called(name, data, perm)

	return args.Error(0)
}

type jsonHandlerMock struct {
	mock.Mock
}

func (m *jsonHandlerMock) Unmarshal(data []byte, v interface{}) error {
	args := m.Called(data, v)

	return args.Error(0)
}

func (m *jsonHandlerMock) MarshalIndent(v interface{}, prefix, indent string) ([]byte, error) {
	args := m.Called(v, prefix, indent)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]byte), args.Error(1)
}

type yamlHandlerMock struct {
	mock.Mock
}

func (m *yamlHandlerMock) Unmarshal(in []byte, out interface{}) error {
	args := m.Called(in, out)

	return args.Error(0)
}

func (m *yamlHandlerMock) Marshal(in interface{}) ([]byte, error) {
	args := m.Called(in)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]byte), args.Error(1)
}

type fileAdapterSuite struct {
	suite.Suite

	testFileName string
	testError    error
}

func (s *fileAdapterSuite) SetupSuite() {
	s.testFileName = "testFile.json"
	s.testError = errors.New("testError")
}

func TestFileAdapterSuite(t *testing.T) {
	suite.Run(t, new(fileAdapterSuite))
}

func (s *fileAdapterSuite) TestNewFileAdapter() {
	testFileType := JSONFile

	adapter := NewFileAdapter(s.testFileName, testFileType)

	assert.NotNil(s.T(), adapter)
	assert.Equal(s.T(), adapter.fileName, s.testFileName)
	assert.Equal(s.T(), adapter.fileType, testFileType)

	assert.IsType(s.T(), adapter.fileHandler, new(defaultFileHandler))
	assert.IsType(s.T(), adapter.jsonHandler, new(defaultJSONHandler))
	assert.IsType(s.T(), adapter.yamlHandler, new(defaultYAMLHandler))

	testFileType = YAMLFile

	adapter = NewFileAdapter(s.testFileName, testFileType)

	assert.NotNil(s.T(), adapter)
	assert.Equal(s.T(), adapter.fileType, testFileType)
}

func (s *fileAdapterSuite) TestSetJSONIndent() {
	adapter := NewFileAdapter(s.testFileName, JSONFile)

	assert.Equal(s.T(), adapter.jsonIndent, defaultJSONIndent)

	testIndent := "  "

	adapter.SetJSONIndent(testIndent)

	assert.Equal(s.T(), adapter.jsonIndent, testIndent)
}

func (s *fileAdapterSuite) TestSetFilePerm() {
	adapter := NewFileAdapter(s.testFileName, JSONFile)

	assert.Equal(s.T(), adapter.filePerm, defaultFilePerm)

	testPerm := FilePerm(0777)

	adapter.SetFilePerm(testPerm)

	assert.Equal(s.T(), adapter.filePerm, testPerm)
}

func (s *fileAdapterSuite) TestLoadPolicy_ReadFile() {
	// Load with working fileHadler.
	testData := []byte("{}")

	workingFileHandler := new(fileHandlerMock)
	workingFileHandler.On(
		"ReadFile",
		mock.Anything,
	).Return(testData, nil)

	adapter := NewFileAdapter(s.testFileName, JSONFile)

	adapter.fileHandler = workingFileHandler

	_, err := adapter.LoadPolicy()

	assert.Nil(s.T(), err)
	workingFileHandler.AssertNumberOfCalls(s.T(), "ReadFile", 1)
	workingFileHandler.AssertCalled(s.T(), "ReadFile", s.testFileName)

	// Load with incorrect fileType
	adapter.fileType = "incorrectFileType"
	_, err = adapter.LoadPolicy()

	assert.Error(s.T(), err)
	assert.IsType(s.T(), new(FileTypeNotSupportedError), err)

	// Load with failing fileHandler

	failingFileHandler := new(fileHandlerMock)
	failingFileHandler.On(
		"ReadFile",
		mock.Anything,
	).Return(nil, s.testError)

	adapter.fileHandler = failingFileHandler

	_, err = adapter.LoadPolicy()

	assert.Error(s.T(), err)
	assert.Equal(s.T(), err, s.testError)
}

func (s *fileAdapterSuite) TestLoadPolicy_JSONFile() {
	// Load with working jsonHandler
	testData := []byte(getBasicPolicyJSONString())

	testFileHandler := new(fileHandlerMock)
	testFileHandler.On(
		"ReadFile",
		mock.Anything,
	).Return(testData, nil)

	workingJSONHandler := new(jsonHandlerMock)
	workingJSONHandler.On(
		"Unmarshal",
		mock.Anything,
		mock.Anything,
	).Return(nil)

	workingYAMLHandler := new(yamlHandlerMock)
	workingYAMLHandler.On(
		"Marshal",
		mock.Anything,
	).Return(testData, nil)

	adapter := NewFileAdapter(s.testFileName, JSONFile)

	adapter.fileHandler = testFileHandler
	adapter.jsonHandler = workingJSONHandler

	policy, err := adapter.LoadPolicy()

	assert.Nil(s.T(), err)
	assert.IsType(s.T(), policy, new(restrict.PolicyDefinition))
	workingJSONHandler.AssertNumberOfCalls(s.T(), "Unmarshal", 1)
	workingJSONHandler.AssertCalled(s.T(), "Unmarshal", testData, mock.Anything)

	// Load with failing jsonHandler
	failingJSONHandler := new(jsonHandlerMock)
	failingJSONHandler.On(
		"Unmarshal",
		mock.Anything,
		mock.Anything,
	).Return(s.testError)

	adapter.jsonHandler = failingJSONHandler

	policy, err = adapter.LoadPolicy()

	assert.Nil(s.T(), policy)
	assert.Error(s.T(), err)
}

func (s *fileAdapterSuite) TestLoadPolicy_JSONReal() {
	testData := []byte(getBasicPolicyJSONString())

	testFileHandler := new(fileHandlerMock)
	testFileHandler.On(
		"ReadFile",
		mock.Anything,
	).Return(testData, nil)

	adapter := NewFileAdapter("test.json", JSONFile)

	adapter.fileHandler = testFileHandler

	policy, err := adapter.LoadPolicy()

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 1, len(policy.Roles))
}

func (s *fileAdapterSuite) TestLoadPolicy_YAMLFile() {
	// Load with working yamlHandler
	testData := []byte(getBasicPolicyYAMLString())

	testFileHandler := new(fileHandlerMock)
	testFileHandler.On(
		"ReadFile",
		mock.Anything,
	).Return(testData, nil)

	workingYAMLHandler := new(yamlHandlerMock)
	workingYAMLHandler.On(
		"Unmarshal",
		mock.Anything,
		mock.Anything,
	).Return(nil)

	adapter := NewFileAdapter(s.testFileName, YAMLFile)

	adapter.fileHandler = testFileHandler
	adapter.yamlHandler = workingYAMLHandler

	policy, err := adapter.LoadPolicy()

	assert.Nil(s.T(), err)
	assert.IsType(s.T(), policy, new(restrict.PolicyDefinition))
	workingYAMLHandler.AssertNumberOfCalls(s.T(), "Unmarshal", 1)
	workingYAMLHandler.AssertCalled(s.T(), "Unmarshal", testData, mock.Anything)

	// Load with failing jsonHandler
	failingYAMLHandler := new(yamlHandlerMock)
	failingYAMLHandler.On(
		"Unmarshal",
		mock.Anything,
		mock.Anything,
	).Return(s.testError)

	adapter.yamlHandler = failingYAMLHandler

	policy, err = adapter.LoadPolicy()

	assert.Nil(s.T(), policy)
	assert.Error(s.T(), err)
}

func (s *fileAdapterSuite) TestLoadPolicy_YAMLReal() {
	testData := []byte(getBasicPolicyYAMLString())

	testFileHandler := new(fileHandlerMock)
	testFileHandler.On(
		"ReadFile",
		mock.Anything,
	).Return(testData, nil)

	adapter := NewFileAdapter("test.yml", YAMLFile)

	adapter.fileHandler = testFileHandler

	policy, err := adapter.LoadPolicy()

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 1, len(policy.Roles))
}

func (s *fileAdapterSuite) TestSavePolicy() {
	testPolicy := getBasicPolicy()

	adapter := NewFileAdapter(s.testFileName, "incorrectFileType")

	err := adapter.SavePolicy(testPolicy)

	assert.NotNil(s.T(), err)
	assert.IsType(s.T(), err, new(FileTypeNotSupportedError))
}

func (s *fileAdapterSuite) TestSavePolicy_WriteFile() {
	// Write with working fileHandler
	testJSONData := []byte(getBasicPolicyJSONString())
	testYAMLData := []byte(getBasicPolicyYAMLString())
	testPolicy := getBasicPolicy()

	workingFileHandler := new(fileHandlerMock)
	workingFileHandler.On(
		"WriteFile",
		mock.Anything,
		mock.Anything,
		mock.Anything,
	).Return(nil)

	workingJSONHandler := new(jsonHandlerMock)
	workingJSONHandler.On(
		"MarshalIndent",
		mock.Anything,
		mock.Anything,
		mock.Anything,
	).Return(testJSONData, nil)

	workingYAMLHandler := new(yamlHandlerMock)
	workingYAMLHandler.On(
		"Marshal",
		mock.Anything,
	).Return(testYAMLData, nil)

	adapter := NewFileAdapter(s.testFileName, JSONFile)

	adapter.fileHandler = workingFileHandler
	adapter.jsonHandler = workingJSONHandler

	err := adapter.SavePolicy(testPolicy)

	assert.Nil(s.T(), err)
	workingFileHandler.AssertNumberOfCalls(s.T(), "WriteFile", 1)
	workingFileHandler.AssertCalled(s.T(), "WriteFile", s.testFileName, mock.Anything, defaultFilePerm)

	// Write with failing fileHandler for JSON
	failingFileHandler := new(fileHandlerMock)
	failingFileHandler.On(
		"WriteFile",
		mock.Anything,
		mock.Anything,
		mock.Anything,
	).Return(s.testError)

	adapter.fileHandler = failingFileHandler

	err = adapter.SavePolicy(testPolicy)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), err, s.testError)

	// Write with failing fileHandler for YAML
	adapter = NewFileAdapter(s.testFileName, YAMLFile)

	adapter.fileHandler = failingFileHandler
	adapter.yamlHandler = workingYAMLHandler

	err = adapter.SavePolicy(testPolicy)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), err, s.testError)
}

func (s *fileAdapterSuite) TestSavePolicy_JSONFile() {
	// Save with working jsonHandler
	testData := []byte(getBasicPolicyJSONString())
	testPolicy := getBasicPolicy()

	workingFileHandler := new(fileHandlerMock)
	workingFileHandler.On(
		"WriteFile",
		mock.Anything,
		mock.Anything,
		mock.Anything,
	).Return(nil)

	workingJSONHandler := new(jsonHandlerMock)
	workingJSONHandler.On(
		"MarshalIndent",
		mock.Anything,
		mock.Anything,
		mock.Anything,
	).Return(testData, nil)

	adapter := NewFileAdapter(s.testFileName, JSONFile)

	adapter.fileHandler = workingFileHandler
	adapter.jsonHandler = workingJSONHandler

	err := adapter.SavePolicy(testPolicy)

	assert.Nil(s.T(), err)
	workingJSONHandler.AssertNumberOfCalls(s.T(), "MarshalIndent", 1)
	workingJSONHandler.AssertCalled(s.T(), "MarshalIndent", testPolicy, "", defaultJSONIndent)

	workingFileHandler.AssertNumberOfCalls(s.T(), "WriteFile", 1)
	workingFileHandler.AssertCalled(s.T(), "WriteFile", s.testFileName, testData, defaultFilePerm)

	// Save with failing jsonHandler
	failingJSONHandler := new(jsonHandlerMock)
	failingJSONHandler.On(
		"MarshalIndent",
		mock.Anything,
		mock.Anything,
		mock.Anything,
	).Return([]byte{}, s.testError)

	adapter.jsonHandler = failingJSONHandler

	err = adapter.SavePolicy(testPolicy)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), err, s.testError)

	failingJSONHandler.AssertNumberOfCalls(s.T(), "MarshalIndent", 1)
	// Since we are reusing workingFileHandler, 1 means that it was not called
	// again if jsonHandler returned error.
	workingFileHandler.AssertNumberOfCalls(s.T(), "WriteFile", 1)
}

func (s *fileAdapterSuite) TestSavePolicy_YAMLFile() {
	// Save with working yamlHandler
	testData := []byte(getBasicPolicyYAMLString())
	testPolicy := getBasicPolicy()

	workingFileHandler := new(fileHandlerMock)
	workingFileHandler.On(
		"WriteFile",
		mock.Anything,
		mock.Anything,
		mock.Anything,
	).Return(nil)

	workingYAMLHandler := new(yamlHandlerMock)
	workingYAMLHandler.On(
		"Marshal",
		mock.Anything,
	).Return(testData, nil)

	adapter := NewFileAdapter(s.testFileName, YAMLFile)

	adapter.fileHandler = workingFileHandler
	adapter.yamlHandler = workingYAMLHandler

	err := adapter.SavePolicy(testPolicy)

	assert.Nil(s.T(), err)
	workingYAMLHandler.AssertNumberOfCalls(s.T(), "Marshal", 1)
	workingYAMLHandler.AssertCalled(s.T(), "Marshal", testPolicy)

	workingFileHandler.AssertNumberOfCalls(s.T(), "WriteFile", 1)
	workingFileHandler.AssertCalled(s.T(), "WriteFile", s.testFileName, testData, defaultFilePerm)

	// Save with failing yamlHandler
	failingYAMLHandler := new(yamlHandlerMock)
	failingYAMLHandler.On(
		"Marshal",
		mock.Anything,
	).Return([]byte{}, s.testError)

	adapter.yamlHandler = failingYAMLHandler

	err = adapter.SavePolicy(testPolicy)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), err, s.testError)

	failingYAMLHandler.AssertNumberOfCalls(s.T(), "Marshal", 1)
	// Since we are reusing workingFileHandler, 1 means that it was not called
	// again if yamlHandler returned error.
	workingFileHandler.AssertNumberOfCalls(s.T(), "WriteFile", 1)
}
