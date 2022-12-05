package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFulfillments_MapByOrderID(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		fulfillments Fulfillments
		expect       map[string]*Fulfillment
	}{
		{
			name: "success",
			fulfillments: Fulfillments{
				{
					OrderID: "order-id01",
				},
				{
					OrderID: "order-id02",
				},
			},
			expect: map[string]*Fulfillment{
				"order-id01": {
					OrderID: "order-id01",
				},
				"order-id02": {
					OrderID: "order-id02",
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.fulfillments.MapByOrderID())
		})
	}
}
