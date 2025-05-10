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
				PhoneNumber:    "090-1234-1234",
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
					PhoneNumber:    "090-1234-1234",
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
				PhoneNumber:    "090-1234-1234",
			},
			expect: nil,
			hasErr: true,
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual, err := NewAddress(tt.params)
			assert.Equal(t, tt.hasErr, err != nil, err)
			if actual != nil {
				actual.ID = ""        // ignore
				actual.AddressID = "" // ignore
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
					PhoneNumber:    "090-1234-1234",
				},
			},
			expect: "&. 購入者",
		},
		{
			name: "success only lastname",
			address: &Address{
				ID:        "address-id",
				UserID:    "user-id",
				IsDefault: true,
				AddressRevision: AddressRevision{
					AddressID:      "address-id",
					Lastname:       "&.",
					Firstname:      "",
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
			expect: "&.",
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.address.Name())
		})
	}
}

func TestAddress_NameKana(t *testing.T) {
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
					PhoneNumber:    "090-1234-1234",
				},
			},
			expect: "あんどどっと こうにゅうしゃ",
		},
		{
			name: "success only lastname",
			address: &Address{
				ID:        "address-id",
				UserID:    "user-id",
				IsDefault: true,
				AddressRevision: AddressRevision{
					AddressID:      "address-id",
					Lastname:       "&.",
					Firstname:      "",
					LastnameKana:   "あんどどっと",
					FirstnameKana:  "",
					PostalCode:     "1000014",
					Prefecture:     "東京都",
					PrefectureCode: 13,
					City:           "千代田区",
					AddressLine1:   "永田町1-7-1",
					AddressLine2:   "",
					PhoneNumber:    "090-1234-1234",
				},
			},
			expect: "あんどどっと",
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.address.NameKana())
		})
	}
}

func TestAddress_FullPath(t *testing.T) {
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
					PhoneNumber:    "090-1234-1234",
				},
			},
			expect: "東京都 千代田区 永田町1-7-1",
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.address.FullPath())
		})
	}
}

func TestAddress_ShortPath(t *testing.T) {
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
					PhoneNumber:    "090-1234-1234",
				},
			},
			expect: "東京都 千代田区 永田町1-7-1",
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.address.ShortPath())
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
				PhoneNumber:    "090-1234-1234",
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
					PhoneNumber:    "090-1234-1234",
				},
			},
		},
	}
	for _, tt := range tests {

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
					AddressRevision: AddressRevision{
						ID:        1,
						AddressID: "address-id01",
					},
					ID:        "address-id01",
					UserID:    "user-id",
					IsDefault: true,
				},
				{
					AddressRevision: AddressRevision{
						ID:        2,
						AddressID: "address-id02",
					},
					ID:        "address-id02",
					UserID:    "user-id",
					IsDefault: false,
				},
			},
			expect: map[string]*Address{
				"address-id01": {
					AddressRevision: AddressRevision{
						ID:        1,
						AddressID: "address-id01",
					},
					ID:        "address-id01",
					UserID:    "user-id",
					IsDefault: true,
				},
				"address-id02": {
					AddressRevision: AddressRevision{
						ID:        2,
						AddressID: "address-id02",
					},
					ID:        "address-id02",
					UserID:    "user-id",
					IsDefault: false,
				},
			},
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.addresses.Map())
		})
	}
}

func TestAddresses_MapByRevision(t *testing.T) {
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
					AddressRevision: AddressRevision{
						ID:        1,
						AddressID: "address-id01",
					},
					ID:        "address-id01",
					UserID:    "user-id",
					IsDefault: true,
				},
				{
					AddressRevision: AddressRevision{
						ID:        2,
						AddressID: "address-id02",
					},
					ID:        "address-id02",
					UserID:    "user-id",
					IsDefault: false,
				},
			},
			expect: map[int64]*Address{
				1: {
					AddressRevision: AddressRevision{
						ID:        1,
						AddressID: "address-id01",
					},
					ID:        "address-id01",
					UserID:    "user-id",
					IsDefault: true,
				},
				2: {
					AddressRevision: AddressRevision{
						ID:        2,
						AddressID: "address-id02",
					},
					ID:        "address-id02",
					UserID:    "user-id",
					IsDefault: false,
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
						PhoneNumber:    "090-1234-1234",
					},
				},
			},
			expect: map[string]*Address{
				"user-id": {
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
						PhoneNumber:    "090-1234-1234",
					},
				},
			},
		},
	}
	for _, tt := range tests {

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
					PhoneNumber:    "090-1234-1234",
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
						PhoneNumber:    "090-1234-1234",
					},
				},
				{
					ID:        "address-id02",
					UserID:    "user-id",
					IsDefault: false,
					AddressRevision: AddressRevision{
						AddressID: "address-id02",
					},
				},
			},
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.addresses.Fill(tt.revisions)
			assert.Equal(t, tt.expect, tt.addresses)
		})
	}
}
