package service

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/media"
	"github.com/and-period/furumaru/api/internal/media/entity"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResizeCoordinatorThumbnail(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		setup  func(ctx context.Context, mocks *mocks)
		input  *media.ResizeFileInput
		expect error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.producer.EXPECT().
					SendMessage(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, buf []byte) (string, error) {
						payload := &entity.ResizerPayload{}
						err := json.Unmarshal(buf, payload)
						require.NoError(t, err)
						expect := &entity.ResizerPayload{
							TargetID: "target-id",
							FileType: entity.FileTypeCoordinatorThumbnail,
							URLs:     []string{"http://example.com/test.jpg"},
						}
						assert.Equal(t, expect, payload)
						return "message-id", nil
					})
			},
			input: &media.ResizeFileInput{
				TargetID: "target-id",
				URLs:     []string{"http://example.com/test.jpg"},
			},
			expect: nil,
		},
		{
			name:   "invalid argument",
			setup:  func(ctx context.Context, mocks *mocks) {},
			input:  &media.ResizeFileInput{},
			expect: exception.ErrInvalidArgument,
		},
		{
			name: "failed to send message",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.producer.EXPECT().SendMessage(ctx, gomock.Any()).Return("", assert.AnError)
			},
			input: &media.ResizeFileInput{
				TargetID: "target-id",
				URLs:     []string{"http://example.com/test.jpg"},
			},
			expect: exception.ErrInternal,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.ResizeCoordinatorThumbnail(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expect)
		}))
	}
}

func TestResizeCoordinatorHeader(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		setup  func(ctx context.Context, mocks *mocks)
		input  *media.ResizeFileInput
		expect error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.producer.EXPECT().
					SendMessage(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, buf []byte) (string, error) {
						payload := &entity.ResizerPayload{}
						err := json.Unmarshal(buf, payload)
						require.NoError(t, err)
						expect := &entity.ResizerPayload{
							TargetID: "target-id",
							FileType: entity.FileTypeCoordinatorHeader,
							URLs:     []string{"http://example.com/test.jpg"},
						}
						assert.Equal(t, expect, payload)
						return "message-id", nil
					})
			},
			input: &media.ResizeFileInput{
				TargetID: "target-id",
				URLs:     []string{"http://example.com/test.jpg"},
			},
			expect: nil,
		},
		{
			name:   "invalid argument",
			setup:  func(ctx context.Context, mocks *mocks) {},
			input:  &media.ResizeFileInput{},
			expect: exception.ErrInvalidArgument,
		},
		{
			name: "failed to send message",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.producer.EXPECT().SendMessage(ctx, gomock.Any()).Return("", assert.AnError)
			},
			input: &media.ResizeFileInput{
				TargetID: "target-id",
				URLs:     []string{"http://example.com/test.jpg"},
			},
			expect: exception.ErrInternal,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.ResizeCoordinatorHeader(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expect)
		}))
	}
}

func TestResizeProducerThumbnail(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		setup  func(ctx context.Context, mocks *mocks)
		input  *media.ResizeFileInput
		expect error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.producer.EXPECT().
					SendMessage(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, buf []byte) (string, error) {
						payload := &entity.ResizerPayload{}
						err := json.Unmarshal(buf, payload)
						require.NoError(t, err)
						expect := &entity.ResizerPayload{
							TargetID: "target-id",
							FileType: entity.FileTypeProducerThumbnail,
							URLs:     []string{"http://example.com/test.jpg"},
						}
						assert.Equal(t, expect, payload)
						return "message-id", nil
					})
			},
			input: &media.ResizeFileInput{
				TargetID: "target-id",
				URLs:     []string{"http://example.com/test.jpg"},
			},
			expect: nil,
		},
		{
			name:   "invalid argument",
			setup:  func(ctx context.Context, mocks *mocks) {},
			input:  &media.ResizeFileInput{},
			expect: exception.ErrInvalidArgument,
		},
		{
			name: "failed to send message",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.producer.EXPECT().SendMessage(ctx, gomock.Any()).Return("", assert.AnError)
			},
			input: &media.ResizeFileInput{
				TargetID: "target-id",
				URLs:     []string{"http://example.com/test.jpg"},
			},
			expect: exception.ErrInternal,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.ResizeProducerThumbnail(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expect)
		}))
	}
}

func TestResizeProducerHeader(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		setup  func(ctx context.Context, mocks *mocks)
		input  *media.ResizeFileInput
		expect error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.producer.EXPECT().
					SendMessage(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, buf []byte) (string, error) {
						payload := &entity.ResizerPayload{}
						err := json.Unmarshal(buf, payload)
						require.NoError(t, err)
						expect := &entity.ResizerPayload{
							TargetID: "target-id",
							FileType: entity.FileTypeProducerHeader,
							URLs:     []string{"http://example.com/test.jpg"},
						}
						assert.Equal(t, expect, payload)
						return "message-id", nil
					})
			},
			input: &media.ResizeFileInput{
				TargetID: "target-id",
				URLs:     []string{"http://example.com/test.jpg"},
			},
			expect: nil,
		},
		{
			name:   "invalid argument",
			setup:  func(ctx context.Context, mocks *mocks) {},
			input:  &media.ResizeFileInput{},
			expect: exception.ErrInvalidArgument,
		},
		{
			name: "failed to send message",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.producer.EXPECT().SendMessage(ctx, gomock.Any()).Return("", assert.AnError)
			},
			input: &media.ResizeFileInput{
				TargetID: "target-id",
				URLs:     []string{"http://example.com/test.jpg"},
			},
			expect: exception.ErrInternal,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.ResizeProducerHeader(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expect)
		}))
	}
}

