package stripe

import (
	"context"
	"testing"

	"github.com/and-period/furumaru/api/internal/store/payment"
	"github.com/stretchr/testify/assert"
)

func TestOrderCreditCard(t *testing.T) {
	t.Parallel()
	p := &provider{}
	result, err := p.OrderCreditCard(context.Background(), &payment.OrderCreditCardParams{
		SessionID: "pi_xxx",
	})
	assert.NoError(t, err)
	assert.Equal(t, &payment.OrderResult{}, result)
}

func TestOrderBankTransfer(t *testing.T) {
	t.Parallel()
	p := &provider{}
	_, err := p.OrderBankTransfer(context.Background(), &payment.OrderBankTransferParams{})
	assert.ErrorIs(t, err, ErrNotSupported)
}

func TestOrderKonbini(t *testing.T) {
	t.Parallel()
	p := &provider{}
	_, err := p.OrderKonbini(context.Background(), &payment.OrderKonbiniParams{})
	assert.ErrorIs(t, err, ErrNotSupported)
}

func TestOrderPayPay(t *testing.T) {
	t.Parallel()
	p := &provider{}
	_, err := p.OrderPayPay(context.Background(), &payment.OrderPayPayParams{})
	assert.ErrorIs(t, err, ErrNotSupported)
}

func TestOrderLinePay(t *testing.T) {
	t.Parallel()
	p := &provider{}
	_, err := p.OrderLinePay(context.Background(), &payment.OrderLinePayParams{})
	assert.ErrorIs(t, err, ErrNotSupported)
}

func TestOrderMerpay(t *testing.T) {
	t.Parallel()
	p := &provider{}
	_, err := p.OrderMerpay(context.Background(), &payment.OrderMerpayParams{})
	assert.ErrorIs(t, err, ErrNotSupported)
}

func TestOrderRakutenPay(t *testing.T) {
	t.Parallel()
	p := &provider{}
	_, err := p.OrderRakutenPay(context.Background(), &payment.OrderRakutenPayParams{})
	assert.ErrorIs(t, err, ErrNotSupported)
}

func TestOrderAUPay(t *testing.T) {
	t.Parallel()
	p := &provider{}
	_, err := p.OrderAUPay(context.Background(), &payment.OrderAUPayParams{})
	assert.ErrorIs(t, err, ErrNotSupported)
}

func TestOrderPaidy(t *testing.T) {
	t.Parallel()
	p := &provider{}
	_, err := p.OrderPaidy(context.Background(), &payment.OrderPaidyParams{})
	assert.ErrorIs(t, err, ErrNotSupported)
}

func TestOrderPayEasy(t *testing.T) {
	t.Parallel()
	p := &provider{}
	_, err := p.OrderPayEasy(context.Background(), &payment.OrderPayEasyParams{})
	assert.ErrorIs(t, err, ErrNotSupported)
}
