package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestActivities_GroupByOrderID(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		activities Activities
		expect     map[string]Activities
	}{
		{
			name: "success",
			activities: Activities{
				{
					ID:      "activity-id01",
					OrderID: "order-id01",
				},
				{
					ID:      "activity-id02",
					OrderID: "order-id01",
				},
				{
					ID:      "activity-id03",
					OrderID: "order-id02",
				},
			},
			expect: map[string]Activities{
				"order-id01": {
					{
						ID:      "activity-id01",
						OrderID: "order-id01",
					},
					{
						ID:      "activity-id02",
						OrderID: "order-id01",
					},
				},
				"order-id02": {
					{
						ID:      "activity-id03",
						OrderID: "order-id02",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.activities.GroupByOrderID())
		})
	}
}
