//go:generate mockgen -source=$GOFILE -package mock_$GOPACKAGE -destination=./../../mock/pkg/$GOPACKAGE/$GOFILE
package youtube

import (
	"context"
	"fmt"
	"net/http"

	"go.uber.org/zap"
	"golang.org/x/oauth2"
	gauth "google.golang.org/api/oauth2/v2"
	"google.golang.org/api/youtube/v3"
)

var (
	ErrBadRequest      = fmt.Errorf("youtube: bad request")
	ErrUnauthorized    = fmt.Errorf("youtube: unauthorized")
	ErrForbidden       = fmt.Errorf("youtube: forbidden")
	ErrNotFound        = fmt.Errorf("youtube: not found")
	ErrTooManyRequests = fmt.Errorf("youtube: too many requests")
	ErrCanceled        = fmt.Errorf("youtube: canceled")
	ErrTimeout         = fmt.Errorf("youtube: timeout")
	ErrUnknown         = fmt.Errorf("youtube: unknown error")
)

type YouTube interface {
	NewAuth() Auth
	NewService(ctx context.Context, token *oauth2.Token) (Service, error)
}

type Auth interface {
	Client(ctx context.Context, t *oauth2.Token) *http.Client
	GetAuthCodeURL(state string) string
	GetToken(ctx context.Context, code string, opts ...oauth2.AuthCodeOption) (*oauth2.Token, error)
	GetTokenInfo(ctx context.Context, token *oauth2.Token) (*gauth.Tokeninfo, error)
}

type Service interface {
	CreateLiveBroadcast(ctx context.Context, params *CreateLiveBroadcastParams) (*youtube.LiveBroadcast, error) // ライブ配信作成
	GetLiveStream(ctx context.Context, streamID string) (*youtube.LiveStream, error)                            // ライブ配信先設定取得
}

type Params struct {
	ClientID        string
	ClientSecret    string
	AuthCallbackURL string
}

type client struct {
	logger          *zap.Logger
	clientID        string
	clientSecret    string
	authCallbackURL string
}

type options struct {
	logger *zap.Logger
}

type Option func(*options)

func WithLogger(logger *zap.Logger) Option {
	return func(opts *options) {
		opts.logger = logger
	}
}

func NewClient(params *Params, opts ...Option) YouTube {
	dopts := &options{
		logger: zap.NewNop(),
	}
	for i := range opts {
		opts[i](dopts)
	}
	return &client{
		logger:          dopts.logger,
		clientID:        params.ClientID,
		clientSecret:    params.ClientSecret,
		authCallbackURL: params.AuthCallbackURL,
	}
}
