package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/coreos/go-oidc/v3/oidc"
)

type OIDCVerifier[Token any] interface {
	VerifyIDToken(ctx context.Context, idToken, nonce string) (*oidc.IDToken, error) // IDトークンの検証
	GetClaims(token *oidc.IDToken) (*Token, error)                                   // IDトークンからクレームの取得
}

type liffVerifier struct {
	verifier *oidc.IDTokenVerifier
}

// LIFFのIDトークンのクレーム
// @see - https://developers.line.biz/ja/docs/basics/user-profile/#profile-information-types
type liffClaims struct {
	Sub         string `json:"sub,omitempty"`          // ユーザーID
	Name        string `json:"name,omitempty"`         // 表示名
	Picture     string `json:"picture,omitempty"`      // プロフィール画像
	Email       string `json:"email,omitempty"`        // メールアドレス
	GivenName   string `json:"given_name,omitempty"`   // 氏名
	FamilyName  string `json:"family_name,omitempty"`  // 氏名
	Gender      string `json:"gender,omitempty"`       // 性別
	Birthdate   string `json:"birthdate,omitempty"`    // 誕生日
	Address     string `json:"address,omitempty"`      // 住所
	PhoneNumber string `json:"phone_number,omitempty"` // 電話番号
}

// LIFFのトークン詳細情報
type LIFFClaims struct {
	Iss     string   `json:"iss"`     // 発行者
	Sub     string   `json:"sub"`     // ユーザーID
	Aud     string   `json:"aud"`     // クライアントID
	Exp     int64    `json:"exp"`     // 有効期限
	Iat     int64    `json:"iat"`     // 発行日時
	Nonce   string   `json:"nonce"`   // リプレイ攻撃防止用の値
	Amr     []string `json:"amr"`     // 認証方法
	Name    string   `json:"name"`    // 表示名
	Picture string   `json:"picture"` // プロフィール画像URL
	Email   string   `json:"email"`   // メールアドレス
}

func NewLIFFVerifier(ctx context.Context) (OIDCVerifier[LIFFClaims], error) {
	const issuer = "https://access.line.me"
	provider, err := oidc.NewProvider(ctx, issuer)
	if err != nil {
		return nil, fmt.Errorf("verifier: failed to line verifier: %w", err)
	}
	v := provider.Verifier(&oidc.Config{
		SkipClientIDCheck: true,
		SkipExpiryCheck:   false,
		SkipIssuerCheck:   false,
	})
	client := &liffVerifier{
		verifier: v,
	}
	return client, nil
}

func (v *liffVerifier) VerifyIDToken(ctx context.Context, idToken, nonce string) (*oidc.IDToken, error) {
	token, err := v.verifier.Verify(ctx, idToken)
	if err != nil {
		return nil, fmt.Errorf("verifier: failed to verify liff id token: %w", err)
	}
	slog.DebugContext(ctx, "Verified ID token", slog.Any("token", token))
	if nonce != "" && nonce != token.Nonce {
		return nil, errors.New("verifier: invalid nonce")
	}
	return token, nil
}

func (v *liffVerifier) GetClaims(token *oidc.IDToken) (*LIFFClaims, error) {
	claims := &LIFFClaims{}
	if err := token.Claims(claims); err != nil {
		return nil, fmt.Errorf("verifier: failed to decode claims: %w", err)
	}
	slog.Debug("Extracted claims", slog.Any("claims", claims))
	return claims, nil
}

func (v *liffVerifier) GetEmail(token *oidc.IDToken) (string, error) {
	claims, err := v.extractClaim(token)
	if err != nil {
		return "", err
	}
	slog.Debug("Extracted claims", slog.Any("claims", claims))
	if claims["email"] == "" {
		return "", ErrEmailNotFound
	}
	return claims["email"], nil
}

func (v *liffVerifier) extractClaim(token *oidc.IDToken) (map[string]string, error) {
	// claims := &liffClaims{}
	claims := make(map[string]string)
	if err := token.Claims(&claims); err != nil {
		return nil, fmt.Errorf("verifier: failed to decode claims: %w", err)
	}
	return claims, nil
}
