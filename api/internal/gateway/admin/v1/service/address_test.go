package service

import (
	"testing"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/stretchr/testify/assert"
)

func TestAddress(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		address *entity.Address
		expect  *Address
	}{
		{
			name: "success",
			address: &entity.Address{
				ID:        "address-id",
				UserID:    "user-id",
				IsDefault: true,
				AddressRevision: entity.AddressRevision{
					ID:             1,
					Lastname:       "&.",
					Firstname:      "購入者",
					LastnameKana:   "あんどどっと",
					FirstnameKana:  "こうにゅうしゃ",
					PostalCode:     "1000014",
					Prefecture:     "東京都",
					PrefectureCode: 13,
					City:           "千代田区",
					AddressLine1:   "永田町1-7-1",
					AddressLine2:   "",
					PhoneNumber:    "090-1234-1234",
				},
			},
			expect: &Address{
				Address: response.Address{
					Lastname:       "&.",
					Firstname:      "購入者",
					LastnameKana:   "あんどどっと",
					FirstnameKana:  "こうにゅうしゃ",
					PostalCode:     "1000014",
					Prefecture:     "東京都",
					PrefectureCode: 13,
					City:           "千代田区",
					AddressLine1:   "永田町1-7-1",
					AddressLine2:   "",
					PhoneNumber:    "090-1234-1234",
				},
				revisionID: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewAddress(tt.address)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestAddress_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		address *Address
		expect  *response.Address
	}{
		{
			name: "success",
			address: &Address{
				Address: response.Address{
					Lastname:       "&.",
					Firstname:      "購入者",
					LastnameKana:   "あんどどっと",
					FirstnameKana:  "こうにゅうしゃ",
					PostalCode:     "1000014",
					PrefectureCode: 13,
					City:           "千代田区",
					AddressLine1:   "永田町1-7-1",
					AddressLine2:   "",
					PhoneNumber:    "090-1234-1234",
				},
				revisionID: 1,
			},
			expect: &response.Address{
				Lastname:       "&.",
				Firstname:      "購入者",
				LastnameKana:   "あんどどっと",
				FirstnameKana:  "こうにゅうしゃ",
				PostalCode:     "1000014",
				PrefectureCode: 13,
				City:           "千代田区",
				AddressLine1:   "永田町1-7-1",
				AddressLine2:   "",
				PhoneNumber:    "090-1234-1234",
			},
		},
		{
			name:    "empty",
			address: nil,
			expect:  nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.address.Response())
		})
	}
}

func TestAddresses(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		addresses entity.Addresses
		expect    Addresses
	}{
		{
			name: "success",
			addresses: entity.Addresses{
				{
					ID:        "address-id",
					UserID:    "user-id",
					IsDefault: true,
					AddressRevision: entity.AddressRevision{
						ID:             1,
						Lastname:       "&.",
						Firstname:      "購入者",
						LastnameKana:   "あんどどっと",
						FirstnameKana:  "こうにゅうしゃ",
						PostalCode:     "1000014",
						Prefecture:     "東京都",
						PrefectureCode: 13,
						City:           "千代田区",
						AddressLine1:   "永田町1-7-1",
						AddressLine2:   "",
						PhoneNumber:    "090-1234-1234",
					},
				},
			},
			expect: Addresses{
				{
					Address: response.Address{
						Lastname:       "&.",
						Firstname:      "購入者",
						LastnameKana:   "あんどどっと",
						FirstnameKana:  "こうにゅうしゃ",
						PostalCode:     "1000014",
						Prefecture:     "東京都",
						PrefectureCode: 13,
						City:           "千代田区",
						AddressLine1:   "永田町1-7-1",
						AddressLine2:   "",
						PhoneNumber:    "090-1234-1234",
					},
					revisionID: 1,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewAddresses(tt.addresses)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestAddresses_MapRevision(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		addresses Addresses
		expect    map[int64]*Address
	}{
		{
			name: "success",
			addresses: Addresses{
				{
					Address: response.Address{
						Lastname:       "&.",
						Firstname:      "購入者",
						LastnameKana:   "あんどどっと",
						FirstnameKana:  "こうにゅうしゃ",
						PostalCode:     "1000014",
						PrefectureCode: 13,
						City:           "千代田区",
						AddressLine1:   "永田町1-7-1",
						AddressLine2:   "",
						PhoneNumber:    "090-1234-1234",
					},
					revisionID: 1,
				},
			},
			expect: map[int64]*Address{
				1: {
					Address: response.Address{
						Lastname:       "&.",
						Firstname:      "購入者",
						LastnameKana:   "あんどどっと",
						FirstnameKana:  "こうにゅうしゃ",
						PostalCode:     "1000014",
						PrefectureCode: 13,
						City:           "千代田区",
						AddressLine1:   "永田町1-7-1",
						AddressLine2:   "",
						PhoneNumber:    "090-1234-1234",
					},
					revisionID: 1,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.addresses.MapByRevision())
		})
	}
}

func TestAddresses_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		addresses Addresses
		expect    []*response.Address
	}{
		{
			name: "success",
			addresses: Addresses{
				{
					Address: response.Address{
						Lastname:       "&.",
						Firstname:      "購入者",
						LastnameKana:   "あんどどっと",
						FirstnameKana:  "こうにゅうしゃ",
						PostalCode:     "1000014",
						PrefectureCode: 13,
						City:           "千代田区",
						AddressLine1:   "永田町1-7-1",
						AddressLine2:   "",
						PhoneNumber:    "090-1234-1234",
					},
					revisionID: 1,
				},
			},
			expect: []*response.Address{
				{
					Lastname:       "&.",
					Firstname:      "購入者",
					LastnameKana:   "あんどどっと",
					FirstnameKana:  "こうにゅうしゃ",
					PostalCode:     "1000014",
					PrefectureCode: 13,
					City:           "千代田区",
					AddressLine1:   "永田町1-7-1",
					AddressLine2:   "",
					PhoneNumber:    "090-1234-1234",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.addresses.Response())
		})
	}
}
