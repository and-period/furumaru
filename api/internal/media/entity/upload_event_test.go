package entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUploadStatus(t *testing.T) {
	t.Parallel()
	now := time.Now()
	params := &UploadEventParams{
		Key:       "dir/key.png",
		FileGroup: "dir",
		FileType:  "image/png",
		UploadURL: "http://example.com/dir/key.png",
		Now:       now,
		TTL:       time.Hour,
	}
	event := NewUploadEvent(params)
	t.Run("new", func(t *testing.T) {
		t.Parallel()
		expect := &UploadEvent{
			Key:          "dir/key.png",
			Status:       UploadStatusWaiting,
			FileGroup:    "dir",
			FileType:     "image/png",
			UploadURL:    "http://example.com/dir/key.png",
			ReferenceURL: "",
			ExpiredAt:    now.Add(time.Hour),
			CreatedAt:    now,
			UpdatedAt:    now,
		}
		assert.Equal(t, expect, event)
	})
	t.Run("table name", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, "upload-events", event.TableName())
	})
	t.Run("primary key", func(t *testing.T) {
		t.Parallel()
		expect := map[string]interface{}{
			"key": "dir/key.png",
		}
		assert.Equal(t, expect, event.PrimaryKey())
	})
}

func TestUploadEvent_SetResult(t *testing.T) {
	t.Parallel()
	now := time.Now()
	type args struct {
		success      bool
		referenceURL string
		now          time.Time
	}
	tests := []struct {
		name   string
		event  *UploadEvent
		args   args
		expect *UploadEvent
	}{
		{
			name: "succeeded",
			event: &UploadEvent{
				Key:          "dir/key.png",
				Status:       UploadStatusWaiting,
				FileGroup:    "dir",
				FileType:     "image/png",
				UploadURL:    "http://example.com/dir/key.png",
				ReferenceURL: "",
				ExpiredAt:    now.Add(time.Hour),
				CreatedAt:    now.Add(-time.Hour),
				UpdatedAt:    now.Add(-time.Hour),
			},
			args: args{
				success:      true,
				referenceURL: "http://example.com/dir/key.png",
				now:          now,
			},
			expect: &UploadEvent{
				Key:          "dir/key.png",
				Status:       UploadStatusSucceeded,
				FileGroup:    "dir",
				FileType:     "image/png",
				UploadURL:    "http://example.com/dir/key.png",
				ReferenceURL: "http://example.com/dir/key.png",
				ExpiredAt:    now.Add(time.Hour),
				CreatedAt:    now.Add(-time.Hour),
				UpdatedAt:    now,
			},
		},
		{
			name: "failed",
			event: &UploadEvent{
				Key:          "dir/key.png",
				Status:       UploadStatusWaiting,
				FileGroup:    "dir",
				FileType:     "image/png",
				UploadURL:    "http://example.com/dir/key.png",
				ReferenceURL: "",
				ExpiredAt:    now.Add(time.Hour),
				CreatedAt:    now.Add(-time.Hour),
				UpdatedAt:    now.Add(-time.Hour),
			},
			args: args{
				success:      false,
				referenceURL: "http://example.com/dir/key.png",
				now:          now,
			},
			expect: &UploadEvent{
				Key:          "dir/key.png",
				Status:       UploadStatusFailed,
				FileGroup:    "dir",
				FileType:     "image/png",
				UploadURL:    "http://example.com/dir/key.png",
				ReferenceURL: "",
				ExpiredAt:    now.Add(time.Hour),
				CreatedAt:    now.Add(-time.Hour),
				UpdatedAt:    now,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.event.SetResult(tt.args.success, tt.args.referenceURL, tt.args.now)
			assert.Equal(t, tt.expect, tt.event)
		})
	}
}

