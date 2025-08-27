package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRefreshToken(t *testing.T) {
	t.Parallel()
	now := time.Now()
	token := &RefreshToken{
		RefreshToken: "refresh-token",
		HashedToken:  "hashed-token",
		UserID:       "user-id",
		FacilityID:   "facility-id",
		ExpiredAt:    now.Add(time.Hour),
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	assert.Equal(t, "refresh-token", token.RefreshToken)
	assert.Equal(t, "hashed-token", token.HashedToken)
	assert.Equal(t, "user-id", token.UserID)
	assert.Equal(t, "facility-id", token.FacilityID)
	assert.Equal(t, now.Add(time.Hour), token.ExpiredAt)
	assert.Equal(t, now, token.CreatedAt)
	assert.Equal(t, now, token.UpdatedAt)
}

func TestRefreshToken_TableName(t *testing.T) {
	t.Parallel()
	token := &RefreshToken{}
	assert.Equal(t, "auth-tokens", token.TableName())
}

func TestRefreshToken_PrimaryKey(t *testing.T) {
	t.Parallel()
	token := &RefreshToken{
		HashedToken: "hashed-token",
	}
	expected := map[string]interface{}{
		"hashed_token": "hashed-token",
	}
	assert.Equal(t, expected, token.PrimaryKey())
}

func TestRefreshToken_Verify(t *testing.T) {
	t.Parallel()

	// Generate a valid refresh token
	refreshToken := newRefreshToken("test-token")
	hashedToken, err := hashRefreshToken(refreshToken)
	require.NoError(t, err)

	tests := []struct {
		name         string
		token        *RefreshToken
		refreshToken string
		wantErr      bool
	}{
		{
			name: "valid token",
			token: &RefreshToken{
				HashedToken: hashedToken,
			},
			refreshToken: refreshToken,
			wantErr:      false,
		},
		{
			name: "invalid token",
			token: &RefreshToken{
				HashedToken: hashedToken,
			},
			refreshToken: "wrong-token",
			wantErr:      true,
		},
		{
			name: "empty refresh token",
			token: &RefreshToken{
				HashedToken: hashedToken,
			},
			refreshToken: "",
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.token.Verify(tt.refreshToken)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNewRefreshToken(t *testing.T) {
	t.Parallel()
	now := time.Now()

	tests := []struct {
		name    string
		params  *RefreshTokenParams
		wantErr bool
	}{
		{
			name: "success",
			params: &RefreshTokenParams{
				UserID:     "user-id",
				FacilityID: "facility-id",
				Now:        now,
				TTL:        time.Hour,
			},
			wantErr: false,
		},
		{
			name: "empty user id",
			params: &RefreshTokenParams{
				UserID:     "",
				FacilityID: "facility-id",
				Now:        now,
				TTL:        time.Hour,
			},
			wantErr: false,
		},
		{
			name: "empty facility id",
			params: &RefreshTokenParams{
				UserID:     "user-id",
				FacilityID: "",
				Now:        now,
				TTL:        time.Hour,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			token, err := NewRefreshToken(tt.params)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, token)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, token)
				assert.NotEmpty(t, token.RefreshToken)
				assert.NotEmpty(t, token.HashedToken)
				assert.Equal(t, tt.params.UserID, token.UserID)
				assert.Equal(t, tt.params.FacilityID, token.FacilityID)
				assert.Equal(t, tt.params.Now.Add(tt.params.TTL), token.ExpiredAt)
				assert.Equal(t, tt.params.Now, token.CreatedAt)
				assert.Equal(t, tt.params.Now, token.UpdatedAt)
			}
		})
	}
}

func TestCompareRefreshToken(t *testing.T) {
	t.Parallel()
	token := "test-token"
	hashedToken, err := hashRefreshToken(token)
	require.NoError(t, err)

	tests := []struct {
		name    string
		hashed  string
		raw     string
		wantErr bool
	}{
		{
			name:    "valid token",
			hashed:  hashedToken,
			raw:     token,
			wantErr: false,
		},
		{
			name:    "invalid token",
			hashed:  hashedToken,
			raw:     "wrong-token",
			wantErr: true,
		},
		{
			name:    "empty token",
			hashed:  hashedToken,
			raw:     "",
			wantErr: true,
		},
		{
			name:    "invalid hash",
			hashed:  "invalid-hash",
			raw:     token,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := compareRefreshToken(tt.hashed, tt.raw)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestParseRSAPrivateKeyFromPEM(t *testing.T) {
	t.Parallel()

	// Generate test keys
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	// PKCS1 format
	pkcs1Bytes := x509.MarshalPKCS1PrivateKey(privateKey)
	pkcs1PEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: pkcs1Bytes,
	})

	// PKCS8 format
	pkcs8Bytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	require.NoError(t, err)
	pkcs8PEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: pkcs8Bytes,
	})

	tests := []struct {
		name    string
		pemData []byte
		wantErr bool
		errType error
	}{
		{
			name:    "valid PKCS1 private key",
			pemData: pkcs1PEM,
			wantErr: false,
		},
		{
			name:    "valid PKCS8 private key",
			pemData: pkcs8PEM,
			wantErr: false,
		},
		{
			name:    "invalid PEM",
			pemData: []byte("invalid pem data"),
			wantErr: true,
			errType: ErrNoPemBlock,
		},
		{
			name:    "empty PEM",
			pemData: []byte{},
			wantErr: true,
			errType: ErrNoPemBlock,
		},
		{
			name: "non-RSA key",
			pemData: pem.EncodeToMemory(&pem.Block{
				Type:  "PRIVATE KEY",
				Bytes: []byte("invalid key data"),
			}),
			wantErr: true,
			errType: ErrNotRSAPrivateKey,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			key, err := parseRSAPrivateKeyFromPEM(tt.pemData)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, key)
				if tt.errType != nil {
					assert.ErrorIs(t, err, tt.errType)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, key)
			}
		})
	}
}

func TestErrors(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "auth: no pem block found", ErrNoPemBlock.Error())
	assert.Equal(t, "auth: not rsa private key", ErrNotRSAPrivateKey.Error())
}
