package auth

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"testing"
	"time"

	mock_dynamodb "github.com/and-period/furumaru/api/mock/pkg/dynamodb"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestAuth(t *testing.T) {
	t.Parallel()
	auth := &Auth{
		AccessToken:  "access-token",
		RefreshToken: "refresh-token",
		ExpiresIn:    3600,
	}
	assert.Equal(t, "access-token", auth.AccessToken)
	assert.Equal(t, "refresh-token", auth.RefreshToken)
	assert.Equal(t, int32(3600), auth.ExpiresIn)
}

func TestJWTGeneratorParams(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	params := &JWTGeneratorParams{
		Cache:      mock_dynamodb.NewMockClient(ctrl),
		Issuer:     "https://example.com",
		PrivateKey: []byte("secret"),
	}
	assert.NotNil(t, params.Cache)
	assert.Equal(t, "https://example.com", params.Issuer)
	assert.Equal(t, []byte("secret"), params.PrivateKey)
}

func TestNewJWTGenerator(t *testing.T) {
	t.Parallel()
	privateKey, _ := generateRSAKeyPair()
	privatePEM := encodeRSAPrivateKey(privateKey)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name    string
		params  *JWTGeneratorParams
		opts    []Option
		wantErr bool
	}{
		{
			name: "success",
			params: &JWTGeneratorParams{
				Cache:      mock_dynamodb.NewMockClient(ctrl),
				Issuer:     "https://example.com",
				PrivateKey: privatePEM,
			},
			opts: []Option{
				WithAccessTokenTTL(time.Hour),
				WithRefreshTokenTTL(24 * time.Hour),
			},
			wantErr: false,
		},
		{
			name: "invalid secret",
			params: &JWTGeneratorParams{
				Cache:      mock_dynamodb.NewMockClient(ctrl),
				Issuer:     "https://example.com",
				PrivateKey: []byte("invalid"),
			},
			opts:    []Option{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			gen, err := NewJWTGenerator(tt.params, tt.opts...)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, gen)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, gen)
			}
		})
	}
}

func TestJWTGenerator_Generate(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	privateKey, _ := generateRSAKeyPair()
	privatePEM := encodeRSAPrivateKey(privateKey)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCache := mock_dynamodb.NewMockClient(ctrl)
	mockCache.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil)

	gen, err := NewJWTGenerator(
		&JWTGeneratorParams{
			Cache:      mockCache,
			Issuer:     "https://example.com",
			PrivateKey: privatePEM,
		},
		WithAccessTokenTTL(time.Hour),
		WithRefreshTokenTTL(24*time.Hour),
	)
	require.NoError(t, err)

	auth, err := gen.Generate(ctx, "user-id", "facility-id")
	assert.NoError(t, err)
	assert.NotEmpty(t, auth.AccessToken)
	assert.NotEmpty(t, auth.RefreshToken)
	assert.Equal(t, int32(3600), auth.ExpiresIn)
}

func TestJWTGenerator_GenerateAccessToken(t *testing.T) {
	t.Parallel()
	privateKey, _ := generateRSAKeyPair()
	privatePEM := encodeRSAPrivateKey(privateKey)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	gen, err := NewJWTGenerator(
		&JWTGeneratorParams{
			Cache:      mock_dynamodb.NewMockClient(ctrl),
			Issuer:     "https://example.com",
			PrivateKey: privatePEM,
		},
		WithAccessTokenTTL(time.Hour),
		WithRefreshTokenTTL(24*time.Hour),
	)
	require.NoError(t, err)

	token, err := gen.GenerateAccessToken("user-id", "facility-id")
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Parse and verify the token
	parsed, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return &privateKey.PublicKey, nil
	})
	assert.NoError(t, err)
	assert.True(t, parsed.Valid)

	claims, ok := parsed.Claims.(*Claims)
	assert.True(t, ok)
	assert.Equal(t, "user-id", claims.Subject)
	assert.Equal(t, "https://example.com", claims.Issuer)
}

