package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShippingRevisions_All(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		revisions ShippingRevisions
	}{
		{
			name: "success",
			revisions: ShippingRevisions{
				{ID: 1, ShippingID: "shipping-01"},
				{ID: 2, ShippingID: "shipping-02"},
			},
		},
		{
			name:      "empty",
			revisions: ShippingRevisions{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var count int
			for range tt.revisions.All() {
				count++
			}
			assert.Equal(t, len(tt.revisions), count)
		})
	}
}

func TestShippingRevisions_IterMapByShippingID(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		revisions ShippingRevisions
	}{
		{
			name: "success",
			revisions: ShippingRevisions{
				{ID: 1, ShippingID: "shipping-01"},
				{ID: 2, ShippingID: "shipping-02"},
			},
		},
		{
			name:      "empty",
			revisions: ShippingRevisions{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := make(map[string]*ShippingRevision)
			for k, v := range tt.revisions.IterMapByShippingID() {
				result[k] = v
			}
			assert.Len(t, result, len(tt.revisions))
			for _, r := range tt.revisions {
				assert.Contains(t, result, r.ShippingID)
			}
		})
	}
}
