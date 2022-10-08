package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCoordinator(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		params *NewCoordinatorParams
		expect *Coordinator
	}{
		{
			name: "success",
			params: &NewCoordinatorParams{
				Admin: &Admin{
					ID:            "admin-id",
					CognitoID:     "cognito-id",
					Lastname:      "&.",
					Firstname:     "スタッフ",
					LastnameKana:  "あんどぴりおど",
					FirstnameKana: "すたっふ",
					Email:         "test-coordinator@and-period.jp",
				},
				PhoneNumber: "+819012345678",
			},
			expect: &Coordinator{
				AdminID:     "admin-id",
				PhoneNumber: "+819012345678",
				Admin: Admin{
					CognitoID:     "cognito-id",
					Lastname:      "&.",
					Firstname:     "スタッフ",
					LastnameKana:  "あんどぴりおど",
					FirstnameKana: "すたっふ",
					Email:         "test-coordinator@and-period.jp",
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewCoordinator(tt.params)
			actual.ID = "" // ignore
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestCoordinator_Fill(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		coordinator *Coordinator
		admin       *Admin
		expect      *Coordinator
	}{
		{
			name: "success",
			coordinator: &Coordinator{
				AdminID: "admin-id",
			},
			admin: &Admin{
				ID:        "admin-id",
				CognitoID: "cognito-id",
			},
			expect: &Coordinator{
				AdminID: "admin-id",
				Admin: Admin{
					ID:        "admin-id",
					CognitoID: "cognito-id",
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.coordinator.Fill(tt.admin)
			assert.Equal(t, tt.expect, tt.coordinator)
		})
	}
}

func TestCoordinators_IDs(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		coordinators Coordinators
		expect       []string
	}{
		{
			name: "success",
			coordinators: Coordinators{
				{AdminID: "coordinator-id01"},
				{AdminID: "coordinator-id02"},
			},
			expect: []string{
				"coordinator-id01",
				"coordinator-id02",
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.coordinators.IDs())
		})
	}
}
