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

// GenerateCoordinatorBonusVideo mocks base method.
func (m *MockService) GenerateCoordinatorBonusVideo(ctx context.Context, in *media.GenerateFileInput) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateCoordinatorBonusVideo", ctx, in)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GenerateCoordinatorBonusVideo indicates an expected call of GenerateCoordinatorBonusVideo.
func (mr *MockServiceMockRecorder) GenerateCoordinatorBonusVideo(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateCoordinatorBonusVideo", reflect.TypeOf((*MockService)(nil).GenerateCoordinatorBonusVideo), ctx, in)
}

// GenerateCoordinatorHeader mocks base method.
func (m *MockService) GenerateCoordinatorHeader(ctx context.Context, in *media.GenerateFileInput) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateCoordinatorHeader", ctx, in)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GenerateCoordinatorHeader indicates an expected call of GenerateCoordinatorHeader.
func (mr *MockServiceMockRecorder) GenerateCoordinatorHeader(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateCoordinatorHeader", reflect.TypeOf((*MockService)(nil).GenerateCoordinatorHeader), ctx, in)
}

// GenerateCoordinatorPromotionVideo mocks base method.
func (m *MockService) GenerateCoordinatorPromotionVideo(ctx context.Context, in *media.GenerateFileInput) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateCoordinatorPromotionVideo", ctx, in)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GenerateCoordinatorPromotionVideo indicates an expected call of GenerateCoordinatorPromotionVideo.
func (mr *MockServiceMockRecorder) GenerateCoordinatorPromotionVideo(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateCoordinatorPromotionVideo", reflect.TypeOf((*MockService)(nil).GenerateCoordinatorPromotionVideo), ctx, in)
}

// GenerateCoordinatorThumbnail mocks base method.
func (m *MockService) GenerateCoordinatorThumbnail(ctx context.Context, in *media.GenerateFileInput) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateCoordinatorThumbnail", ctx, in)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GenerateCoordinatorThumbnail indicates an expected call of GenerateCoordinatorThumbnail.
func (mr *MockServiceMockRecorder) GenerateCoordinatorThumbnail(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateCoordinatorThumbnail", reflect.TypeOf((*MockService)(nil).GenerateCoordinatorThumbnail), ctx, in)
}

// GenerateProducerBonusVideo mocks base method.
func (m *MockService) GenerateProducerBonusVideo(ctx context.Context, in *media.GenerateFileInput) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateProducerBonusVideo", ctx, in)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GenerateProducerBonusVideo indicates an expected call of GenerateProducerBonusVideo.
func (mr *MockServiceMockRecorder) GenerateProducerBonusVideo(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateProducerBonusVideo", reflect.TypeOf((*MockService)(nil).GenerateProducerBonusVideo), ctx, in)
}

// GenerateProducerHeader mocks base method.
func (m *MockService) GenerateProducerHeader(ctx context.Context, in *media.GenerateFileInput) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateProducerHeader", ctx, in)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GenerateProducerHeader indicates an expected call of GenerateProducerHeader.
func (mr *MockServiceMockRecorder) GenerateProducerHeader(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateProducerHeader", reflect.TypeOf((*MockService)(nil).GenerateProducerHeader), ctx, in)
}

// GenerateProducerPromotionVideo mocks base method.
func (m *MockService) GenerateProducerPromotionVideo(ctx context.Context, in *media.GenerateFileInput) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateProducerPromotionVideo", ctx, in)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GenerateProducerPromotionVideo indicates an expected call of GenerateProducerPromotionVideo.
func (mr *MockServiceMockRecorder) GenerateProducerPromotionVideo(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateProducerPromotionVideo", reflect.TypeOf((*MockService)(nil).GenerateProducerPromotionVideo), ctx, in)
}

// GenerateProducerThumbnail mocks base method.
func (m *MockService) GenerateProducerThumbnail(ctx context.Context, in *media.GenerateFileInput) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateProducerThumbnail", ctx, in)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GenerateProducerThumbnail indicates an expected call of GenerateProducerThumbnail.
func (mr *MockServiceMockRecorder) GenerateProducerThumbnail(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateProducerThumbnail", reflect.TypeOf((*MockService)(nil).GenerateProducerThumbnail), ctx, in)
}

// GenerateProductMediaImage mocks base method.
func (m *MockService) GenerateProductMediaImage(ctx context.Context, in *media.GenerateFileInput) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateProductMediaImage", ctx, in)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GenerateProductMediaImage indicates an expected call of GenerateProductMediaImage.
func (mr *MockServiceMockRecorder) GenerateProductMediaImage(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateProductMediaImage", reflect.TypeOf((*MockService)(nil).GenerateProductMediaImage), ctx, in)
}