func TestResizeUserThumbnail(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		setup  func(ctx context.Context, mocks *mocks)
		input  *media.ResizeFileInput
		expect error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.producer.EXPECT().
					SendMessage(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, buf []byte) (string, error) {
						payload := &entity.ResizerPayload{}
						err := json.Unmarshal(buf, payload)
						require.NoError(t, err)
						expect := &entity.ResizerPayload{
							TargetID: "target-id",
							FileType: entity.FileTypeUserThumbnail,
							URLs:     []string{"http://example.com/test.jpg"},
						}
						assert.Equal(t, expect, payload)
						return "message-id", nil
					})
			},
			input: &media.ResizeFileInput{
				TargetID: "target-id",
				URLs:     []string{"http://example.com/test.jpg"},
			},
			expect: nil,
		},
		{
			name:   "invalid argument",
			setup:  func(ctx context.Context, mocks *mocks) {},
			input:  &media.ResizeFileInput{},
			expect: exception.ErrInvalidArgument,
		},
		{
			name: "failed to send message",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.producer.EXPECT().SendMessage(ctx, gomock.Any()).Return("", assert.AnError)
			},
			input: &media.ResizeFileInput{
				TargetID: "target-id",
				URLs:     []string{"http://example.com/test.jpg"},
			},
			expect: exception.ErrInternal,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.ResizeUserThumbnail(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expect)
		}))
	}
}

func TestResizeProductMedia(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		setup  func(ctx context.Context, mocks *mocks)
		input  *media.ResizeFileInput
		expect error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.producer.EXPECT().
					SendMessage(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, buf []byte) (string, error) {
						payload := &entity.ResizerPayload{}
						err := json.Unmarshal(buf, payload)
						require.NoError(t, err)
						expect := &entity.ResizerPayload{
							TargetID: "target-id",
							FileType: entity.FileTypeProductMedia,
							URLs:     []string{"http://example.com/test.jpg"},
						}
						assert.Equal(t, expect, payload)
						return "message-id", nil
					})
			},
			input: &media.ResizeFileInput{
				TargetID: "target-id",
				URLs:     []string{"http://example.com/test.jpg"},
			},
			expect: nil,
		},
		{
			name:   "invalid argument",
			setup:  func(ctx context.Context, mocks *mocks) {},
			input:  &media.ResizeFileInput{},
			expect: exception.ErrInvalidArgument,
		},
		{
			name: "failed to send message",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.producer.EXPECT().SendMessage(ctx, gomock.Any()).Return("", assert.AnError)
			},
			input: &media.ResizeFileInput{
				TargetID: "target-id",
				URLs:     []string{"http://example.com/test.jpg"},
			},
			expect: exception.ErrInternal,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.ResizeProductMedia(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expect)
		}))
	}
}

func TestResizeProductTypeIcon(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		setup  func(ctx context.Context, mocks *mocks)
		input  *media.ResizeFileInput
		expect error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.producer.EXPECT().
					SendMessage(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, buf []byte) (string, error) {
						payload := &entity.ResizerPayload{}
						err := json.Unmarshal(buf, payload)
						require.NoError(t, err)
						expect := &entity.ResizerPayload{
							TargetID: "target-id",
							FileType: entity.FileTypeProductTypeIcon,
							URLs:     []string{"http://example.com/test.jpg"},
						}
						assert.Equal(t, expect, payload)
						return "message-id", nil
					})
			},
			input: &media.ResizeFileInput{
				TargetID: "target-id",
				URLs:     []string{"http://example.com/test.jpg"},
			},
			expect: nil,
		},
		{
			name:   "invalid argument",
			setup:  func(ctx context.Context, mocks *mocks) {},
			input:  &media.ResizeFileInput{},
			expect: exception.ErrInvalidArgument,
		},
		{
			name: "failed to send message",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.producer.EXPECT().SendMessage(ctx, gomock.Any()).Return("", assert.AnError)
			},
			input: &media.ResizeFileInput{
				TargetID: "target-id",
				URLs:     []string{"http://example.com/test.jpg"},
			},
			expect: exception.ErrInternal,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.ResizeProductTypeIcon(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expect)
		}))
	}
}

func TestResizeScheduleThumbnail(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		setup  func(ctx context.Context, mocks *mocks)
		input  *media.ResizeFileInput
		expect error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.producer.EXPECT().
					SendMessage(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, buf []byte) (string, error) {
						payload := &entity.ResizerPayload{}
						err := json.Unmarshal(buf, payload)
						require.NoError(t, err)
						expect := &entity.ResizerPayload{
							TargetID: "target-id",
							FileType: entity.FileTypeScheduleThumbnail,
							URLs:     []string{"http://example.com/test.jpg"},
						}
						assert.Equal(t, expect, payload)
						return "message-id", nil
					})
			},
			input: &media.ResizeFileInput{
				TargetID: "target-id",
				URLs:     []string{"http://example.com/test.jpg"},
			},
			expect: nil,
		},
		{
			name:   "invalid argument",
			setup:  func(ctx context.Context, mocks *mocks) {},
			input:  &media.ResizeFileInput{},
			expect: exception.ErrInvalidArgument,
		},
		{
			name: "failed to send message",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.producer.EXPECT().SendMessage(ctx, gomock.Any()).Return("", assert.AnError)
			},
			input: &media.ResizeFileInput{
				TargetID: "target-id",
				URLs:     []string{"http://example.com/test.jpg"},
			},
			expect: exception.ErrInternal,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.ResizeScheduleThumbnail(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expect)
		}))
	}
}
