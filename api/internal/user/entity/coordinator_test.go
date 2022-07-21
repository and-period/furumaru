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
				Lastname:      "&.",
				Firstname:     "スタッフ",
				LastnameKana:  "あんどぴりおど",
				FirstnameKana: "すたっふ",
				Email:         "test-coordinator@and-period.jp",
				PhoneNumber:   "+819012345678",
			},
			expect: &Coordinator{
				Lastname:      "&.",
				Firstname:     "スタッフ",
				LastnameKana:  "あんどぴりおど",
				FirstnameKana: "すたっふ",
				Email:         "test-coordinator@and-period.jp",
				PhoneNumber:   "+819012345678",
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

func TestCoordinator_Name(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		coordinator *Coordinator
		expect      string
	}{
		{
			name:        "success",
			coordinator: &Coordinator{Lastname: "&.", Firstname: "スタッフ"},
			expect:      "&. スタッフ",
		},
		{
			name:        "success only lastname",
			coordinator: &Coordinator{Lastname: "&.", Firstname: ""},
			expect:      "&.",
		},
		{
			name:        "success only firstname",
			coordinator: &Coordinator{Lastname: "", Firstname: "スタッフ"},
			expect:      "スタッフ",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.coordinator.Name())
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
				{ID: "coordinator-id01"},
				{ID: "coordinator-id02"},
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
