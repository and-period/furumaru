package database

import (
	"time"

	"github.com/and-period/furumaru/api/internal/store/entity"
)

func testOrderPayment(id, transactionID, orderID, promotionID, paymentID string, now time.Time) *entity.OrderPayment {
	return &entity.OrderPayment{
		ID:             id,
		TransactionID:  transactionID,
		OrderID:        orderID,
		PromotionID:    promotionID,
		PaymentID:      paymentID,
		PaymentType:    entity.PaymentTypeCard,
		Subtotal:       100,
		Discount:       0,
		ShippingCharge: 500,
		Tax:            60,
		Total:          660,
		Lastname:       "&.",
		Firstname:      "スタッフ",
		PostalCode:     "1000014",
		Prefecture:     "東京都",
		City:           "千代田区",
		AddressLine1:   "永田町1-7-1",
		AddressLine2:   "",
		PhoneNumber:    "+819012345678",
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}
