package auth

import (
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/uuid"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidAccessToken  = errors.New("auth: invalid access token")
	ErrInvalidRefreshToken = errors.New("auth: invalid refresh token")
	ErrRefreshTokenExpired = errors.New("auth: refresh token expired")
	ErrEmailNotFound       = errors.New("auth: email not found in claims")
	ErrEmailUnverified     = errors.New("auth: email not verified in claims")
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
4%3
// $ openssl genrsa 2048
func parseRSAPrivateKey(pemStr []byte) (*rsa.PrivateKey, error) {
	convertedPem, err := convertPemFromString(pemStr)
	if err != nil {
		return nil, fmt.Errorf("auth: failed to convert pem string: %w", err)
	}

	block, _ := pem.Decode(convertedPem)
	if block == nil {
		return nil, errors.New("auth: failed to parse PEM block containing the private key")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		// PKCS#8形式も試してみる
		key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("auth: failed to parse private key: %w", err)
		}
		rsaPrivateKey, ok := key.(*rsa.PrivateKey)
		if !ok {
			return nil, errors.New("auth: key type is not RSA")
		}
		return rsaPrivateKey, nil
	}

	return privateKey, nil
}

func convertPemFromString(pemStr []byte) ([]byte, error) {
	pemString := string(pemStr)

	// \nを実際の改行コードに変換
	pemString = strings.ReplaceAll(pemString, "\\n", "\n")

	// 既に適切な改行コードが使われている場合はそのまま返す
	if strings.Contains(pemString, "\n") {
		return []byte(pemString), nil
	}

	return nil, errors.New("auth: invalid PEM format")
}
