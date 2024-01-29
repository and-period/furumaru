// Code generated by MockGen. DO NOT EDIT.
// Source: service.go

// Package mock_media is a generated GoMock package.
package mock_media

import (
	context "context"
	reflect "reflect"

	media "github.com/and-period/furumaru/api/internal/media"
	entity "github.com/and-period/furumaru/api/internal/media/entity"
	gomock "go.uber.org/mock/gomock"
)

// MockService is a mock of Service interface.
type MockService struct {
	ctrl     *gomock.Controller
	recorder *MockServiceMockRecorder
}

// MockServiceMockRecorder is the mock recorder for MockService.
type MockServiceMockRecorder struct {
	mock *MockService
}

// NewMockService creates a new mock instance.
func NewMockService(ctrl *gomock.Controller) *MockService {
	mock := &MockService{ctrl: ctrl}
	mock.recorder = &MockServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockService) EXPECT() *MockServiceMockRecorder {
	return m.recorder
}

// ActivateBroadcastMP4 mocks base method.
func (m *MockService) ActivateBroadcastMP4(ctx context.Context, in *media.ActivateBroadcastMP4Input) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ActivateBroadcastMP4", ctx, in)
	ret0, _ := ret[0].(error)
	return ret0
}

// ActivateBroadcastMP4 indicates an expected call of ActivateBroadcastMP4.
func (mr *MockServiceMockRecorder) ActivateBroadcastMP4(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ActivateBroadcastMP4", reflect.TypeOf((*MockService)(nil).ActivateBroadcastMP4), ctx, in)
}

// ActivateBroadcastRTMP mocks base method.
func (m *MockService) ActivateBroadcastRTMP(ctx context.Context, in *media.ActivateBroadcastRTMPInput) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ActivateBroadcastRTMP", ctx, in)
	ret0, _ := ret[0].(error)
	return ret0
}

// ActivateBroadcastRTMP indicates an expected call of ActivateBroadcastRTMP.
func (mr *MockServiceMockRecorder) ActivateBroadcastRTMP(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ActivateBroadcastRTMP", reflect.TypeOf((*MockService)(nil).ActivateBroadcastRTMP), ctx, in)
}

// ActivateBroadcastStaticImage mocks base method.
func (m *MockService) ActivateBroadcastStaticImage(ctx context.Context, in *media.ActivateBroadcastStaticImageInput) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ActivateBroadcastStaticImage", ctx, in)
	ret0, _ := ret[0].(error)
	return ret0
}

// ActivateBroadcastStaticImage indicates an expected call of ActivateBroadcastStaticImage.
func (mr *MockServiceMockRecorder) ActivateBroadcastStaticImage(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ActivateBroadcastStaticImage", reflect.TypeOf((*MockService)(nil).ActivateBroadcastStaticImage), ctx, in)
}

