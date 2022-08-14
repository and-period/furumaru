package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAdministrator(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		params *NewAdministratorParams
		expect *Administrator
	}{
		{
			name: "success",
			params: &NewAdministratorParams{
				Lastname:      "&.",
				Firstname:     "スタッフ",
				LastnameKana:  "あんどぴりおど",
				FirstnameKana: "すたっふ",
				Email:         "test-admin@and-period.jp",
				PhoneNumber:   "+819012345678",
			},
			expect: &Administrator{
				Lastname:      "&.",
				Firstname:     "スタッフ",
				LastnameKana:  "あんどぴりおど",
				FirstnameKana: "すたっふ",
				Email:         "test-admin@and-period.jp",
				PhoneNumber:   "+819012345678",
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewAdministrator(tt.params)
			actual.ID = "" // ignore
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestAdministrator_Name(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		administrator *Administrator
		expect        string
	}{
		{
			name:          "success",
			administrator: &Administrator{Lastname: "&.", Firstname: "スタッフ"},
			expect:        "&. スタッフ",
		},
		{
			name:          "success only lastname",
			administrator: &Administrator{Lastname: "&.", Firstname: ""},
			expect:        "&.",
		},
		{
			name:          "success only firstname",
			administrator: &Administrator{Lastname: "", Firstname: "スタッフ"},
			expect:        "スタッフ",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.administrator.Name())
		})
	}
}

func TestAdministrators_IDs(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name           string
		administrators Administrators
		expect         []string
	}{
		{
			name: "success",
			administrators: Administrators{
				{ID: "administrator-id01"},
				{ID: "administrator-id02"},
			},
			expect: []string{
				"administrator-id01",
				"administrator-id02",
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.administrators.IDs())
		})
	}
}
