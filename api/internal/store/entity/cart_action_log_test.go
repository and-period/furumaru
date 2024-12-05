package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddCartItemActionLog(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		params *AddCartItemActionLogParams
		expect *CartActionLog
	}{
		{
			name: "success",
			params: &AddCartItemActionLogParams{
				SessionID: "session-id",
				UserID:    "user-id",
				UserAgent: "user-agent",
				ClientIP:  "client-ip",
				ProductID: "product-id",
			},
			expect: &CartActionLog{
				SessionID: "session-id",
				Type:      CartActionLogTypeAddCartItem,
				UserID:    "user-id",
				UserAgent: "user-agent",
				ClientIP:  "client-ip",
				ProductID: "product-id",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewAddCartItemActionLog(tt.params)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestRemoveCartItemActionLog(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		params *RemoveCartItemActionLogParams
		expect *CartActionLog
	}{
		{
			name: "success",
			params: &RemoveCartItemActionLogParams{
				SessionID: "session-id",
				UserID:    "user-id",
				UserAgent: "user-agent",
				ClientIP:  "client-ip",
				ProductID: "product-id",
			},
			expect: &CartActionLog{
				SessionID: "session-id",
				Type:      CartActionLogTypeRemoveCartItem,
				UserID:    "user-id",
				UserAgent: "user-agent",
				ClientIP:  "client-ip",
				ProductID: "product-id",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewRemoveCartItemActionLog(tt.params)
			assert.Equal(t, tt.expect, actual)
		})
	}
}