// CreateBroadcast mocks base method.
func (m *MockService) CreateBroadcast(ctx context.Context, in *media.CreateBroadcastInput) (*entity.Broadcast, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateBroadcast", ctx, in)
	ret0, _ := ret[0].(*entity.Broadcast)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateBroadcast indicates an expected call of CreateBroadcast.
func (mr *MockServiceMockRecorder) CreateBroadcast(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateBroadcast", reflect.TypeOf((*MockService)(nil).CreateBroadcast), ctx, in)
}

// CreateBroadcastViewerLog mocks base method.
func (m *MockService) CreateBroadcastViewerLog(ctx context.Context, in *media.CreateBroadcastViewerLogInput) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateBroadcastViewerLog", ctx, in)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateBroadcastViewerLog indicates an expected call of CreateBroadcastViewerLog.
func (mr *MockServiceMockRecorder) CreateBroadcastViewerLog(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateBroadcastViewerLog", reflect.TypeOf((*MockService)(nil).CreateBroadcastViewerLog), ctx, in)
}

// DeactivateBroadcastStaticImage mocks base method.
func (m *MockService) DeactivateBroadcastStaticImage(ctx context.Context, in *media.DeactivateBroadcastStaticImageInput) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeactivateBroadcastStaticImage", ctx, in)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeactivateBroadcastStaticImage indicates an expected call of DeactivateBroadcastStaticImage.
func (mr *MockServiceMockRecorder) DeactivateBroadcastStaticImage(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeactivateBroadcastStaticImage", reflect.TypeOf((*MockService)(nil).DeactivateBroadcastStaticImage), ctx, in)
}

// GetBroadcastArchiveMP4UploadURL mocks base method.
func (m *MockService) GetBroadcastArchiveMP4UploadURL(ctx context.Context, in *media.GenerateBroadcastArchiveMP4UploadInput) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBroadcastArchiveMP4UploadURL", ctx, in)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBroadcastArchiveMP4UploadURL indicates an expected call of GetBroadcastArchiveMP4UploadURL.
func (mr *MockServiceMockRecorder) GetBroadcastArchiveMP4UploadURL(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBroadcastArchiveMP4UploadURL", reflect.TypeOf((*MockService)(nil).GetBroadcastArchiveMP4UploadURL), ctx, in)
}

// GetBroadcastByScheduleID mocks base method.
func (m *MockService) GetBroadcastByScheduleID(ctx context.Context, in *media.GetBroadcastByScheduleIDInput) (*entity.Broadcast, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBroadcastByScheduleID", ctx, in)
	ret0, _ := ret[0].(*entity.Broadcast)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBroadcastByScheduleID indicates an expected call of GetBroadcastByScheduleID.
func (mr *MockServiceMockRecorder) GetBroadcastByScheduleID(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBroadcastByScheduleID", reflect.TypeOf((*MockService)(nil).GetBroadcastByScheduleID), ctx, in)
}

// GetBroadcastLiveMP4UploadURL mocks base method.
func (m *MockService) GetBroadcastLiveMP4UploadURL(ctx context.Context, in *media.GenerateUploadURLInput) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBroadcastLiveMP4UploadURL", ctx, in)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBroadcastLiveMP4UploadURL indicates an expected call of GetBroadcastLiveMP4UploadURL.
func (mr *MockServiceMockRecorder) GetBroadcastLiveMP4UploadURL(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBroadcastLiveMP4UploadURL", reflect.TypeOf((*MockService)(nil).GetBroadcastLiveMP4UploadURL), ctx, in)
}

// GetCoordinatorBonusVideoUploadURL mocks base method.
func (m *MockService) GetCoordinatorBonusVideoUploadURL(ctx context.Context, in *media.GenerateUploadURLInput) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCoordinatorBonusVideoUploadURL", ctx, in)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCoordinatorBonusVideoUploadURL indicates an expected call of GetCoordinatorBonusVideoUploadURL.
func (mr *MockServiceMockRecorder) GetCoordinatorBonusVideoUploadURL(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCoordinatorBonusVideoUploadURL", reflect.TypeOf((*MockService)(nil).GetCoordinatorBonusVideoUploadURL), ctx, in)
}

// GetCoordinatorHeaderUploadURL mocks base method.
func (m *MockService) GetCoordinatorHeaderUploadURL(ctx context.Context, in *media.GenerateUploadURLInput) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCoordinatorHeaderUploadURL", ctx, in)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCoordinatorHeaderUploadURL indicates an expected call of GetCoordinatorHeaderUploadURL.
func (mr *MockServiceMockRecorder) GetCoordinatorHeaderUploadURL(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCoordinatorHeaderUploadURL", reflect.TypeOf((*MockService)(nil).GetCoordinatorHeaderUploadURL), ctx, in)
}

// GetCoordinatorPromotionVideoUploadURL mocks base method.
func (m *MockService) GetCoordinatorPromotionVideoUploadURL(ctx context.Context, in *media.GenerateUploadURLInput) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCoordinatorPromotionVideoUploadURL", ctx, in)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCoordinatorPromotionVideoUploadURL indicates an expected call of GetCoordinatorPromotionVideoUploadURL.
func (mr *MockServiceMockRecorder) GetCoordinatorPromotionVideoUploadURL(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCoordinatorPromotionVideoUploadURL", reflect.TypeOf((*MockService)(nil).GetCoordinatorPromotionVideoUploadURL), ctx, in)
}

// GetCoordinatorThumbnailUploadURL mocks base method.
func (m *MockService) GetCoordinatorThumbnailUploadURL(ctx context.Context, in *media.GenerateUploadURLInput) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCoordinatorThumbnailUploadURL", ctx, in)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCoordinatorThumbnailUploadURL indicates an expected call of GetCoordinatorThumbnailUploadURL.
func (mr *MockServiceMockRecorder) GetCoordinatorThumbnailUploadURL(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCoordinatorThumbnailUploadURL", reflect.TypeOf((*MockService)(nil).GetCoordinatorThumbnailUploadURL), ctx, in)
}

// GetProducerBonusVideoUploadURL mocks base method.
func (m *MockService) GetProducerBonusVideoUploadURL(ctx context.Context, in *media.GenerateUploadURLInput) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProducerBonusVideoUploadURL", ctx, in)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProducerBonusVideoUploadURL indicates an expected call of GetProducerBonusVideoUploadURL.
func (mr *MockServiceMockRecorder) GetProducerBonusVideoUploadURL(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProducerBonusVideoUploadURL", reflect.TypeOf((*MockService)(nil).GetProducerBonusVideoUploadURL), ctx, in)
}

// GetProducerHeaderUploadURL mocks base method.
func (m *MockService) GetProducerHeaderUploadURL(ctx context.Context, in *media.GenerateUploadURLInput) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProducerHeaderUploadURL", ctx, in)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProducerHeaderUploadURL indicates an expected call of GetProducerHeaderUploadURL.
func (mr *MockServiceMockRecorder) GetProducerHeaderUploadURL(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProducerHeaderUploadURL", reflect.TypeOf((*MockService)(nil).GetProducerHeaderUploadURL), ctx, in)
}

// GetProducerPromotionVideoUploadURL mocks base method.
func (m *MockService) GetProducerPromotionVideoUploadURL(ctx context.Context, in *media.GenerateUploadURLInput) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProducerPromotionVideoUploadURL", ctx, in)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProducerPromotionVideoUploadURL indicates an expected call of GetProducerPromotionVideoUploadURL.
func (mr *MockServiceMockRecorder) GetProducerPromotionVideoUploadURL(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProducerPromotionVideoUploadURL", reflect.TypeOf((*MockService)(nil).GetProducerPromotionVideoUploadURL), ctx, in)
}

// GetProducerThumbnailUploadURL mocks base method.
func (m *MockService) GetProducerThumbnailUploadURL(ctx context.Context, in *media.GenerateUploadURLInput) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProducerThumbnailUploadURL", ctx, in)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProducerThumbnailUploadURL indicates an expected call of GetProducerThumbnailUploadURL.
func (mr *MockServiceMockRecorder) GetProducerThumbnailUploadURL(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProducerThumbnailUploadURL", reflect.TypeOf((*MockService)(nil).GetProducerThumbnailUploadURL), ctx, in)
}

// GetProductMediaImageUploadURL mocks base method.
func (m *MockService) GetProductMediaImageUploadURL(ctx context.Context, in *media.GenerateUploadURLInput) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProductMediaImageUploadURL", ctx, in)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProductMediaImageUploadURL indicates an expected call of GetProductMediaImageUploadURL.
func (mr *MockServiceMockRecorder) GetProductMediaImageUploadURL(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProductMediaImageUploadURL", reflect.TypeOf((*MockService)(nil).GetProductMediaImageUploadURL), ctx, in)
}

// GetProductMediaVideoUploadURL mocks base method.
func (m *MockService) GetProductMediaVideoUploadURL(ctx context.Context, in *media.GenerateUploadURLInput) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProductMediaVideoUploadURL", ctx, in)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProductMediaVideoUploadURL indicates an expected call of GetProductMediaVideoUploadURL.
func (mr *MockServiceMockRecorder) GetProductMediaVideoUploadURL(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProductMediaVideoUploadURL", reflect.TypeOf((*MockService)(nil).GetProductMediaVideoUploadURL), ctx, in)
}

// GetProductTypeIconUploadURL mocks base method.
func (m *MockService) GetProductTypeIconUploadURL(ctx context.Context, in *media.GenerateUploadURLInput) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProductTypeIconUploadURL", ctx, in)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProductTypeIconUploadURL indicates an expected call of GetProductTypeIconUploadURL.
func (mr *MockServiceMockRecorder) GetProductTypeIconUploadURL(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProductTypeIconUploadURL", reflect.TypeOf((*MockService)(nil).GetProductTypeIconUploadURL), ctx, in)
}

// GetScheduleImageUploadURL mocks base method.
func (m *MockService) GetScheduleImageUploadURL(ctx context.Context, in *media.GenerateUploadURLInput) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetScheduleImageUploadURL", ctx, in)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetScheduleImageUploadURL indicates an expected call of GetScheduleImageUploadURL.
func (mr *MockServiceMockRecorder) GetScheduleImageUploadURL(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetScheduleImageUploadURL", reflect.TypeOf((*MockService)(nil).GetScheduleImageUploadURL), ctx, in)
}

// GetScheduleOpeningVideoUploadURL mocks base method.
func (m *MockService) GetScheduleOpeningVideoUploadURL(ctx context.Context, in *media.GenerateUploadURLInput) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetScheduleOpeningVideoUploadURL", ctx, in)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetScheduleOpeningVideoUploadURL indicates an expected call of GetScheduleOpeningVideoUploadURL.
func (mr *MockServiceMockRecorder) GetScheduleOpeningVideoUploadURL(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetScheduleOpeningVideoUploadURL", reflect.TypeOf((*MockService)(nil).GetScheduleOpeningVideoUploadURL), ctx, in)
}

// GetScheduleThumbnailUploadURL mocks base method.
func (m *MockService) GetScheduleThumbnailUploadURL(ctx context.Context, in *media.GenerateUploadURLInput) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetScheduleThumbnailUploadURL", ctx, in)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetScheduleThumbnailUploadURL indicates an expected call of GetScheduleThumbnailUploadURL.
func (mr *MockServiceMockRecorder) GetScheduleThumbnailUploadURL(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetScheduleThumbnailUploadURL", reflect.TypeOf((*MockService)(nil).GetScheduleThumbnailUploadURL), ctx, in)
}

// GetUploadEvent mocks base method.
func (m *MockService) GetUploadEvent(ctx context.Context, in *media.GetUploadEventInput) (*entity.UploadEvent, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUploadEvent", ctx, in)
	ret0, _ := ret[0].(*entity.UploadEvent)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUploadEvent indicates an expected call of GetUploadEvent.
func (mr *MockServiceMockRecorder) GetUploadEvent(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUploadEvent", reflect.TypeOf((*MockService)(nil).GetUploadEvent), ctx, in)
}

// GetUserThumbnailUploadURL mocks base method.
func (m *MockService) GetUserThumbnailUploadURL(ctx context.Context, in *media.GenerateUploadURLInput) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserThumbnailUploadURL", ctx, in)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserThumbnailUploadURL indicates an expected call of GetUserThumbnailUploadURL.
func (mr *MockServiceMockRecorder) GetUserThumbnailUploadURL(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserThumbnailUploadURL", reflect.TypeOf((*MockService)(nil).GetUserThumbnailUploadURL), ctx, in)
}

// ListBroadcasts mocks base method.
func (m *MockService) ListBroadcasts(ctx context.Context, in *media.ListBroadcastsInput) (entity.Broadcasts, int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListBroadcasts", ctx, in)
	ret0, _ := ret[0].(entity.Broadcasts)
	ret1, _ := ret[1].(int64)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ListBroadcasts indicates an expected call of ListBroadcasts.
func (mr *MockServiceMockRecorder) ListBroadcasts(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListBroadcasts", reflect.TypeOf((*MockService)(nil).ListBroadcasts), ctx, in)
}

// PauseBroadcast mocks base method.
func (m *MockService) PauseBroadcast(ctx context.Context, in *media.PauseBroadcastInput) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PauseBroadcast", ctx, in)
	ret0, _ := ret[0].(error)
	return ret0
}

// PauseBroadcast indicates an expected call of PauseBroadcast.
func (mr *MockServiceMockRecorder) PauseBroadcast(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PauseBroadcast", reflect.TypeOf((*MockService)(nil).PauseBroadcast), ctx, in)
}

// ResizeCoordinatorHeader mocks base method.
func (m *MockService) ResizeCoordinatorHeader(ctx context.Context, in *media.ResizeFileInput) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ResizeCoordinatorHeader", ctx, in)
	ret0, _ := ret[0].(error)
	return ret0
}

// ResizeCoordinatorHeader indicates an expected call of ResizeCoordinatorHeader.
func (mr *MockServiceMockRecorder) ResizeCoordinatorHeader(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ResizeCoordinatorHeader", reflect.TypeOf((*MockService)(nil).ResizeCoordinatorHeader), ctx, in)
}

// ResizeCoordinatorThumbnail mocks base method.
func (m *MockService) ResizeCoordinatorThumbnail(ctx context.Context, in *media.ResizeFileInput) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ResizeCoordinatorThumbnail", ctx, in)
	ret0, _ := ret[0].(error)
	return ret0
}

// ResizeCoordinatorThumbnail indicates an expected call of ResizeCoordinatorThumbnail.
func (mr *MockServiceMockRecorder) ResizeCoordinatorThumbnail(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ResizeCoordinatorThumbnail", reflect.TypeOf((*MockService)(nil).ResizeCoordinatorThumbnail), ctx, in)
}

// ResizeProducerHeader mocks base method.
func (m *MockService) ResizeProducerHeader(ctx context.Context, in *media.ResizeFileInput) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ResizeProducerHeader", ctx, in)
	ret0, _ := ret[0].(error)
	return ret0
}

