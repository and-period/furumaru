package database

import (
	"time"

	"github.com/and-period/furumaru/api/internal/store/entity"
)

func testPayment(orderID, addressID, transactionID, paymentMethodID string, now time.Time) *entity.Payment {
	return &entity.Payment{
		OrderID:       orderID,
		AddressID:     addressID,
		TransactionID: transactionID,
		MethodType:    entity.PaymentMethodTypeCreditCard,
		Subtotal:      1800,
		Discount:      0,
		ShippingFee:   500,
		Tax:           230,
		Total:         2530,
		RefundTotal:   0,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}
