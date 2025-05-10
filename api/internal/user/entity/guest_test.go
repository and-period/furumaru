package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGuests_Map(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		guests Guests
		expect map[string]*Guest
	}{
		{
			name: "success",
			guests: Guests{
				{
					UserID: "user-id01",
					Email:  "test-admin01@and-period.jp",
				},
				{
					UserID: "user-id02",
					Email:  "test-admin02@and-period.jp",
				},
			},
			expect: map[string]*Guest{
				"user-id01": {
					UserID: "user-id01",
					Email:  "test-admin01@and-period.jp",
				},
				"user-id02": {
					UserID: "user-id02",
					Email:  "test-admin02@and-period.jp",
				},
			},
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.guests.Map())
		})
	}
}