// ResizeProducerHeader indicates an expected call of ResizeProducerHeader.
func (mr *MockServiceMockRecorder) ResizeProducerHeader(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ResizeProducerHeader", reflect.TypeOf((*MockService)(nil).ResizeProducerHeader), ctx, in)
}

// ResizeProducerThumbnail mocks base method.
func (m *MockService) ResizeProducerThumbnail(ctx context.Context, in *media.ResizeFileInput) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ResizeProducerThumbnail", ctx, in)
	ret0, _ := ret[0].(error)
	return ret0
}

// ResizeProducerThumbnail indicates an expected call of ResizeProducerThumbnail.
func (mr *MockServiceMockRecorder) ResizeProducerThumbnail(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ResizeProducerThumbnail", reflect.TypeOf((*MockService)(nil).ResizeProducerThumbnail), ctx, in)
}

// ResizeProductMedia mocks base method.
func (m *MockService) ResizeProductMedia(ctx context.Context, in *media.ResizeFileInput) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ResizeProductMedia", ctx, in)
	ret0, _ := ret[0].(error)
	return ret0
}

// ResizeProductMedia indicates an expected call of ResizeProductMedia.
func (mr *MockServiceMockRecorder) ResizeProductMedia(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ResizeProductMedia", reflect.TypeOf((*MockService)(nil).ResizeProductMedia), ctx, in)
}

