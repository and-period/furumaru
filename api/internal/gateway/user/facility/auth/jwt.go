package auth

import (
	"context"
	"crypto/rsa"
	"fmt"
	"time"

	"github.com/and-period/furumaru/api/pkg/dynamodb"
	"github.com/golang-jwt/jwt/v5"
)

type JWTGenerator interface {
	Generate(ctx context.Context, sub, facilityID string) (*Auth, error)
	GenerateAccessToken(sub, facilityID string) (string, error)
	GenerateRefreshToken(ctx context.Context, sub, facilityID string) (string, error)
	RefreshAccessToken(ctx context.Context, refreshToken string) (*Auth, error)
	DeleteRefreshToken(ctx context.Context, userID string) error
}

type JWTGeneratorParams struct {
	Cache      dynamodb.Client
	Issuer     string
	PrivateKey []byte
}

type jwtGenerator struct {
	issuer          string
	secret          *rsa.PrivateKey
	signingMethod   *jwt.SigningMethodRSA
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
	cache           dynamodb.Client
	now             func() time.Time
	generateID      func() string
}

func NewJWTGenerator(params *JWTGeneratorParams, opts ...Option) (JWTGenerator, error) {
	dopts := buildOptions(opts...)
	secret, err := parseRSAPrivateKey(params.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("auth: failed to parse private key: %w", err)
	}
	client := &jwtGenerator{
		issuer:          params.Issuer,
		secret:          secret,
		signingMethod:   jwt.SigningMethodRS256,
		accessTokenTTL:  dopts.accessTokenTTL,
		refreshTokenTTL: dopts.refreshTokenTTL,
		cache:           params.Cache,
		now:             dopts.now,
		generateID:      dopts.generateID,
	}
	return client, nil
}

func (g *jwtGenerator) Generate(ctx context.Context, sub, facilityID string) (*Auth, error) {
	accessToken, err := g.GenerateAccessToken(sub, facilityID)
	if err != nil {
		return nil, fmt.Errorf("auth: failed to generate access token: %w", err)
	}
	refreshToken, err := g.GenerateRefreshToken(ctx, sub, facilityID)
	if err != nil {
		return nil, fmt.Errorf("auth: failed to generate refresh token: %w", err)
	}
	res := &Auth{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int32(g.accessTokenTTL.Seconds()),
	}
	return res, nil
}

func (g *jwtGenerator) GenerateAccessToken(sub, facilityID string) (string, error) {
	claims := &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    g.issuer,
			Subject:   sub,
			IssuedAt:  jwt.NewNumericDate(g.now()),
			ExpiresAt: jwt.NewNumericDate(g.now().Add(g.accessTokenTTL)),
			ID:        g.generateID(),
		},
		FacilityID: facilityID,
	}
	token, err := jwt.NewWithClaims(g.signingMethod, claims).SignedString(g.secret)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (g *jwtGenerator) GenerateRefreshToken(ctx context.Context, sub, facilityID string) (string, error) {
	params := &RefreshTokenParams{
		UserID:     sub,
		FacilityID: facilityID,
		Now:        g.now(),
		TTL:        g.refreshTokenTTL,
	}
	token, err := NewRefreshToken(params)
	if err != nil {
		return "", fmt.Errorf("auth: failed to create auth token: %w", err)
	}
	if err := g.cache.Insert(ctx, token); err != nil {
		return "", fmt.Errorf("auth: failed to insert auth token into cache: %w", err)
	}
	return token.RefreshToken, nil
}

func (g *jwtGenerator) RefreshAccessToken(ctx context.Context, refreshToken string) (*Auth, error) {
	if refreshToken == "" {
		return nil, ErrInvalidRefreshToken
	}
	hashedToken, err := hashRefreshToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("auth: failed to hash refresh token: %w", err)
	}
	token := &RefreshToken{
		HashedToken: hashedToken,
	}
	if err := g.cache.Get(ctx, token); err != nil {
		return nil, fmt.Errorf("auth: failed to get auth token from cache: %w", err)
	}
	accessToken, err := g.GenerateAccessToken(token.UserID, token.FacilityID)
	if err != nil {
		return nil, fmt.Errorf("auth: failed to generate access token: %w", err)
	}
	res := &Auth{
		AccessToken:  accessToken,
		RefreshToken: "",
		ExpiresIn:    int32(g.accessTokenTTL.Seconds()),
	}
	return res, nil
}

func (g *jwtGenerator) DeleteRefreshToken(ctx context.Context, userID string) error {
	// ユーザーIDに紐づく全てのリフレッシュトークンを検索
	filter := map[string]interface{}{
		"user_id": userID,
	}
	tokens := RefreshTokens{}
	if err := g.cache.Scan(ctx, tokens, filter); err != nil {
		return fmt.Errorf("auth: failed to scan refresh tokens: %w", err)
	}
	if len(tokens) == 0 {
		return nil
	}

	// バッチ削除を実行
	if err := g.cache.BatchDelete(ctx, tokens); err != nil {
		return fmt.Errorf("auth: failed to batch delete refresh tokens: %w", err)
	}
	return nil
}

type Claims struct {
	jwt.RegisteredClaims
	FacilityID string `json:"facilityId,omitempty"`
}

type JWTVerifier interface {
	VerifyAccessToken(accessToken string) (*Claims, error)
	VerifyRefreshToken(ctx context.Context, refreshToken string) (*RefreshToken, error)
}

type JWTVerifierParams struct {
	Cache      dynamodb.Client
	Issuer     string
	PrivateKey []byte
}

type jwtVerifier struct {
	issuer        string
	secret        *rsa.PrivateKey
	signingMethod *jwt.SigningMethodRSA
	cache         dynamodb.Client
	now           func() time.Time
}

func NewJWTVerifier(params *JWTVerifierParams, opts ...Option) (JWTVerifier, error) {
	dopts := buildOptions(opts...)
	secret, err := parseRSAPrivateKey(params.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("auth: failed to parse private key: %w", err)
	}
	client := &jwtVerifier{
		issuer:        params.Issuer,
		secret:        secret,
		signingMethod: dopts.signingMethod,
		cache:         params.Cache,
		now:           dopts.now,
	}
	return client, nil
}

func (v *jwtVerifier) VerifyAccessToken(accessToken string) (*Claims, error) {
	parser := jwt.NewParser(
		jwt.WithValidMethods([]string{v.signingMethod.Alg()}),
		jwt.WithIssuer(v.issuer),
		jwt.WithIssuedAt(),
	)
	claims := &Claims{}
	keyFunc := func(token *jwt.Token) (any, error) {
		return &v.secret.PublicKey, nil
	}
	token, err := parser.ParseWithClaims(accessToken, claims, keyFunc)
	if err != nil {
		return nil, fmt.Errorf("auth: failed to parse access token: %w", err)
	}
	if !token.Valid {
		return nil, ErrInvalidAccessToken
	}
	return claims, nil
}

func (v *jwtVerifier) VerifyRefreshToken(ctx context.Context, refreshToken string) (*RefreshToken, error) {
	hashedToken, err := hashRefreshToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("auth: failed to hash refresh token: %w", err)
	}
	token := &RefreshToken{
		HashedToken: hashedToken,
	}
	if err := v.cache.Get(ctx, token); err != nil {
		return nil, fmt.Errorf("auth: failed to get auth token from cache: %w", err)
	}
	if err := token.Verify(refreshToken); err != nil {
		return nil, fmt.Errorf("auth: invalid refresh token: %w", err)
	}
	if v.now().After(token.ExpiredAt) {
		return nil, ErrRefreshTokenExpired
	}
	return token, nil
}
