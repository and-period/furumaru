package auth

import (
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"time"

	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/uuid"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidAccessToken  = errors.New("auth: invalid access token")
	ErrInvalidRefreshToken = errors.New("auth: invalid refresh token")
	ErrNoPemBlock          = errors.New("auth: no pem block found")
	ErrNotRSAPrivateKey    = errors.New("auth: not rsa private key")
	ErrRefreshTokenExpired = errors.New("auth: refresh token expired")
)

type options struct {
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
	signingMethod   *jwt.SigningMethodRSA
	now             func() time.Time
	generateID      func() string
}

type Option func(opts *options)

func WithAccessTokenTTL(ttl time.Duration) Option {
	return func(opts *options) {
		opts.accessTokenTTL = ttl
	}
}

func WithRefreshTokenTTL(ttl time.Duration) Option {
	return func(opts *options) {
		opts.refreshTokenTTL = ttl
	}
}

func WithSigningMethod(method *jwt.SigningMethodRSA) Option {
	return func(opts *options) {
		opts.signingMethod = method
	}
}

func WithNow(now func() time.Time) Option {
	return func(opts *options) {
		opts.now = now
	}
}

func WithGenerateID(f func() string) Option {
	return func(opts *options) {
		opts.generateID = f
	}
}

func buildOptions(opts ...Option) *options {
	dopts := &options{
		accessTokenTTL:  30 * time.Minute,    // 30分
		refreshTokenTTL: 14 * 24 * time.Hour, // 2週間
		signingMethod:   jwt.SigningMethodRS256,
		now:             jst.Now,
		generateID:      uuid.New,
	}
	for i := range opts {
		opts[i](dopts)
	}
	return dopts
}

type RefreshToken struct {
	RefreshToken string    `dynamodbav:"-"`                   // リフレッシュトークン
	HashedToken  string    `dynamodbav:"hashed_token"`        // リフレッシュトークン（ハッシュ値）
	UserID       string    `dynamodbav:"user_id"`             // 施設利用者ID
	FacilityID   string    `dynamodbav:"facility_id"`         // 生産者ID
	ExpiredAt    time.Time `dynamodbav:"expired_at,unixtime"` // 有効期限
	CreatedAt    time.Time `dynamodbav:"created_at"`          // 登録日時
	UpdatedAt    time.Time `dynamodbav:"updated_at"`          // 更新日時
}

type RefreshTokenParams struct {
	UserID     string
	FacilityID string
	Now        time.Time
	TTL        time.Duration
}

func NewRefreshToken(params *RefreshTokenParams) (*RefreshToken, error) {
	refreshToken := newRefreshToken(uuid.New())
	hashedToken, err := hashRefreshToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("auth: failed to hash refresh token: %w", err)
	}
	res := &RefreshToken{
		RefreshToken: refreshToken,
		HashedToken:  hashedToken,
		UserID:       params.UserID,
		FacilityID:   params.FacilityID,
		ExpiredAt:    params.Now.Add(params.TTL),
		CreatedAt:    params.Now,
		UpdatedAt:    params.Now,
	}
	return res, nil
}

func (t *RefreshToken) TableName() string {
	return "auth-tokens"
}

func (t *RefreshToken) PrimaryKey() map[string]interface{} {
	return map[string]interface{}{
		"hashed_token": t.HashedToken,
	}
}

func (t *RefreshToken) Verify(refreshToken string) error {
	return compareRefreshToken(t.HashedToken, refreshToken)
}

func newRefreshToken(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return base64.RawURLEncoding.EncodeToString(sum[:])
}

func hashRefreshToken(raw string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(raw), bcrypt.DefaultCost)
	return string(hash), err
}

func compareRefreshToken(hashed string, raw string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(raw))
}

func parseRSAPrivateKeyFromPEM(pemBytes []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, ErrNoPemBlock
	}

	if key, err := x509.ParsePKCS1PrivateKey(block.Bytes); err == nil {
		return key, nil
	}
	if keyAny, err := x509.ParsePKCS8PrivateKey(block.Bytes); err == nil {
		if rsaKey, ok := keyAny.(*rsa.PrivateKey); ok {
			return rsaKey, nil
		}
	}
	return nil, ErrNotRSAPrivateKey
}