// ResizeProductTypeIcon mocks base method.
func (m *MockService) ResizeProductTypeIcon(ctx context.Context, in *media.ResizeFileInput) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ResizeProductTypeIcon", ctx, in)
	ret0, _ := ret[0].(error)
	return ret0
}

// ResizeProductTypeIcon indicates an expected call of ResizeProductTypeIcon.
func (mr *MockServiceMockRecorder) ResizeProductTypeIcon(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ResizeProductTypeIcon", reflect.TypeOf((*MockService)(nil).ResizeProductTypeIcon), ctx, in)
}

// ResizeScheduleThumbnail mocks base method.
func (m *MockService) ResizeScheduleThumbnail(ctx context.Context, in *media.ResizeFileInput) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ResizeScheduleThumbnail", ctx, in)
	ret0, _ := ret[0].(error)
	return ret0
}

// ResizeScheduleThumbnail indicates an expected call of ResizeScheduleThumbnail.
func (mr *MockServiceMockRecorder) ResizeScheduleThumbnail(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ResizeScheduleThumbnail", reflect.TypeOf((*MockService)(nil).ResizeScheduleThumbnail), ctx, in)
}

// ResizeUserThumbnail mocks base method.
func (m *MockService) ResizeUserThumbnail(ctx context.Context, in *media.ResizeFileInput) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ResizeUserThumbnail", ctx, in)
	ret0, _ := ret[0].(error)
	return ret0
}

