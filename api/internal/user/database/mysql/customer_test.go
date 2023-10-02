package mysql

import (
	"time"

	"github.com/and-period/furumaru/api/internal/codes"
	"github.com/and-period/furumaru/api/internal/user/entity"
)

func testCustomer(id string, now time.Time) *entity.Customer {
	return &entity.Customer{
		UserID:        id,
		Lastname:      "&.",
		Firstname:     "スタッフ",
		LastnameKana:  "あんどぴりおど",
		FirstnameKana: "すたっふ",
		PostalCode:    "1000014",
		Prefecture:    codes.PrefectureValues["tokyo"],
		City:          "千代田区",
		AddressLine1:  "永田町1-7-1",
		AddressLine2:  "",
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}
