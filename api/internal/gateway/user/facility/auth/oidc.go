package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/coreos/go-oidc/v3/oidc"
)

type OIDCVerifier interface {
	VerifyIDToken(ctx context.Context, idToken, nonce string) (*oidc.IDToken, error) // IDトークンの検証
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
