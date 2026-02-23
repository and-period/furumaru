package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMultiOrderMetadata_All(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		metadata MultiOrderMetadata
	}{
		{
			name: "success",
			metadata: MultiOrderMetadata{
				{OrderID: "order-01"},
				{OrderID: "order-02"},
			},
		},
		{
			name:     "empty",
			metadata: MultiOrderMetadata{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var count int
			for range tt.metadata.All() {
				count++
			}
			assert.Equal(t, len(tt.metadata), count)
		})
	}
}

func TestMultiOrderMetadata_IterMapByOrderID(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		metadata MultiOrderMetadata
	}{
		{
			name: "success",
			metadata: MultiOrderMetadata{
				{OrderID: "order-01"},
				{OrderID: "order-02"},
			},
		},
		{
			name:     "empty",
			metadata: MultiOrderMetadata{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := make(map[string]*OrderMetadata)
			for k, v := range tt.metadata.IterMapByOrderID() {
				result[k] = v
			}
			assert.Len(t, result, len(tt.metadata))
			for _, m := range tt.metadata {
				assert.Contains(t, result, m.OrderID)
			}
		})
	}
}
