package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddress(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		params *NewAddressParams
		expect *Address
		hasErr bool
	}{
		{
			name: "success",
			params: &NewAddressParams{
				UserID:         "user-id",
				IsDefault:      true,
				Lastname:       "&.",
				Firstname:      "購入者",
				LastnameKana:   "あんどどっと",
				FirstnameKana:  "こうにゅうしゃ",
				PostalCode:     "1000014",
				PrefectureCode: 13,
				City:           "千代田区",
				AddressLine1:   "永田町1-7-1",
				AddressLine2:   "",
				PhoneNumber:    "+819012345678",
			},
			expect: &Address{
				UserID:    "user-id",
				IsDefault: true,
				AddressRevision: AddressRevision{
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
					PhoneNumber:    "+819012345678",
				},
			},
			hasErr: false,
		},
		{
			name: "invalid prefecture",
			params: &NewAddressParams{
				UserID:         "user-id",
				IsDefault:      true,
				Lastname:       "&.",
				Firstname:      "購入者",
				LastnameKana:   "あんどどっと",
				FirstnameKana:  "こうにゅうしゃ",
				PostalCode:     "1000014",
				PrefectureCode: 0,
				City:           "千代田区",
				AddressLine1:   "永田町1-7-1",
				AddressLine2:   "",
				PhoneNumber:    "+819012345678",
			},
			expect: nil,
			hasErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual, err := NewAddress(tt.params)
			assert.Equal(t, tt.hasErr, err != nil, err)
			if actual != nil {
				actual.ID = ""                        // ignore
				actual.AddressRevision.AddressID = "" // ignore
			}
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestAddress_Name(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		address *Address
		expect  string
	}{
		{
			name: "success",
			address: &Address{
				ID:        "address-id",
				UserID:    "user-id",
				IsDefault: true,
				AddressRevision: AddressRevision{
					AddressID:      "address-id",
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
					PhoneNumber:    "+819012345678",
				},
			},
			expect: "&. 購入者",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.address.Name())
		})
	}
}

func TestAddress_String(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		address *Address
		expect  string
	}{
		{
			name: "success",
			address: &Address{
				ID:        "address-id",
				UserID:    "user-id",
				IsDefault: true,
				AddressRevision: AddressRevision{
					AddressID:      "address-id",
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
					PhoneNumber:    "+819012345678",
				},
			},
			expect: "東京都 千代田区 永田町1-7-1",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.address.String())
		})
	}
}

func TestAddress_Fill(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		address  *Address
		revision *AddressRevision
		expect   *Address
	}{
		{
			name: "success",
			address: &Address{
				ID:        "address-id",
				UserID:    "user-id",
				IsDefault: true,
			},
			revision: &AddressRevision{
				AddressID:      "address-id",
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
				PhoneNumber:    "+819012345678",
			},
			expect: &Address{
				ID:        "address-id",
				UserID:    "user-id",
				IsDefault: true,
				AddressRevision: AddressRevision{
					AddressID:      "address-id",
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
					PhoneNumber:    "+819012345678",
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.address.Fill(tt.revision)
			assert.Equal(t, tt.expect, tt.address)
		})
	}
}

func TestAddresses_IDs(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		addresses Addresses
		expect    []string
	}{
		{
			name: "success",
			addresses: Addresses{
				{
					ID:        "address-id",
					UserID:    "user-id",
					IsDefault: true,
				},
			},
			expect: []string{"address-id"},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.addresses.IDs())
		})
	}
}

func TestAddresses_Map(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		addresses Addresses
		expect    map[string]*Address
	}{
		{
			name: "success",
			addresses: Addresses{
				{
					ID:        "address-id01",
					UserID:    "user-id",
					IsDefault: true,
				},
				{
					ID:        "address-id02",
					UserID:    "user-id",
					IsDefault: false,
				},
			},
			expect: map[string]*Address{
				"address-id01": {
					ID:        "address-id01",
					UserID:    "user-id",
					IsDefault: true,
				},
				"address-id02": {
					ID:        "address-id02",
					UserID:    "user-id",
					IsDefault: false,
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.addresses.Map())
		})
	}
}

func TestAddresses_MapByUserID(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		addresses Addresses
		expect    map[string]*Address
	}{
		{
			name: "success",
			addresses: Addresses{
				{
					ID:        "address-id01",
					UserID:    "user-id01",
					IsDefault: true,
				},
				{
					ID:        "address-id02",
					UserID:    "user-id02",
					IsDefault: false,
				},
			},
			expect: map[string]*Address{
				"user-id01": {
					ID:        "address-id01",
					UserID:    "user-id01",
					IsDefault: true,
				},
				"user-id02": {
					ID:        "address-id02",
					UserID:    "user-id02",
					IsDefault: false,
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.addresses.MapByUserID())
		})
	}
}

func TestAddresses_Fill(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		addresses Addresses
		revisions map[string]*AddressRevision
		expect    Addresses
	}{
		{
			name: "success",
			addresses: Addresses{
				{
					ID:        "address-id01",
					UserID:    "user-id",
					IsDefault: true,
				},
				{
					ID:        "address-id02",
					UserID:    "user-id",
					IsDefault: false,
				},
			},
			revisions: map[string]*AddressRevision{
				"address-id01": {
					AddressID:      "address-id01",
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
					PhoneNumber:    "+819012345678",
				},
			},
			expect: Addresses{
				{
					ID:        "address-id01",
					UserID:    "user-id",
					IsDefault: true,
					AddressRevision: AddressRevision{
						AddressID:      "address-id01",
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
						PhoneNumber:    "+819012345678",
					},
				},
				{
					ID:        "address-id02",
					UserID:    "user-id",
					IsDefault: false,
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.addresses.Fill(tt.revisions)
			assert.Equal(t, tt.expect, tt.addresses)
		})
	}
}
