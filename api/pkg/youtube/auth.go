package youtube

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/and-period/furumaru/api/pkg/log"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/googleapi"
	gauth "google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
	youtube "google.golang.org/api/youtube/v3"
)

type auth struct {
	*oauth2.Config
}

func (c *client) NewAuth() Auth {
	config := &oauth2.Config{
		ClientID:     c.clientID,
		ClientSecret: c.clientSecret,
		Endpoint:     google.Endpoint,
		Scopes:       []string{"openid", "email", "profile", youtube.YoutubeScope},
		RedirectURL:  c.authCallbackURL,
	}
	return &auth{
		Config: config,
	}
}

func (a *auth) GetAuthCodeURL(state string) string {
	opts := []oauth2.AuthCodeOption{
		oauth2.AccessTypeOffline,
		oauth2.ApprovalForce,
	}
	return a.AuthCodeURL(state, opts...)
}

func (a *auth) GetToken(ctx context.Context, code string, opts ...oauth2.AuthCodeOption) (*oauth2.Token, error) {
	token, err := a.Exchange(ctx, code, opts...)
	if err != nil {
		return nil, a.internalError(err)
	}
	if !token.Valid() {
		return nil, fmt.Errorf("%w: token is invalid", ErrUnauthorized)
	}
	return token, nil
}

func (a *auth) GetTokenInfo(ctx context.Context, token *oauth2.Token) (*gauth.Tokeninfo, error) {
	client := a.Client(ctx, token)
	service, err := gauth.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, a.internalError(err)
	}
	user, err := service.Tokeninfo().AccessToken(token.AccessToken).Context(ctx).Do()
	if err != nil {
		return nil, a.internalError(err)
	}
	return user, nil
}

func (a *auth) internalError(err error) error {
	if err == nil {
		return nil
	}
	slog.Error("Failed to oauth2 api", log.Error(err))

	code := ErrUnknown

	var e *googleapi.Error
	if errors.As(err, &e) {
		switch e.Code {
		case http.StatusBadRequest:
			code = ErrBadRequest
		case http.StatusUnauthorized:
			code = ErrUnauthorized
		case http.StatusForbidden:
			code = ErrForbidden
		case http.StatusTooManyRequests:
			code = ErrTooManyRequests
		}
		return fmt.Errorf("%w: %s", code, e.Error())
	}

	var ae *oauth2.RetrieveError
	if errors.As(err, &ae) {
		switch ae.Response.StatusCode {
		case http.StatusBadRequest:
			code = ErrBadRequest
		case http.StatusUnauthorized:
			code = ErrUnauthorized
		case http.StatusForbidden:
			code = ErrForbidden
		}
		return fmt.Errorf("%w: %s", code, ae.Error())
	}

	return fmt.Errorf("%w: %s", code, err.Error())
}