// GenerateProductMediaVideo mocks base method.
func (m *MockService) GenerateProductMediaVideo(ctx context.Context, in *media.GenerateFileInput) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateProductMediaVideo", ctx, in)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GenerateProductMediaVideo indicates an expected call of GenerateProductMediaVideo.
func (mr *MockServiceMockRecorder) GenerateProductMediaVideo(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateProductMediaVideo", reflect.TypeOf((*MockService)(nil).GenerateProductMediaVideo), ctx, in)
}

// GenerateProductTypeIcon mocks base method.
func (m *MockService) GenerateProductTypeIcon(ctx context.Context, in *media.GenerateFileInput) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateProductTypeIcon", ctx, in)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GenerateProductTypeIcon indicates an expected call of GenerateProductTypeIcon.
func (mr *MockServiceMockRecorder) GenerateProductTypeIcon(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateProductTypeIcon", reflect.TypeOf((*MockService)(nil).GenerateProductTypeIcon), ctx, in)
}

// GenerateScheduleImage mocks base method.
func (m *MockService) GenerateScheduleImage(ctx context.Context, in *media.GenerateFileInput) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateScheduleImage", ctx, in)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GenerateScheduleImage indicates an expected call of GenerateScheduleImage.
func (mr *MockServiceMockRecorder) GenerateScheduleImage(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateScheduleImage", reflect.TypeOf((*MockService)(nil).GenerateScheduleImage), ctx, in)
}

// GenerateScheduleOpeningVideo mocks base method.
func (m *MockService) GenerateScheduleOpeningVideo(ctx context.Context, in *media.GenerateFileInput) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateScheduleOpeningVideo", ctx, in)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GenerateScheduleOpeningVideo indicates an expected call of GenerateScheduleOpeningVideo.
func (mr *MockServiceMockRecorder) GenerateScheduleOpeningVideo(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateScheduleOpeningVideo", reflect.TypeOf((*MockService)(nil).GenerateScheduleOpeningVideo), ctx, in)
}

