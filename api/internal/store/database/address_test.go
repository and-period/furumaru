package database

import (
	"time"

	"github.com/and-period/furumaru/api/internal/store/entity"
)

func testAddress(id, userID string, now time.Time) *entity.Address {
	return &entity.Address{
		ID:             id,
		UserID:         userID,
		Hash:           entity.NewAddressHash(userID, "1000014", "永田町1-7-1", ""),
		IsDefault:      false,
		Lastname:       "&.",
		Firstname:      "購入者",
		PostalCode:     "1000014",
		Prefecture:     "東京都",
		PrefectureCode: 13,
		City:           "千代田区",
		AddressLine1:   "永田町1-7-1",
		AddressLine2:   "",
		PhoneNumber:    "+819012345678",
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}