func TestJWTGenerator_GenerateRefreshToken(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	privateKey, _ := generateRSAKeyPair()
	privatePEM := encodeRSAPrivateKey(privateKey)

	tests := []struct {
		name       string
		sub        string
		facilityID string
		setupMock  func(*mock_dynamodb.MockClient)
		wantErr    bool
	}{
		{
			name:       "success",
			sub:        "user-id",
			facilityID: "facility-id",
			setupMock: func(m *mock_dynamodb.MockClient) {
				m.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil)
			},
			wantErr: false,
		},
		{
			name:       "cache insert error",
			sub:        "user-id",
			facilityID: "facility-id",
			setupMock: func(m *mock_dynamodb.MockClient) {
				m.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(errors.New("cache error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockCache := mock_dynamodb.NewMockClient(ctrl)
			tt.setupMock(mockCache)

			gen, err := NewJWTGenerator(
				&JWTGeneratorParams{
					Cache:      mockCache,
					Issuer:     "https://example.com",
					PrivateKey: privatePEM,
				},
				WithAccessTokenTTL(time.Hour),
				WithRefreshTokenTTL(24*time.Hour),
			)
			require.NoError(t, err)

			token, err := gen.GenerateRefreshToken(ctx, tt.sub, tt.facilityID)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, token)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)
			}
		})
	}
}

func TestJWTGenerator_RefreshAccessToken(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	privateKey, _ := generateRSAKeyPair()
	privatePEM := encodeRSAPrivateKey(privateKey)

	tests := []struct {
		name         string
		refreshToken string
		setupMock    func(*mock_dynamodb.MockClient)
		wantErr      bool
	}{
		{
			name:         "success",
			refreshToken: "valid-refresh-token",
			setupMock: func(m *mock_dynamodb.MockClient) {
				hashedToken, _ := hashRefreshToken("valid-refresh-token")
				m.EXPECT().Get(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, token *RefreshToken) error {
					token.HashedToken = hashedToken
					token.UserID = "user-id"
					token.FacilityID = "facility-id"
					token.ExpiredAt = time.Now().Add(time.Hour)
					return nil
				})
			},
			wantErr: false,
		},
		{
			name:         "empty refresh token",
			refreshToken: "",
			setupMock:    func(m *mock_dynamodb.MockClient) {},
			wantErr:      true,
		},
		{
			name:         "cache get error",
			refreshToken: "valid-refresh-token",
			setupMock: func(m *mock_dynamodb.MockClient) {
				m.EXPECT().Get(gomock.Any(), gomock.Any()).Return(errors.New("cache error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockCache := mock_dynamodb.NewMockClient(ctrl)
			tt.setupMock(mockCache)

			gen, err := NewJWTGenerator(
				&JWTGeneratorParams{
					Cache:      mockCache,
					Issuer:     "https://example.com",
					PrivateKey: privatePEM,
				},
				WithAccessTokenTTL(time.Hour),
				WithRefreshTokenTTL(24*time.Hour),
			)
			require.NoError(t, err)

			auth, err := gen.RefreshAccessToken(ctx, tt.refreshToken)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, auth)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, auth)
				assert.NotEmpty(t, auth.AccessToken)
				assert.Empty(t, auth.RefreshToken) // RefreshAccessTokenではRefreshTokenは空
				assert.Equal(t, int32(3600), auth.ExpiresIn) // 1時間 = 3600秒
			}
		})
	}
}

