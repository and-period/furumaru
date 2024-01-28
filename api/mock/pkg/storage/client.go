// Code generated by MockGen. DO NOT EDIT.
// Source: client.go

// Package mock_storage is a generated GoMock package.
package mock_storage

import (
	context "context"
	io "io"
	url "net/url"
	reflect "reflect"
	time "time"

	storage "github.com/and-period/furumaru/api/pkg/storage"
	gomock "go.uber.org/mock/gomock"
)

// MockBucket is a mock of Bucket interface.
type MockBucket struct {
	ctrl     *gomock.Controller
	recorder *MockBucketMockRecorder
}

// MockBucketMockRecorder is the mock recorder for MockBucket.
type MockBucketMockRecorder struct {
	mock *MockBucket
}

// NewMockBucket creates a new mock instance.
func NewMockBucket(ctrl *gomock.Controller) *MockBucket {
	mock := &MockBucket{ctrl: ctrl}
	mock.recorder = &MockBucketMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBucket) EXPECT() *MockBucketMockRecorder {
	return m.recorder
}

// Copy mocks base method.
func (m *MockBucket) Copy(ctx context.Context, source, key string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Copy", ctx, source, key)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Copy indicates an expected call of Copy.
func (mr *MockBucketMockRecorder) Copy(ctx, source, key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Copy", reflect.TypeOf((*MockBucket)(nil).Copy), ctx, source, key)
}

// Download mocks base method.
func (m *MockBucket) Download(ctx context.Context, url string) (io.Reader, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Download", ctx, url)
	ret0, _ := ret[0].(io.Reader)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Download indicates an expected call of Download.
func (mr *MockBucketMockRecorder) Download(ctx, url interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Download", reflect.TypeOf((*MockBucket)(nil).Download), ctx, url)
}

// DownloadAndReadAll mocks base method.
func (m *MockBucket) DownloadAndReadAll(ctx context.Context, url string) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DownloadAndReadAll", ctx, url)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DownloadAndReadAll indicates an expected call of DownloadAndReadAll.
func (mr *MockBucketMockRecorder) DownloadAndReadAll(ctx, url interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DownloadAndReadAll", reflect.TypeOf((*MockBucket)(nil).DownloadAndReadAll), ctx, url)
}

// GenerateObjectURL mocks base method.
func (m *MockBucket) GenerateObjectURL(path string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateObjectURL", path)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GenerateObjectURL indicates an expected call of GenerateObjectURL.
func (mr *MockBucketMockRecorder) GenerateObjectURL(path interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateObjectURL", reflect.TypeOf((*MockBucket)(nil).GenerateObjectURL), path)
}

// GeneratePresignUploadURI mocks base method.
func (m *MockBucket) GeneratePresignUploadURI(key string, expiresIn time.Duration) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GeneratePresignUploadURI", key, expiresIn)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GeneratePresignUploadURI indicates an expected call of GeneratePresignUploadURI.
func (mr *MockBucketMockRecorder) GeneratePresignUploadURI(key, expiresIn interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GeneratePresignUploadURI", reflect.TypeOf((*MockBucket)(nil).GeneratePresignUploadURI), key, expiresIn)
}

// GenerateS3URI mocks base method.
func (m *MockBucket) GenerateS3URI(path string) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateS3URI", path)
	ret0, _ := ret[0].(string)
	return ret0
}

// GenerateS3URI indicates an expected call of GenerateS3URI.
func (mr *MockBucketMockRecorder) GenerateS3URI(path interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateS3URI", reflect.TypeOf((*MockBucket)(nil).GenerateS3URI), path)
}

// GetFQDN mocks base method.
func (m *MockBucket) GetFQDN() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFQDN")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetFQDN indicates an expected call of GetFQDN.
func (mr *MockBucketMockRecorder) GetFQDN() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFQDN", reflect.TypeOf((*MockBucket)(nil).GetFQDN))
}

// GetHost mocks base method.
func (m *MockBucket) GetHost() (*url.URL, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetHost")
	ret0, _ := ret[0].(*url.URL)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetHost indicates an expected call of GetHost.
func (mr *MockBucketMockRecorder) GetHost() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetHost", reflect.TypeOf((*MockBucket)(nil).GetHost))
}

// GetMetadata mocks base method.
func (m *MockBucket) GetMetadata(ctx context.Context, key string) (*storage.Metadata, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMetadata", ctx, key)
	ret0, _ := ret[0].(*storage.Metadata)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMetadata indicates an expected call of GetMetadata.
func (mr *MockBucketMockRecorder) GetMetadata(ctx, key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMetadata", reflect.TypeOf((*MockBucket)(nil).GetMetadata), ctx, key)
}

// IsMyHost mocks base method.
func (m *MockBucket) IsMyHost(url string) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsMyHost", url)
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsMyHost indicates an expected call of IsMyHost.
func (mr *MockBucketMockRecorder) IsMyHost(url interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsMyHost", reflect.TypeOf((*MockBucket)(nil).IsMyHost), url)
}

// ReplaceURLToS3URI mocks base method.
func (m *MockBucket) ReplaceURLToS3URI(rawURL string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReplaceURLToS3URI", rawURL)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReplaceURLToS3URI indicates an expected call of ReplaceURLToS3URI.
func (mr *MockBucketMockRecorder) ReplaceURLToS3URI(rawURL interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReplaceURLToS3URI", reflect.TypeOf((*MockBucket)(nil).ReplaceURLToS3URI), rawURL)
}

// Upload mocks base method.
func (m *MockBucket) Upload(ctx context.Context, path string, body io.Reader) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Upload", ctx, path, body)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Upload indicates an expected call of Upload.
func (mr *MockBucketMockRecorder) Upload(ctx, path, body interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Upload", reflect.TypeOf((*MockBucket)(nil).Upload), ctx, path, body)
}
