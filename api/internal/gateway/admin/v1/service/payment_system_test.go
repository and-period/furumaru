package service

import (
	"testing"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/stretchr/testify/assert"
)

func TestPaymentSystemStatus(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		status entity.PaymentSystemStatus
		expect PaymentSystemStatus
	}{
		{
			name:   "in use",
			status: entity.PaymentSystemStatusInUse,
			expect: PaymentSystemStatusInUse,
		},
		{
			name:   "outage",
			status: entity.PaymentSystemStatusOutage,
			expect: PaymentSystemStatusOutage,
		},
		{
			name:   "unknown",
			status: entity.PaymentSystemStatusUnknown,
			expect: PaymentSystemStatusUnknown,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewPaymentSystemStatus(tt.status)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestPaymentSystemStatus_StoreEntity(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		status PaymentSystemStatus
		expect entity.PaymentSystemStatus
	}{
		{
			name:   "in use",
			status: PaymentSystemStatusInUse,
			expect: entity.PaymentSystemStatusInUse,
		},
		{
			name:   "outage",
			status: PaymentSystemStatusOutage,
			expect: entity.PaymentSystemStatusOutage,
		},
		{
			name:   "unknown",
			status: PaymentSystemStatusUnknown,
			expect: entity.PaymentSystemStatusUnknown,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.status.StoreEntity())
		})
	}
}

func TestPaymentSystemStatus_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		status PaymentSystemStatus
		expect int32
	}{
		{
			name:   "success",
			status: PaymentSystemStatusInUse,
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

func TestPaymentSystem(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		system *entity.PaymentSystem
		expect *PaymentSystem
	}{
		{
			name: "success",
			system: &entity.PaymentSystem{
				MethodType: entity.PaymentMethodTypeCreditCard,
				Status:     entity.PaymentSystemStatusInUse,
				CreatedAt:  jst.Date(2022, 1, 1, 0, 0, 0, 0),
				UpdatedAt:  jst.Date(2022, 1, 1, 0, 0, 0, 0),
			},
			expect: &PaymentSystem{
				PaymentSystem: response.PaymentSystem{
					MethodType: PaymentMethodTypeCreditCard.Response(),
					Status:     PaymentSystemStatusInUse.Response(),
					CreatedAt:  1640962800,
					UpdatedAt:  1640962800,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewPaymentSystem(tt.system)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestPaymentSystem_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		system *PaymentSystem
		expect *response.PaymentSystem
	}{
		{
			name: "success",
			system: &PaymentSystem{
				PaymentSystem: response.PaymentSystem{
					MethodType: PaymentMethodTypeCreditCard.Response(),
					Status:     PaymentSystemStatusInUse.Response(),
					CreatedAt:  1640962800,
					UpdatedAt:  1640962800,
				},
			},
			expect: &response.PaymentSystem{
				MethodType: PaymentMethodTypeCreditCard.Response(),
				Status:     PaymentSystemStatusInUse.Response(),
				CreatedAt:  1640962800,
				UpdatedAt:  1640962800,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.system.Response())
		})
	}
}

func TestPaymentSystems(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		systems entity.PaymentSystems
		expect  PaymentSystems
	}{
		{
			name: "success",
			systems: entity.PaymentSystems{
				{
					MethodType: entity.PaymentMethodTypeCreditCard,
					Status:     entity.PaymentSystemStatusInUse,
					CreatedAt:  jst.Date(2022, 1, 1, 0, 0, 0, 0),
					UpdatedAt:  jst.Date(2022, 1, 1, 0, 0, 0, 0),
				},
			},
			expect: PaymentSystems{
				{
					PaymentSystem: response.PaymentSystem{
						MethodType: PaymentMethodTypeCreditCard.Response(),
						Status:     PaymentSystemStatusInUse.Response(),
						CreatedAt:  1640962800,
						UpdatedAt:  1640962800,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewPaymentSystems(tt.systems)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestPaymentSystems_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		systems PaymentSystems
		expect  []*response.PaymentSystem
	}{
		{
			name: "success",
			systems: PaymentSystems{
				{
					PaymentSystem: response.PaymentSystem{
						MethodType: PaymentMethodTypeCreditCard.Response(),
						Status:     PaymentSystemStatusInUse.Response(),
						CreatedAt:  1640962800,
						UpdatedAt:  1640962800,
					},
				},
			},
			expect: []*response.PaymentSystem{
				{
					MethodType: PaymentMethodTypeCreditCard.Response(),
					Status:     PaymentSystemStatusInUse.Response(),
					CreatedAt:  1640962800,
					UpdatedAt:  1640962800,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.systems.Response())
		})
	}
}