func TestJWTGenerator_DeleteRefreshToken(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	privateKey, _ := generateRSAKeyPair()
	privatePEM := encodeRSAPrivateKey(privateKey)

	tests := []struct {
		name      string
		userID    string
		setupMock func(*mock_dynamodb.MockClient)
		wantErr   bool
	}{
		{
			name:   "success",
			userID: "user-id",
			setupMock: func(m *mock_dynamodb.MockClient) {
				m.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				m.EXPECT().BatchDelete(gomock.Any(), gomock.Any()).Return(nil)
			},
			wantErr: false,
		},
		{
			name:   "scan error",
			userID: "user-id",
			setupMock: func(m *mock_dynamodb.MockClient) {
				m.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("scan error"))
			},
			wantErr: true,
		},
		{
			name:   "batch delete error",
			userID: "user-id",
			setupMock: func(m *mock_dynamodb.MockClient) {
				m.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				m.EXPECT().BatchDelete(gomock.Any(), gomock.Any()).Return(errors.New("batch delete error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockCache := mock_dynamodb.NewMockClient(ctrl)
			tt.setupMock(mockCache)

			gen, err := NewJWTGenerator(
				&JWTGeneratorParams{
					Cache:      mockCache,
					Issuer:     "https://example.com",
					PrivateKey: privatePEM,
				},
				WithAccessTokenTTL(time.Hour),
				WithRefreshTokenTTL(24*time.Hour),
			)
			require.NoError(t, err)

			err = gen.DeleteRefreshToken(ctx, tt.userID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestClaims(t *testing.T) {
	t.Parallel()
	claims := &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:  "https://example.com",
			Subject: "user-id",
		},
		FacilityID: "facility-id",
	}
	assert.Equal(t, "https://example.com", claims.Issuer)
	assert.Equal(t, "user-id", claims.Subject)
	assert.Equal(t, "facility-id", claims.FacilityID)
}

func TestJWTVerifierParams(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	params := &JWTVerifierParams{
		Cache:      mock_dynamodb.NewMockClient(ctrl),
		Issuer:     "https://example.com",
		PrivateKey: []byte("secret"),
	}
	assert.NotNil(t, params.Cache)
	assert.Equal(t, "https://example.com", params.Issuer)
	assert.Equal(t, []byte("secret"), params.PrivateKey)
}

func TestNewJWTVerifier(t *testing.T) {
	t.Parallel()
	privateKey, _ := generateRSAKeyPair()
	privatePEM := encodeRSAPrivateKey(privateKey)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name    string
		params  *JWTVerifierParams
		opts    []Option
		wantErr bool
	}{
		{
			name: "success",
			params: &JWTVerifierParams{
				Cache:      mock_dynamodb.NewMockClient(ctrl),
				Issuer:     "https://example.com",
				PrivateKey: privatePEM,
			},
			opts:    []Option{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			verifier, err := NewJWTVerifier(tt.params, tt.opts...)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, verifier)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, verifier)
			}
		})
	}
}

func TestJWTVerifier_VerifyAccessToken(t *testing.T) {
	t.Parallel()
	privateKey, _ := generateRSAKeyPair()
	privatePEM := encodeRSAPrivateKey(privateKey)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	verifier, err := NewJWTVerifier(&JWTVerifierParams{
		Cache:      mock_dynamodb.NewMockClient(ctrl),
		Issuer:     "https://example.com",
		PrivateKey: privatePEM,
	})
	require.NoError(t, err)

	// Generate a valid token
	now := time.Now()
	claims := &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "https://example.com",
			Subject:   "user-id",
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour)),
			ID:        "token-id",
		},
		FacilityID: "facility-id",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenString, err := token.SignedString(privateKey)
	require.NoError(t, err)

	tests := []struct {
		name        string
		accessToken string
		wantErr     bool
		wantClaims  *Claims
	}{
		{
			name:        "valid token",
			accessToken: tokenString,
			wantErr:     false,
			wantClaims:  claims,
		},
		{
			name:        "invalid token",
			accessToken: "invalid.token.here",
			wantErr:     true,
			wantClaims:  nil,
		},
		{
			name: "wrong issuer",
			accessToken: func() string {
				wrongClaims := &Claims{
					RegisteredClaims: jwt.RegisteredClaims{
						Issuer:    "https://wrong.com",
						Subject:   "user-id",
						IssuedAt:  jwt.NewNumericDate(now),
						ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour)),
					},
				}
				wrongToken := jwt.NewWithClaims(jwt.SigningMethodRS256, wrongClaims)
				wrongTokenString, _ := wrongToken.SignedString(privateKey)
				return wrongTokenString
			}(),
			wantErr:    true,
			wantClaims: nil,
		},
		{
			name: "expired token",
			accessToken: func() string {
				expiredClaims := &Claims{
					RegisteredClaims: jwt.RegisteredClaims{
						Issuer:    "https://example.com",
						Subject:   "user-id",
						IssuedAt:  jwt.NewNumericDate(now.Add(-2 * time.Hour)),
						ExpiresAt: jwt.NewNumericDate(now.Add(-time.Hour)),
					},
				}
				expiredToken := jwt.NewWithClaims(jwt.SigningMethodRS256, expiredClaims)
				expiredTokenString, _ := expiredToken.SignedString(privateKey)
				return expiredTokenString
			}(),
			wantErr:    true,
			wantClaims: nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result, err := verifier.VerifyAccessToken(tt.accessToken)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.wantClaims.Subject, result.Subject)
				assert.Equal(t, tt.wantClaims.Issuer, result.Issuer)
				assert.Equal(t, tt.wantClaims.FacilityID, result.FacilityID)
			}
		})
	}
}

