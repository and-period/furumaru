package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCustomer_Fill(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		customer *Customer
		expect   *Customer
	}{
		{
			name: "success",
			customer: &Customer{
				PrefectureCode: 13,
			},
			expect: &Customer{
				Prefecture:     "東京都",
				PrefectureCode: 13,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.customer.Fill()
			assert.Equal(t, tt.expect, tt.customer)
		})
	}
}

func TestCustomers_Map(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		customers Customers
		expect    map[string]*Customer
	}{
		{
			name: "success",
			customers: Customers{
				{
					UserID:    "user-id01",
					Lastname:  "&.",
					Firstname: "スタッフ",
				},
				{
					UserID:    "user-id02",
					Lastname:  "&.",
					Firstname: "スタッフ",
				},
			},
			expect: map[string]*Customer{
				"user-id01": {
					UserID:    "user-id01",
					Lastname:  "&.",
					Firstname: "スタッフ",
				},
				"user-id02": {
					UserID:    "user-id02",
					Lastname:  "&.",
					Firstname: "スタッフ",
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.customers.Map())
		})
	}
}

func TestCustomers_Fill(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		customers Customers
		expect    Customers
	}{
		{
			name: "success",
			customers: Customers{
				{
					PrefectureCode: 13,
				},
			},
			expect: Customers{
				{
					Prefecture:     "東京都",
					PrefectureCode: 13,
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.customers.Fill()
			assert.Equal(t, tt.expect, tt.customers)
		})
	}
}
