package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddressRevision(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		params *NewAddressRevisionParams
		expect *AddressRevision
		hasErr bool
	}{
		{
			name: "success",
			params: &NewAddressRevisionParams{
				Lastname:       "&.",
				Firstname:      "購入者",
				PostalCode:     "1000014",
				PrefectureCode: 13,
				City:           "千代田区",
				AddressLine1:   "永田町1-7-1",
				AddressLine2:   "",
				PhoneNumber:    "+819012345678",
			},
			expect: &AddressRevision{
				Lastname:       "&.",
				Firstname:      "購入者",
				PostalCode:     "1000014",
				Prefecture:     "東京都",
				PrefectureCode: 13,
				City:           "千代田区",
				AddressLine1:   "永田町1-7-1",
				AddressLine2:   "",
				PhoneNumber:    "+819012345678",
			},
			hasErr: false,
		},
		{
			name: "invalid prefecture",
			params: &NewAddressRevisionParams{
				Lastname:       "&.",
				Firstname:      "購入者",
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
			actual, err := NewAddressRevision(tt.params)
			assert.Equal(t, tt.hasErr, err != nil, err)
			if actual != nil {
				actual.AddressID = "" // ignore
			}
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestAddressRevision_Fill(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		revision *AddressRevision
		expect   *AddressRevision
		hasErr   bool
	}{
		{
			name: "success",
			revision: &AddressRevision{
				Lastname:       "&.",
				Firstname:      "購入者",
				PostalCode:     "1000014",
				Prefecture:     "",
				PrefectureCode: 13,
				City:           "千代田区",
				AddressLine1:   "永田町1-7-1",
				AddressLine2:   "",
				PhoneNumber:    "+819012345678",
			},
			expect: &AddressRevision{
				Lastname:       "&.",
				Firstname:      "購入者",
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
			name: "invalid prefecture",
			revision: &AddressRevision{
				Lastname:       "&.",
				Firstname:      "購入者",
				PostalCode:     "1000014",
				Prefecture:     "",
				PrefectureCode: 0,
				City:           "千代田区",
				AddressLine1:   "永田町1-7-1",
				AddressLine2:   "",
				PhoneNumber:    "+819012345678",
			},
			expect: &AddressRevision{
				Lastname:       "&.",
				Firstname:      "購入者",
				PostalCode:     "1000014",
				Prefecture:     "",
				PrefectureCode: 0,
				City:           "千代田区",
				AddressLine1:   "永田町1-7-1",
				AddressLine2:   "",
				PhoneNumber:    "+819012345678",
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.revision.Fill()
			assert.Equal(t, tt.expect, tt.revision)
		})
	}
}

func TestAddressRevisions_Map(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		revisions AddressRevisions
		expect    map[string]*AddressRevision
		hasErr    bool
	}{
		{
			name: "success",
			revisions: AddressRevisions{
				{
					AddressID:      "address-id",
					Lastname:       "&.",
					Firstname:      "購入者",
					PostalCode:     "1000014",
					Prefecture:     "東京都",
					PrefectureCode: 13,
					City:           "千代田区",
					AddressLine1:   "永田町1-7-1",
					AddressLine2:   "",
					PhoneNumber:    "+819012345678",
				},
			},
			expect: map[string]*AddressRevision{
				"address-id": {
					AddressID:      "address-id",
					Lastname:       "&.",
					Firstname:      "購入者",
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
			assert.Equal(t, tt.expect, tt.revisions.Map())
		})
	}
}

func TestAddressRevisions_Fill(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		revisions AddressRevisions
		expect    AddressRevisions
		hasErr    bool
	}{
		{
			name: "success",
			revisions: AddressRevisions{
				{
					AddressID:      "address-id",
					Lastname:       "&.",
					Firstname:      "購入者",
					PostalCode:     "1000014",
					Prefecture:     "",
					PrefectureCode: 13,
					City:           "千代田区",
					AddressLine1:   "永田町1-7-1",
					AddressLine2:   "",
					PhoneNumber:    "+819012345678",
				},
			},
			expect: AddressRevisions{
				{
					AddressID:      "address-id",
					Lastname:       "&.",
					Firstname:      "購入者",
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
		{
			name: "invalid prefecture",
			revisions: AddressRevisions{
				{
					AddressID:      "address-id",
					Lastname:       "&.",
					Firstname:      "購入者",
					PostalCode:     "1000014",
					Prefecture:     "",
					PrefectureCode: 0,
					City:           "千代田区",
					AddressLine1:   "永田町1-7-1",
					AddressLine2:   "",
					PhoneNumber:    "+819012345678",
				},
			},
			expect: AddressRevisions{
				{
					AddressID:      "address-id",
					Lastname:       "&.",
					Firstname:      "購入者",
					PostalCode:     "1000014",
					Prefecture:     "",
					PrefectureCode: 0,
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
			tt.revisions.Fill()
			assert.Equal(t, tt.expect, tt.revisions)
		})
	}
}