func TestJWTVerifier_VerifyRefreshToken(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	privateKey, _ := generateRSAKeyPair()
	privatePEM := encodeRSAPrivateKey(privateKey)

	now := time.Now()
	validToken := &RefreshToken{
		RefreshToken: "valid-refresh-token",
		HashedToken:  "$2a$10$validhash",
		UserID:       "user-id",
		FacilityID:   "facility-id",
		ExpiredAt:    now.Add(time.Hour),
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	expiredToken := &RefreshToken{
		RefreshToken: "expired-refresh-token",
		HashedToken:  "$2a$10$expiredhash",
		UserID:       "user-id",
		FacilityID:   "facility-id",
		ExpiredAt:    now.Add(-time.Hour),
		CreatedAt:    now.Add(-2 * time.Hour),
		UpdatedAt:    now.Add(-2 * time.Hour),
	}

	tests := []struct {
		name         string
		refreshToken string
		setupMock    func(*mock_dynamodb.MockClient)
		setupTime    func() time.Time
		wantErr      bool
	}{
		{
			name:         "success",
			refreshToken: "valid-refresh-token",
			setupMock: func(m *mock_dynamodb.MockClient) {
				m.EXPECT().Get(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, token *RefreshToken) error {
					*token = *validToken
					token.HashedToken, _ = hashRefreshToken("valid-refresh-token")
					return nil
				})
			},
			setupTime: func() time.Time { return now },
			wantErr:   false,
		},
		{
			name:         "cache get error",
			refreshToken: "valid-refresh-token",
			setupMock: func(m *mock_dynamodb.MockClient) {
				m.EXPECT().Get(gomock.Any(), gomock.Any()).Return(errors.New("cache error"))
			},
			setupTime: func() time.Time { return now },
			wantErr:   true,
		},
		{
			name:         "expired token",
			refreshToken: "expired-refresh-token",
			setupMock: func(m *mock_dynamodb.MockClient) {
				m.EXPECT().Get(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, token *RefreshToken) error {
					*token = *expiredToken
					token.HashedToken, _ = hashRefreshToken("expired-refresh-token")
					return nil
				})
			},
			setupTime: func() time.Time { return now },
			wantErr:   true,
		},
		{
			name:         "invalid token",
			refreshToken: "wrong-refresh-token",
			setupMock: func(m *mock_dynamodb.MockClient) {
				m.EXPECT().Get(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, token *RefreshToken) error {
					*token = *validToken
					token.HashedToken, _ = hashRefreshToken("valid-refresh-token")
					return nil
				})
			},
			setupTime: func() time.Time { return now },
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockCache := mock_dynamodb.NewMockClient(ctrl)
			tt.setupMock(mockCache)

			verifier, err := NewJWTVerifier(&JWTVerifierParams{
				Cache:      mockCache,
				Issuer:     "https://example.com",
				PrivateKey: privatePEM,
			})
			require.NoError(t, err)

			v := verifier.(*jwtVerifier)
			v.now = tt.setupTime

			result, err := verifier.VerifyRefreshToken(ctx, tt.refreshToken)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, "user-id", result.UserID)
				assert.Equal(t, "facility-id", result.FacilityID)
			}
		})
	}
}

