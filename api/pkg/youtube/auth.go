package youtube

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/youtube/v3"
)

type auth struct {
	*oauth2.Config
	logger *zap.Logger
}

func (c *client) NewAuth() Auth {
	config := &oauth2.Config{
		ClientID:     c.clientID,
		ClientSecret: c.clientSecret,
		Endpoint:     google.Endpoint,
		Scopes:       []string{youtube.YoutubeScope},
		RedirectURL:  c.authCallbackURL,
	}
	return &auth{
		Config: config,
		logger: c.logger,
	}
}

func (a *auth) AuthCodeURL(state string) string {
	opts := []oauth2.AuthCodeOption{
		oauth2.AccessTypeOffline,
		oauth2.ApprovalForce,
	}
	return a.Config.AuthCodeURL(state, opts...)
}

func (a *auth) Exchange(ctx context.Context, code string, opts ...oauth2.AuthCodeOption) (*oauth2.Token, error) {
	token, err := a.Config.Exchange(ctx, code, opts...)
	if err == nil {
		return token, nil
	}
	a.logger.Error("failed to exchange token", zap.Error(err))

	var e *oauth2.RetrieveError
	if !errors.As(err, &e) {
		return nil, fmt.Errorf("%w: %s", ErrUnknown, err.Error())
	}
	switch e.Response.StatusCode {
	case http.StatusBadRequest:
		return nil, fmt.Errorf("%w: %s", ErrBadRequest, e.Error())
	case http.StatusUnauthorized:
		return nil, fmt.Errorf("%w: %s", ErrUnauthorized, e.Error())
	case http.StatusForbidden:
		return nil, fmt.Errorf("%w: %s", ErrForbidden, e.Error())
	default:
		return nil, fmt.Errorf("%w: %s", ErrUnknown, e.Error())
	}
}