// GenerateScheduleThumbnail mocks base method.
func (m *MockService) GenerateScheduleThumbnail(ctx context.Context, in *media.GenerateFileInput) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateScheduleThumbnail", ctx, in)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GenerateScheduleThumbnail indicates an expected call of GenerateScheduleThumbnail.
func (mr *MockServiceMockRecorder) GenerateScheduleThumbnail(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateScheduleThumbnail", reflect.TypeOf((*MockService)(nil).GenerateScheduleThumbnail), ctx, in)
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

// UploadCoordinatorBonusVideo mocks base method.
func (m *MockService) UploadCoordinatorBonusVideo(ctx context.Context, in *media.UploadFileInput) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UploadCoordinatorBonusVideo", ctx, in)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UploadCoordinatorBonusVideo indicates an expected call of UploadCoordinatorBonusVideo.
func (mr *MockServiceMockRecorder) UploadCoordinatorBonusVideo(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UploadCoordinatorBonusVideo", reflect.TypeOf((*MockService)(nil).UploadCoordinatorBonusVideo), ctx, in)
}

// UploadCoordinatorHeader mocks base method.
func (m *MockService) UploadCoordinatorHeader(ctx context.Context, in *media.UploadFileInput) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UploadCoordinatorHeader", ctx, in)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UploadCoordinatorHeader indicates an expected call of UploadCoordinatorHeader.
func (mr *MockServiceMockRecorder) UploadCoordinatorHeader(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UploadCoordinatorHeader", reflect.TypeOf((*MockService)(nil).UploadCoordinatorHeader), ctx, in)
}

// UploadCoordinatorPromotionVideo mocks base method.
func (m *MockService) UploadCoordinatorPromotionVideo(ctx context.Context, in *media.UploadFileInput) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UploadCoordinatorPromotionVideo", ctx, in)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UploadCoordinatorPromotionVideo indicates an expected call of UploadCoordinatorPromotionVideo.
func (mr *MockServiceMockRecorder) UploadCoordinatorPromotionVideo(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UploadCoordinatorPromotionVideo", reflect.TypeOf((*MockService)(nil).UploadCoordinatorPromotionVideo), ctx, in)
}

// UploadCoordinatorThumbnail mocks base method.
func (m *MockService) UploadCoordinatorThumbnail(ctx context.Context, in *media.UploadFileInput) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UploadCoordinatorThumbnail", ctx, in)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UploadCoordinatorThumbnail indicates an expected call of UploadCoordinatorThumbnail.
func (mr *MockServiceMockRecorder) UploadCoordinatorThumbnail(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UploadCoordinatorThumbnail", reflect.TypeOf((*MockService)(nil).UploadCoordinatorThumbnail), ctx, in)
}

// UploadProducerBonusVideo mocks base method.
func (m *MockService) UploadProducerBonusVideo(ctx context.Context, in *media.UploadFileInput) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UploadProducerBonusVideo", ctx, in)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UploadProducerBonusVideo indicates an expected call of UploadProducerBonusVideo.
func (mr *MockServiceMockRecorder) UploadProducerBonusVideo(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UploadProducerBonusVideo", reflect.TypeOf((*MockService)(nil).UploadProducerBonusVideo), ctx, in)
}

// UploadProducerHeader mocks base method.
func (m *MockService) UploadProducerHeader(ctx context.Context, in *media.UploadFileInput) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UploadProducerHeader", ctx, in)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UploadProducerHeader indicates an expected call of UploadProducerHeader.
func (mr *MockServiceMockRecorder) UploadProducerHeader(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UploadProducerHeader", reflect.TypeOf((*MockService)(nil).UploadProducerHeader), ctx, in)
}

// UploadProducerPromotionVideo mocks base method.
func (m *MockService) UploadProducerPromotionVideo(ctx context.Context, in *media.UploadFileInput) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UploadProducerPromotionVideo", ctx, in)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UploadProducerPromotionVideo indicates an expected call of UploadProducerPromotionVideo.
func (mr *MockServiceMockRecorder) UploadProducerPromotionVideo(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UploadProducerPromotionVideo", reflect.TypeOf((*MockService)(nil).UploadProducerPromotionVideo), ctx, in)
}

// UploadProducerThumbnail mocks base method.
func (m *MockService) UploadProducerThumbnail(ctx context.Context, in *media.UploadFileInput) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UploadProducerThumbnail", ctx, in)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UploadProducerThumbnail indicates an expected call of UploadProducerThumbnail.
func (mr *MockServiceMockRecorder) UploadProducerThumbnail(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UploadProducerThumbnail", reflect.TypeOf((*MockService)(nil).UploadProducerThumbnail), ctx, in)
}

// UploadProductMedia mocks base method.
func (m *MockService) UploadProductMedia(ctx context.Context, in *media.UploadFileInput) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UploadProductMedia", ctx, in)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UploadProductMedia indicates an expected call of UploadProductMedia.
func (mr *MockServiceMockRecorder) UploadProductMedia(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UploadProductMedia", reflect.TypeOf((*MockService)(nil).UploadProductMedia), ctx, in)
}

// UploadProductTypeIcon mocks base method.
func (m *MockService) UploadProductTypeIcon(ctx context.Context, in *media.UploadFileInput) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UploadProductTypeIcon", ctx, in)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UploadProductTypeIcon indicates an expected call of UploadProductTypeIcon.
func (mr *MockServiceMockRecorder) UploadProductTypeIcon(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UploadProductTypeIcon", reflect.TypeOf((*MockService)(nil).UploadProductTypeIcon), ctx, in)
}

// UploadScheduleImage mocks base method.
func (m *MockService) UploadScheduleImage(ctx context.Context, in *media.UploadFileInput) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UploadScheduleImage", ctx, in)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UploadScheduleImage indicates an expected call of UploadScheduleImage.
func (mr *MockServiceMockRecorder) UploadScheduleImage(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UploadScheduleImage", reflect.TypeOf((*MockService)(nil).UploadScheduleImage), ctx, in)
}

// UploadScheduleOpeningVideo mocks base method.
func (m *MockService) UploadScheduleOpeningVideo(ctx context.Context, in *media.UploadFileInput) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UploadScheduleOpeningVideo", ctx, in)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UploadScheduleOpeningVideo indicates an expected call of UploadScheduleOpeningVideo.
func (mr *MockServiceMockRecorder) UploadScheduleOpeningVideo(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UploadScheduleOpeningVideo", reflect.TypeOf((*MockService)(nil).UploadScheduleOpeningVideo), ctx, in)
}

// UploadScheduleThumbnail mocks base method.
func (m *MockService) UploadScheduleThumbnail(ctx context.Context, in *media.UploadFileInput) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UploadScheduleThumbnail", ctx, in)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UploadScheduleThumbnail indicates an expected call of UploadScheduleThumbnail.
func (mr *MockServiceMockRecorder) UploadScheduleThumbnail(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UploadScheduleThumbnail", reflect.TypeOf((*MockService)(nil).UploadScheduleThumbnail), ctx, in)
}
