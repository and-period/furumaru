package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrderExperiences_All(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		experiences OrderExperiences
	}{
		{
			name: "success",
			experiences: OrderExperiences{
				{OrderID: "order-01", ExperienceRevisionID: 1},
				{OrderID: "order-02", ExperienceRevisionID: 2},
			},
		},
		{
			name:        "empty",
			experiences: OrderExperiences{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var count int
			for range tt.experiences.All() {
				count++
			}
			assert.Equal(t, len(tt.experiences), count)
		})
	}
}

func TestOrderExperiences_IterMapByOrderID(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		experiences OrderExperiences
	}{
		{
			name: "success",
			experiences: OrderExperiences{
				{OrderID: "order-01", ExperienceRevisionID: 1},
				{OrderID: "order-02", ExperienceRevisionID: 2},
			},
		},
		{
			name:        "empty",
			experiences: OrderExperiences{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := make(map[string]*OrderExperience)
			for k, v := range tt.experiences.IterMapByOrderID() {
				result[k] = v
			}
			assert.Len(t, result, len(tt.experiences))
			for _, e := range tt.experiences {
				assert.Contains(t, result, e.OrderID)
			}
		})
	}
}
