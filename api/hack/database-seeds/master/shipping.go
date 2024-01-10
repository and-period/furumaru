package master

import "github.com/and-period/furumaru/api/internal/store/entity"

var DefaultShipping = &entity.Shipping{
	ID:            "default",
	CoordinatorID: "",
}

var DefaultShippingRevision = &entity.ShippingRevision{
	ID:         1,
	ShippingID: "default",
	Box60Rates: []*entity.ShippingRate{
		{
			Number: 1,
			Name:   "本州・四国・九州",
			Price:  1060,
			PrefectureCodes: []int32{
				2, 3, 4, 5, 6, 7, 8, 9, 10,
				11, 12, 13, 14, 15, 16, 17, 18, 19, 20,
				21, 22, 23, 24, 25, 26, 27, 28, 29, 30,
				31, 32, 33, 34, 35, 36, 37, 38, 39, 40,
				41, 42, 43, 44, 45, 46,
			},
		},
		{
			Number: 2,
			Name:   "北海道・沖縄",
			Price:  1920,
			PrefectureCodes: []int32{
				1, 47,
			},
		},
	},
	Box60Frozen: 220,
	Box80Rates: []*entity.ShippingRate{
		{
			Number: 1,
			Name:   "本州・四国・九州",
			Price:  1350,
			PrefectureCodes: []int32{
				2, 3, 4, 5, 6, 7, 8, 9, 10,
				11, 12, 13, 14, 15, 16, 17, 18, 19, 20,
				21, 22, 23, 24, 25, 26, 27, 28, 29, 30,
				31, 32, 33, 34, 35, 36, 37, 38, 39, 40,
				41, 42, 43, 44, 45, 46,
			},
		},
		{
			Number: 2,
			Name:   "北海道・沖縄",
			Price:  2200,
			PrefectureCodes: []int32{
				1, 47,
			},
		},
	},
	Box80Frozen: 220,
	Box100Rates: []*entity.ShippingRate{
		{
			Number: 1,
			Name:   "本州・四国・九州",
			Price:  1650,
			PrefectureCodes: []int32{
				2, 3, 4, 5, 6, 7, 8, 9, 10,
				11, 12, 13, 14, 15, 16, 17, 18, 19, 20,
				21, 22, 23, 24, 25, 26, 27, 28, 29, 30,
				31, 32, 33, 34, 35, 36, 37, 38, 39, 40,
				41, 42, 43, 44, 45, 46,
			},
		},
		{
			Number: 2,
			Name:   "北海道・沖縄",
			Price:  2510,
			PrefectureCodes: []int32{
				1, 47,
			},
		},
	},
	Box100Frozen:      220,
	HasFreeShipping:   false,
	FreeShippingRates: 0,
}