// ResizeUserThumbnail indicates an expected call of ResizeUserThumbnail.
func (mr *MockServiceMockRecorder) ResizeUserThumbnail(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ResizeUserThumbnail", reflect.TypeOf((*MockService)(nil).ResizeUserThumbnail), ctx, in)
}

// UnpauseBroadcast mocks base method.
func (m *MockService) UnpauseBroadcast(ctx context.Context, in *media.UnpauseBroadcastInput) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UnpauseBroadcast", ctx, in)
	ret0, _ := ret[0].(error)
	return ret0
}

// UnpauseBroadcast indicates an expected call of UnpauseBroadcast.
func (mr *MockServiceMockRecorder) UnpauseBroadcast(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UnpauseBroadcast", reflect.TypeOf((*MockService)(nil).UnpauseBroadcast), ctx, in)
}

// UpdateBroadcastArchive mocks base method.
func (m *MockService) UpdateBroadcastArchive(ctx context.Context, in *media.UpdateBroadcastArchiveInput) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateBroadcastArchive", ctx, in)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateBroadcastArchive indicates an expected call of UpdateBroadcastArchive.
func (mr *MockServiceMockRecorder) UpdateBroadcastArchive(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateBroadcastArchive", reflect.TypeOf((*MockService)(nil).UpdateBroadcastArchive), ctx, in)
}
