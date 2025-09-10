package service

import (
	"testing"

	"github.com/and-period/furumaru/api/internal/gateway/user/facility/types"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/stretchr/testify/assert"
)

func TestShippingSize(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		status entity.ShippingSize
		expect types.ShippingSize
	}{
		{
			name:   "size 60",
			status: entity.ShippingSize60,
			expect: types.ShippingSize60,
		},
		{
			name:   "size 80",
			status: entity.ShippingSize80,
			expect: types.ShippingSize80,
		},
		{
			name:   "size 100",
			status: entity.ShippingSize100,
			expect: types.ShippingSize100,
		},
		{
			name:   "unknown",
			status: entity.ShippingSizeUnknown,
			expect: types.ShippingSizeUnknown,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, NewShippingSize(tt.status))
		})
	}
}

func TestShippingSize_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		status ShippingSize
		expect int32
	}{
		{
			name:   "success",
			status: ShippingSize(types.ShippingSize60),
			expect: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.status.Response())
		})
	}
}

func TestShippingType(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		status entity.ShippingType
		expect types.ShippingType
	}{
		{
			name:   "normal",
			status: entity.ShippingTypeNormal,
			expect: types.ShippingTypeNormal,
		},
		{
			name:   "frozen",
			status: entity.ShippingTypeFrozen,
			expect: types.ShippingTypeFrozen,
		},
		{
			name:   "unknown",
			status: entity.ShippingTypeUnknown,
			expect: types.ShippingTypeUnknown,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, NewShippingType(tt.status))
		})
	}
}

func TestShippingType_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		status ShippingType
		expect int32
	}{
		{
			name:   "success",
			status: ShippingType(types.ShippingTypeNormal),
			expect: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.status.Response())
		})
	}
}
