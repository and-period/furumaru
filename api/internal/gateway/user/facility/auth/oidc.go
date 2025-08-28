package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/coreos/go-oidc/v3/oidc"
)

type OIDCVerifier interface {
	VerifyIDToken(ctx context.Context, idToken, nonce string) (*oidc.IDToken, error) // IDトークンの検証
	GetEmail(token *oidc.IDToken) (string, error)                                    // IDトークンからメールアドレスの取得
}

type lineVerifier struct {
	verifier *oidc.IDTokenVerifier
}

func NewLineVerifier(ctx context.Context, channelID string) (OIDCVerifier, error) {
	const issuer = "https://access.line.me"
	provider, err := oidc.NewProvider(ctx, issuer)
	if err != nil {
		return nil, fmt.Errorf("verifier: failed to line verifier: %w", err)
	}
	v := provider.Verifier(&oidc.Config{
		ClientID:        channelID,
		SkipExpiryCheck: false,
	})
	client := &lineVerifier{
		verifier: v,
	}
	return client, nil
}

func (v *lineVerifier) VerifyIDToken(ctx context.Context, idToken, nonce string) (*oidc.IDToken, error) {
	token, err := v.verifier.Verify(ctx, idToken)
	if err != nil {
		return nil, fmt.Errorf("verifier: failed to verify line id token: %w", err)
	}
	if nonce != "" && nonce != token.Nonce {
		return nil, errors.New("verifier: invalid nonce")
	}
	return token, nil
}

func (v *lineVerifier) GetEmail(token *oidc.IDToken) (string, error) {
	claims, err := v.extractClaim(token)
	if err != nil {
		return "", err
	}
	if claims.Email == "" {
		return "", ErrEmailNotFound
	}
	if !claims.EmailVerified {
		return "", ErrEmailNotFoundVerified
	}
	return claims.Email, nil
}

type lineClaims struct {
	Sub           string `json:"sub,omitempty"`
	Name          string `json:"name,omitempty"`
	Picture       string `json:"picture,omitempty"`
	Email         string `json:"email,omitempty"`
	EmailVerified bool   `json:"email_verified,omitempty"`
}

func (v *lineVerifier) extractClaim(token *oidc.IDToken) (*lineClaims, error) {
	claims := &lineClaims{}
	if err := token.Claims(&claims); err != nil {
		return nil, fmt.Errorf("verifier: failed to decode claims: %w", err)
	}
	return claims, nil
}