func TestUploadEvent_Regulation(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		event     *UploadEvent
		expect    *Regulation
		expectErr error
	}{
		// ライブ配信関連
		{
			name:      "broadcast archive mp4",
			event:     &UploadEvent{FileGroup: BroadcastArchiveMP4Path},
			expect:    BroadcastArchiveRegulation,
			expectErr: nil,
		},
		{
			name:      "broadcast live mp4",
			event:     &UploadEvent{FileGroup: BroadcastLiveMP4Path},
			expect:    BroadcastLiveMP4Regulation,
			expectErr: nil,
		},
		// コーディネータ関連
		{
			name:      "coordinator thumbnail",
			event:     &UploadEvent{FileGroup: CoordinatorThumbnailPath},
			expect:    CoordinatorThumbnailRegulation,
			expectErr: nil,
		},
		{
			name:      "coordinator header",
			event:     &UploadEvent{FileGroup: CoordinatorHeaderPath},
			expect:    CoordinatorHeaderRegulation,
			expectErr: nil,
		},
		{
			name:      "coordinator promotion video",
			event:     &UploadEvent{FileGroup: CoordinatorPromotionVideoPath},
			expect:    CoordinatorPromotionVideoRegulation,
			expectErr: nil,
		},
		{
			name:      "coordinator bonus video",
			event:     &UploadEvent{FileGroup: CoordinatorBonusVideoPath},
			expect:    CoordinatorBonusVideoRegulation,
			expectErr: nil,
		},
		// 体験関連
		{
			name:      "experience media image",
			event:     &UploadEvent{FileGroup: ExperienceMediaImagePath},
			expect:    ExperienceMediaImageRegulation,
			expectErr: nil,
		},
		{
			name:      "experience media video",
			event:     &UploadEvent{FileGroup: ExperienceMediaVideoPath},
			expect:    ExperienceMediaVideoRegulation,
			expectErr: nil,
		},
		{
			name:      "experience promotion video",
			event:     &UploadEvent{FileGroup: ExperiencePromotionVideoPath},
			expect:    ExperiencePromotionVideoRegulation,
			expectErr: nil,
		},
		// 生産者関連
		{
			name:      "producer thumbnail",
			event:     &UploadEvent{FileGroup: ProducerThumbnailPath},
			expect:    ProducerThumbnailRegulation,
			expectErr: nil,
		},
		{
			name:      "producer header",
			event:     &UploadEvent{FileGroup: ProducerHeaderPath},
			expect:    ProducerHeaderRegulation,
			expectErr: nil,
		},
		{
			name:      "producer promotion video",
			event:     &UploadEvent{FileGroup: ProducerPromotionVideoPath},
			expect:    ProducerPromotionVideoRegulation,
			expectErr: nil,
		},
		{
			name:      "producer bonus video",
			event:     &UploadEvent{FileGroup: ProducerBonusVideoPath},
			expect:    ProducerBonusVideoRegulation,
			expectErr: nil,
		},
		// 購入者関連
		{
			name:      "user thumbnail",
			event:     &UploadEvent{FileGroup: UserThumbnailPath},
			expect:    UserThumbnailRegulation,
			expectErr: nil,
		},
		// 商品関連
		{
			name:      "product media image",
			event:     &UploadEvent{FileGroup: ProductMediaImagePath},
			expect:    ProductMediaImageRegulation,
			expectErr: nil,
		},
		{
			name:      "product media video",
			event:     &UploadEvent{FileGroup: ProductMediaVideoPath},
			expect:    ProductMediaVideoRegulation,
			expectErr: nil,
		},
		// 品目関連
		{
			name:      "product type icon",
			event:     &UploadEvent{FileGroup: ProductTypeIconPath},
			expect:    ProductTypeIconRegulation,
			expectErr: nil,
		},
		// 開催スケジュール関連
		{
			name:      "schedule thumbnail",
			event:     &UploadEvent{FileGroup: ScheduleThumbnailPath},
			expect:    ScheduleThumbnailRegulation,
			expectErr: nil,
		},
		{
			name:      "schedule image",
			event:     &UploadEvent{FileGroup: ScheduleImagePath},
			expect:    ScheduleImageRegulation,
			expectErr: nil,
		},
		{
			name:      "schedule opening video",
			event:     &UploadEvent{FileGroup: ScheduleOpeningVideoPath},
			expect:    ScheduleOpeningVideoRegulation,
			expectErr: nil,
		},
		{
			name:      "not found",
			event:     &UploadEvent{},
			expect:    nil,
			expectErr: ErrNotFoundReguration,
		},
		// オンデマンド配信関連
		{
			name:      "video thumbnail",
			event:     &UploadEvent{FileGroup: VideoThumbnailPath},
			expect:    VideoThumbnailRegulation,
			expectErr: nil,
		},
		{
			name:      "video mp4",
			event:     &UploadEvent{FileGroup: VideoMP4Path},
			expect:    VideoMP4Regulation,
			expectErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual, err := tt.event.Reguration()
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, actual)
		})
	}
}
