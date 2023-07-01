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
				UserID:       "user-id",
				IsDefault:    true,
				Lastname:     "&.",
				Firstname:    "購入者",
				PostalCode:   "1000014",
				Prefecture:   13,
				City:         "千代田区",
				AddressLine1: "永田町1-7-1",
				AddressLine2: "",
				PhoneNumber:  "+819012345678",
			},
			expect: &Address{
				UserID:       "user-id",
				Hash:         "789ef22a79a364f95c66a3d3b1fda213c1316a6c7f8b6306b493d8c46d2dce75",
				IsDefault:    true,
				Lastname:     "&.",
				Firstname:    "購入者",
				PostalCode:   "1000014",
				Prefecture:   13,
				City:         "千代田区",
				AddressLine1: "永田町1-7-1",
				AddressLine2: "",
				PhoneNumber:  "+819012345678",
			},
			hasErr: false,
		},
		{
			name: "invalid prefecture",
			params: &NewAddressParams{
				UserID:       "user-id",
				IsDefault:    true,
				Lastname:     "&.",
				Firstname:    "購入者",
				PostalCode:   "1000014",
				Prefecture:   0,
				City:         "千代田区",
				AddressLine1: "永田町1-7-1",
				AddressLine2: "",
				PhoneNumber:  "+819012345678",
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
				actual.ID = "" // ignore
			}
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestAddressHash(t *testing.T) {
	t.Parallel()
	type args struct {
		userID       string
		postalCode   string
		addressLine1 string
		addressLine2 string
	}
	tests := []struct {
		name   string
		args   args
		expect string
	}{
		{
			name: "success",
			args: args{
				userID:       "user-id",
				postalCode:   "1000014",
				addressLine1: "永田町1-7-1",
				addressLine2: "",
			},
			expect: "789ef22a79a364f95c66a3d3b1fda213c1316a6c7f8b6306b493d8c46d2dce75",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewAddressHash(tt.args.userID, tt.args.postalCode, tt.args.addressLine1, tt.args.addressLine2)
			assert.Equal(t, tt.expect, actual)
		})
	}
}
