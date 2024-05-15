package youtube

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

// CreateLiveBroadcastParams - ライブ配信作成パラメータ
// @see - https://developers.google.com/youtube/v3/live/docs/liveBroadcasts/insert?hl=ja
type CreateLiveBroadcastParams struct {
	Title       string    // ライブ配信タイトル
	Description string    // ライブ配信説明
	StartAt     time.Time // 配信開始予定日時
	EndAt       time.Time // 配信終了予定日時
	Public      bool      // 配信を公開設定（true:公開,false:限定公開）
}

type service struct {
	service       *youtube.Service
	logger        *zap.Logger
	livePublished bool
}

func (c *client) NewService(ctx context.Context, token *oauth2.Token) (Service, error) {
	httpClient := c.NewAuth().Client(ctx, token)
	srv, err := youtube.NewService(ctx, option.WithHTTPClient(httpClient))
	if err != nil {
		return nil, err
	}
	return &service{
		service: srv,
		logger:  c.logger,
	}, nil
}

func (s *service) CreateLiveBroadcast(ctx context.Context, params *CreateLiveBroadcastParams) (*youtube.LiveBroadcast, error) {
	privacyStatus := "unlisted"
	if s.livePublished {
		privacyStatus = "public"
	}
	in := &youtube.LiveBroadcast{
		Snippet: &youtube.LiveBroadcastSnippet{
			Title:              params.Title,
			Description:        params.Description,
			ScheduledStartTime: params.StartAt.Format(time.RFC3339),
			ScheduledEndTime:   params.EndAt.Format(time.RFC3339),
		},
		Status: &youtube.LiveBroadcastStatus{
			PrivacyStatus:           privacyStatus,
			SelfDeclaredMadeForKids: false,
		},
		ContentDetails: &youtube.LiveBroadcastContentDetails{
			EnableAutoStart: true,
			EnableAutoStop:  true,
			EnableDvr:       true,
			RecordFromStart: true,
		},
	}
	part := []string{"id", "snippet", "contentDetails", "status"}
	res, err := s.service.LiveBroadcasts.Insert(part, in).Context(ctx).Do()
	if err != nil {
		return nil, s.internalError(err)
	}
	return res, nil
}

func (s *service) GetLiveStream(ctx context.Context, streamID string) (*youtube.LiveStream, error) {
	part := []string{"id", "snippet", "cdn"}
	out, err := s.service.LiveStreams.List(part).Id(streamID).Context(ctx).Do()
	if err != nil {
		return nil, s.internalError(err)
	}
	if len(out.Items) == 0 {
		return nil, ErrNotFound
	}
	return out.Items[0], nil
}

func (s *service) internalError(err error) error {
	if err == nil {
		return nil
	}
	s.logger.Debug("Failed to youtube api", zap.Error(err))

	switch {
	case errors.Is(err, context.Canceled):
		return fmt.Errorf("%w: %s", ErrCanceled, err.Error())
	case errors.Is(err, context.DeadlineExceeded):
		return fmt.Errorf("%w: %s", ErrTimeout, err.Error())
	}

	var e *googleapi.Error
	if !errors.As(err, &e) {
		return fmt.Errorf("%w: %s", ErrUnknown, err.Error())
	}

	switch e.Code {
	case http.StatusBadRequest:
		return fmt.Errorf("%w: %s", ErrBadRequest, e.Message)
	case http.StatusUnauthorized:
		return fmt.Errorf("%w: %s", ErrUnauthorized, e.Message)
	case http.StatusForbidden:
		return fmt.Errorf("%w: %s", ErrForbidden, e.Message)
	case http.StatusTooManyRequests:
		return fmt.Errorf("%w: %s", ErrTooManyRequests, e.Message)
	default:
		return fmt.Errorf("%w: %s", ErrUnknown, e.Message)
	}
}
