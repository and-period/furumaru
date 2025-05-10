package entity

import (
	"testing"

	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/stretchr/testify/assert"
)

func TestCrditCard(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		params *NewCreditCardParams
		expect *CreditCard
	}{
		{
			name: "card brand is visa",
			params: &NewCreditCardParams{
				Name:   "ＦＵＲＵＭＡＲＵ　TARO",
				Number: "4242424242424242",
				Month:  12,
				Year:   2024,
				CVV:    "123",
			},
			expect: &CreditCard{
				Name:   "FURUMARU TARO",
				Number: "4242424242424242",
				Month:  12,
				Year:   2024,
				CVV:    "123",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, NewCreditCard(tt.params))
		})
	}
}

func TestCrditCard_Validate(t *testing.T) {
	t.Parallel()
	now := jst.Date(2024, 1, 2, 18, 30, 0, 0)
	tests := []struct {
		name   string
		card   *CreditCard
		expect error
	}{
		{
			name: "validate",
			card: &CreditCard{
				Number: "4242424242424242",
				Month:  12,
				Year:   2024,
				CVV:    "123",
			},
			expect: nil,
		},
		{
			name: "invalid card number for too short",
			card: &CreditCard{
				Number: "123456789012",
				Month:  12,
				Year:   2024,
				CVV:    "123",
			},
			expect: errCreditCardInvalidNumber,
		},
		{
			name: "invalid card number for too long",
			card: &CreditCard{
				Number: "12345678901234567",
				Month:  12,
				Year:   2024,
				CVV:    "123",
			},
			expect: errCreditCardInvalidNumber,
		},
		{
			name: "invalid card number for luhn algorithm",
			card: &CreditCard{
				Number: "1234567890123456",
				Month:  12,
				Year:   2024,
				CVV:    "123",
			},
			expect: errCreditCardInvalidNumber,
		},
		{
			name: "invalid cvv for too short",
			card: &CreditCard{
				Number: "4242424242424242",
				Month:  12,
				Year:   2024,
				CVV:    "12",
			},
			expect: errCreditCardInvalidCVV,
		},
		{
			name: "invalid cvv for too long",
			card: &CreditCard{
				Number: "4242424242424242",
				Month:  12,
				Year:   2024,
				CVV:    "12345",
			},
			expect: errCreditCardInvalidCVV,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.card.Validate(now)
			assert.ErrorIs(t, err, tt.expect)
		})
	}
}