// Helper functions for tests
func generateRSAKeyPair() (*rsa.PrivateKey, error) {
	return rsa.GenerateKey(rand.Reader, 2048)
}

func encodeRSAPrivateKey(key *rsa.PrivateKey) []byte {
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(key)
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	}
	return pem.EncodeToMemory(block)
}

func encodeRSAPublicKey(key *rsa.PublicKey) []byte {
	publicKeyBytes, _ := x509.MarshalPKIXPublicKey(key)
	block := &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyBytes,
	}
	return pem.EncodeToMemory(block)
}

func encodeRSAPublicKeyPKCS1(key *rsa.PublicKey) []byte {
	publicKeyBytes := x509.MarshalPKCS1PublicKey(key)
	block := &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyBytes,
	}
	return pem.EncodeToMemory(block)
}

func TestOptions(t *testing.T) {
	t.Parallel()

	t.Run("WithAccessTokenTTL", func(t *testing.T) {
		t.Parallel()
		ttl := 2 * time.Hour
		opt := WithAccessTokenTTL(ttl)
		opts := &options{}
		opt(opts)
		assert.Equal(t, ttl, opts.accessTokenTTL)
	})

	t.Run("WithRefreshTokenTTL", func(t *testing.T) {
		t.Parallel()
		ttl := 7 * 24 * time.Hour
		opt := WithRefreshTokenTTL(ttl)
		opts := &options{}
		opt(opts)
		assert.Equal(t, ttl, opts.refreshTokenTTL)
	})

	t.Run("WithSigningMethod", func(t *testing.T) {
		t.Parallel()
		method := jwt.SigningMethodRS256
		opt := WithSigningMethod(method)
		opts := &options{}
		opt(opts)
		assert.Equal(t, method, opts.signingMethod)
	})

	t.Run("WithNow", func(t *testing.T) {
		t.Parallel()
		now := func() time.Time { return time.Unix(1234567890, 0) }
		opt := WithNow(now)
		opts := &options{}
		opt(opts)
		assert.NotNil(t, opts.now)
		assert.Equal(t, time.Unix(1234567890, 0), opts.now())
	})

	t.Run("WithGenerateID", func(t *testing.T) {
		t.Parallel()
		genID := func() string { return "test-id" }
		opt := WithGenerateID(genID)
		opts := &options{}
		opt(opts)
		assert.NotNil(t, opts.generateID)
		assert.Equal(t, "test-id", opts.generateID())
	})
}

func TestBuildOptions(t *testing.T) {
	t.Parallel()

	t.Run("default options", func(t *testing.T) {
		t.Parallel()
		opts := buildOptions()
		assert.Equal(t, 30*time.Minute, opts.accessTokenTTL)
		assert.Equal(t, 14*24*time.Hour, opts.refreshTokenTTL)
		assert.Equal(t, jwt.SigningMethodRS256, opts.signingMethod)
		assert.NotNil(t, opts.now)
		assert.NotNil(t, opts.generateID)
	})

	t.Run("with custom options", func(t *testing.T) {
		t.Parallel()
		customTTL := 1 * time.Hour
		customMethod := jwt.SigningMethodRS256
		opts := buildOptions(
			WithAccessTokenTTL(customTTL),
			WithSigningMethod(customMethod),
		)
		assert.Equal(t, customTTL, opts.accessTokenTTL)
		assert.Equal(t, 14*24*time.Hour, opts.refreshTokenTTL)
		assert.Equal(t, customMethod, opts.signingMethod)
	})
}
