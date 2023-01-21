package handler

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/and-period/furumaru/api/internal/common"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/request"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/service"
	"github.com/and-period/furumaru/api/internal/store"
	sentity "github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/internal/user"
	uentity "github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGetSchedule(t *testing.T) {
	t.Parallel()

	producersIn := &user.MultiGetProducersInput{
		ProducerIDs: []string{"producer-id"},
	}
	producers := uentity.Producers{
		{
			Admin: uentity.Admin{
				ID:            "producer-id",
				Lastname:      "&.",
				Firstname:     "管理者",
				LastnameKana:  "あんどどっと",
				FirstnameKana: "かんりしゃ",
				Email:         "test-producer@and-period.jp",
			},
			AdminID:       "producer-id",
			CoordinatorID: "coordinator-id",
			StoreName:     "&.農園",
			ThumbnailURL:  "https://and-period.jp/thumbnail.png",
			HeaderURL:     "https://and-period.jp/header.png",
			PhoneNumber:   "+819012345678",
			PostalCode:    "1000014",
			Prefecture:    "東京都",
			City:          "千代田区",
			CreatedAt:     jst.Date(2022, 1, 1, 0, 0, 0, 0),
			UpdatedAt:     jst.Date(2022, 1, 1, 0, 0, 0, 0),
		},
	}

	shippingIn := &store.GetShippingInput{
		ShippingID: "shipping-id",
	}
	shipping := &sentity.Shipping{
		ID:   "shipping-id",
		Name: "デフォルト配送設定",
		Box60Rates: sentity.ShippingRates{
			{Number: 1, Name: "東京都", Price: 0, Prefectures: []int64{13}},
		},
		Box60Refrigerated:  500,
		Box60Frozen:        800,
		Box80Rates:         sentity.ShippingRates{},
		Box80Refrigerated:  500,
		Box80Frozen:        800,
		Box100Rates:        sentity.ShippingRates{},
		Box100Refrigerated: 500,
		Box100Frozen:       800,
		HasFreeShipping:    true,
		FreeShippingRates:  3000,
		CreatedAt:          jst.Date(2022, 1, 1, 0, 0, 0, 0),
		UpdatedAt:          jst.Date(2022, 1, 1, 0, 0, 0, 0),
	}
	categoriesIn := &store.MultiGetCategoriesInput{
		CategoryIDs: []string{"category-id"},
	}
	categories := sentity.Categories{
		{
			ID:        "category-id",
			Name:      "野菜",
			CreatedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
			UpdatedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
		},
	}
	productTypesIn := &store.MultiGetProductTypesInput{
		ProductTypeIDs: []string{"product-type-id"},
	}
	productTypes := sentity.ProductTypes{
		{
			ID:         "product-type-id",
			Name:       "じゃがいも",
			CategoryID: "category-id",
			IconURL:    "https://and-period.jp/icon.png",
			CreatedAt:  jst.Date(2022, 1, 1, 0, 0, 0, 0),
			UpdatedAt:  jst.Date(2022, 1, 1, 0, 0, 0, 0),
		},
	}
	productsIn := &store.MultiGetProductsInput{
		ProductIDs: []string{"product-id"},
	}
	product := &sentity.Product{
		ID:              "product-id",
		TypeID:          "product-type-id",
		ProducerID:      "producer-id",
		Name:            "新鮮なじゃがいも",
		Description:     "新鮮なじゃがいもをお届けします。",
		Public:          true,
		Inventory:       100,
		Weight:          1300,
		WeightUnit:      sentity.WeightUnitGram,
		Item:            1,
		ItemUnit:        "袋",
		ItemDescription: "1袋あたり100gのじゃがいも",
		Media: sentity.MultiProductMedia{
			{
				URL:         "https://and-period.jp/thumbnail01.png",
				IsThumbnail: true,
				Images: common.Images{
					{URL: "https://and-period.jp/thumbnail01_240.png", Size: common.ImageSizeSmall},
					{URL: "https://and-period.jp/thumbnail01_675.png", Size: common.ImageSizeMedium},
					{URL: "https://and-period.jp/thumbnail01_900.png", Size: common.ImageSizeLarge},
				},
			},
			{
				URL:         "https://and-period.jp/thumbnail02.png",
				IsThumbnail: false,
				Images: common.Images{
					{URL: "https://and-period.jp/thumbnail02_240.png", Size: common.ImageSizeSmall},
					{URL: "https://and-period.jp/thumbnail02_675.png", Size: common.ImageSizeMedium},
					{URL: "https://and-period.jp/thumbnail02_900.png", Size: common.ImageSizeLarge},
				},
			},
		},
		Price:            400,
		DeliveryType:     sentity.DeliveryTypeNormal,
		Box60Rate:        50,
		Box80Rate:        40,
		Box100Rate:       30,
		OriginPrefecture: "滋賀県",
		OriginCity:       "彦根市",
		CreatedAt:        jst.Date(2022, 1, 1, 0, 0, 0, 0),
		UpdatedAt:        jst.Date(2022, 1, 1, 0, 0, 0, 0),
	}
	products := sentity.Products{product}

	scheduleIn := &store.GetScheduleInput{
		ScheduleID: "schedule-id",
	}
	livesIn := &store.ListLivesByScheduleIDInput{
		ScheduleID: "schedule-id",
	}
	schedule := &sentity.Schedule{
		ID:            "schedule-id",
		CoordinatorID: "coordinator-id",
		ShippingID:    "shipping-id",
		Title:         "スケジュールタイトル",
		Description:   "スケジュールの説明",
		ThumbnailURL:  "https://and-period.jp/thumbnail01.png",
		StartAt:       jst.Date(2022, 1, 1, 0, 0, 0, 0),
		EndAt:         jst.Date(2022, 1, 1, 0, 0, 0, 0),
		CreatedAt:     jst.Date(2022, 1, 1, 0, 0, 0, 0),
		UpdatedAt:     jst.Date(2022, 1, 1, 0, 0, 0, 0),
	}
	live := &sentity.Live{
		LiveProducts: sentity.LiveProducts{
			{
				LiveID:    "live-id",
				ProductID: "product-id",
				CreatedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
				UpdatedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
			},
		},
		ID:          "live-id",
		ScheduleID:  "schedule-id",
		Title:       "配信タイトル",
		Description: "配信の説明",
		Status:      sentity.LiveStatusWaiting,
		Published:   false,
		Canceled:    false,
		ProducerID:  "producer-id",
		StartAt:     jst.Date(2022, 1, 1, 0, 0, 0, 0),
		EndAt:       jst.Date(2022, 1, 1, 0, 0, 0, 0),
		CreatedAt:   jst.Date(2022, 1, 1, 0, 0, 0, 0),
		UpdatedAt:   jst.Date(2022, 1, 1, 0, 0, 0, 0),
	}
	lives := sentity.Lives{live}
	tests := []struct {
		name       string
		setup      func(t *testing.T, mocks *mocks, ctrl *gomock.Controller)
		scheduleID string
		expect     *testResponse
	}{
		{
			name: "success",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.store.EXPECT().GetSchedule(gomock.Any(), scheduleIn).Return(schedule, nil)
				mocks.user.EXPECT().MultiGetProducers(gomock.Any(), producersIn).Return(producers, nil).Times(2)
				mocks.store.EXPECT().GetShipping(gomock.Any(), shippingIn).Return(shipping, nil)
				mocks.store.EXPECT().ListLivesByScheduleID(gomock.Any(), livesIn).Return(lives, nil)
				mocks.store.EXPECT().MultiGetCategories(gomock.Any(), categoriesIn).Return(categories, nil)
				mocks.store.EXPECT().MultiGetProductTypes(gomock.Any(), productTypesIn).Return(productTypes, nil)
				mocks.store.EXPECT().MultiGetProducts(gomock.Any(), productsIn).Return(products, nil)
			},
			scheduleID: "schedule-id",
			expect: &testResponse{
				code: http.StatusOK,
				body: &response.ScheduleResponse{
					Schedule: &response.Schedule{
						ID:            "schedule-id",
						CoordinatorID: "coordinator-id",
						ShippingID:    "shipping-id",
						ShippingName:  "デフォルト配送設定",
						Title:         "スケジュールタイトル",
						Description:   "スケジュールの説明",
						ThumbnailURL:  "https://and-period.jp/thumbnail01.png",
						StartAt:       1640962800,
						EndAt:         1640962800,
						CreatedAt:     1640962800,
						UpdatedAt:     1640962800,
					},
					Lives: []*response.Live{
						{
							ID:           "live-id",
							ScheduleID:   "schedule-id",
							Title:        "配信タイトル",
							Description:  "配信の説明",
							Status:       1,
							Published:    false,
							Canceled:     false,
							ProducerID:   "producer-id",
							ProducerName: "&. 管理者",
							Products: []*response.Product{
								{
									ID:              "product-id",
									TypeID:          "product-type-id",
									TypeName:        "じゃがいも",
									TypeIconURL:     "https://and-period.jp/icon.png",
									CategoryID:      "category-id",
									CategoryName:    "野菜",
									ProducerID:      "producer-id",
									StoreName:       "&.農園",
									Name:            "新鮮なじゃがいも",
									Description:     "新鮮なじゃがいもをお届けします。",
									Public:          true,
									Inventory:       100,
									Weight:          1.3,
									ItemUnit:        "袋",
									ItemDescription: "1袋あたり100gのじゃがいも",
									Media: []*response.ProductMedia{
										{
											URL:         "https://and-period.jp/thumbnail01.png",
											IsThumbnail: true,
											Images: []*response.Image{
												{URL: "https://and-period.jp/thumbnail01_240.png", Size: int32(service.ImageSizeSmall)},
												{URL: "https://and-period.jp/thumbnail01_675.png", Size: int32(service.ImageSizeMedium)},
												{URL: "https://and-period.jp/thumbnail01_900.png", Size: int32(service.ImageSizeLarge)},
											},
										},
										{
											URL:         "https://and-period.jp/thumbnail02.png",
											IsThumbnail: false,
											Images: []*response.Image{
												{URL: "https://and-period.jp/thumbnail02_240.png", Size: int32(service.ImageSizeSmall)},
												{URL: "https://and-period.jp/thumbnail02_675.png", Size: int32(service.ImageSizeMedium)},
												{URL: "https://and-period.jp/thumbnail02_900.png", Size: int32(service.ImageSizeLarge)},
											},
										},
									},
									Price:            400,
									DeliveryType:     1,
									Box60Rate:        50,
									Box80Rate:        40,
									Box100Rate:       30,
									OriginPrefecture: "滋賀県",
									OriginCity:       "彦根市",
									CreatedAt:        1640962800,
									UpdatedAt:        1640962800,
								},
							},
							StartAt:   1640962800,
							EndAt:     1640962800,
							CreatedAt: 1640962800,
							UpdatedAt: 1640962800,
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			const format = "/v1/schedules/%s"
			path := fmt.Sprintf(format, tt.scheduleID)
			testGet(t, tt.setup, tt.expect, path)
		})
	}
}

func TestCreateSchedule(t *testing.T) {
	t.Parallel()

	coordinatorIn := &user.GetCoordinatorInput{
		CoordinatorID: "coordinator-id",
	}
	productsIn := &store.MultiGetProductsInput{
		ProductIDs: []string{"product-id"},
	}
	scheduleIn := &store.CreateScheduleInput{
		CoordinatorID: "coordinator-id",
		ShippingID:    "shipping-id",
		Title:         "スケジュールタイトル",
		Description:   "スケジュールの説明",
		ThumbnailURL:  "https://and-period.jp/thumbnail01.png",
		StartAt:       jst.Date(2022, 1, 1, 0, 0, 0, 0),
		EndAt:         jst.Date(2022, 1, 1, 0, 0, 0, 0),
		Lives: []*store.CreateScheduleLive{
			{
				Title:       "配信タイトル",
				Description: "配信の説明",
				ProducerID:  "producer-id",
				ProductIDs:  []string{"product-id"},
				StartAt:     jst.Date(2022, 1, 1, 0, 0, 0, 0),
				EndAt:       jst.Date(2022, 1, 1, 0, 0, 0, 0),
			},
		},
	}
	coordinator := &uentity.Coordinator{
		Admin: uentity.Admin{
			ID:            "coordinator-id",
			Lastname:      "&.",
			Firstname:     "管理者",
			LastnameKana:  "あんどどっと",
			FirstnameKana: "かんりしゃ",
			Email:         "test-coordinator@and-period.jp",
		},
		AdminID:          "coordinator-id",
		CompanyName:      "&.株式会社",
		StoreName:        "&.農園",
		ThumbnailURL:     "https://and-period.jp/thumbnail.png",
		HeaderURL:        "https://and-period.jp/header.png",
		TwitterAccount:   "twitter-id",
		InstagramAccount: "instagram-id",
		FacebookAccount:  "facebook-id",
		PhoneNumber:      "+819012345678",
		PostalCode:       "1000014",
		Prefecture:       "東京都",
		City:             "千代田区",
		CreatedAt:        jst.Date(2022, 1, 1, 0, 0, 0, 0),
		UpdatedAt:        jst.Date(2022, 1, 1, 0, 0, 0, 0),
	}
	producersIn := &user.MultiGetProducersInput{
		ProducerIDs: []string{"producer-id"},
	}
	producers := uentity.Producers{
		{
			Admin: uentity.Admin{
				ID:            "producer-id",
				Lastname:      "&.",
				Firstname:     "管理者",
				LastnameKana:  "あんどどっと",
				FirstnameKana: "かんりしゃ",
				Email:         "test-producer@and-period.jp",
			},
			AdminID:       "producer-id",
			CoordinatorID: "coordinator-id",
			StoreName:     "&.農園",
			ThumbnailURL:  "https://and-period.jp/thumbnail.png",
			HeaderURL:     "https://and-period.jp/header.png",
			PhoneNumber:   "+819012345678",
			PostalCode:    "1000014",
			Prefecture:    "東京都",
			City:          "千代田区",
			CreatedAt:     jst.Date(2022, 1, 1, 0, 0, 0, 0),
			UpdatedAt:     jst.Date(2022, 1, 1, 0, 0, 0, 0),
		},
	}
	shippingIn := &store.GetShippingInput{
		ShippingID: "shipping-id",
	}
	shipping := &sentity.Shipping{
		ID:   "shipping-id",
		Name: "デフォルト配送設定",
		Box60Rates: sentity.ShippingRates{
			{Number: 1, Name: "東京都", Price: 0, Prefectures: []int64{13}},
		},
		Box60Refrigerated:  500,
		Box60Frozen:        800,
		Box80Rates:         sentity.ShippingRates{},
		Box80Refrigerated:  500,
		Box80Frozen:        800,
		Box100Rates:        sentity.ShippingRates{},
		Box100Refrigerated: 500,
		Box100Frozen:       800,
		HasFreeShipping:    true,
		FreeShippingRates:  3000,
		CreatedAt:          jst.Date(2022, 1, 1, 0, 0, 0, 0),
		UpdatedAt:          jst.Date(2022, 1, 1, 0, 0, 0, 0),
	}
	categoriesIn := &store.MultiGetCategoriesInput{
		CategoryIDs: []string{"category-id"},
	}
	categories := sentity.Categories{
		{
			ID:        "category-id",
			Name:      "野菜",
			CreatedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
			UpdatedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
		},
	}
	productTypesIn := &store.MultiGetProductTypesInput{
		ProductTypeIDs: []string{"product-type-id"},
	}
	productTypes := sentity.ProductTypes{
		{
			ID:         "product-type-id",
			Name:       "じゃがいも",
			CategoryID: "category-id",
			IconURL:    "https://and-period.jp/icon.png",
			CreatedAt:  jst.Date(2022, 1, 1, 0, 0, 0, 0),
			UpdatedAt:  jst.Date(2022, 1, 1, 0, 0, 0, 0),
		},
	}
	product := &sentity.Product{
		ID:              "product-id",
		TypeID:          "product-type-id",
		ProducerID:      "producer-id",
		Name:            "新鮮なじゃがいも",
		Description:     "新鮮なじゃがいもをお届けします。",
		Public:          true,
		Inventory:       100,
		Weight:          1300,
		WeightUnit:      sentity.WeightUnitGram,
		Item:            1,
		ItemUnit:        "袋",
		ItemDescription: "1袋あたり100gのじゃがいも",
		Media: sentity.MultiProductMedia{
			{
				URL:         "https://and-period.jp/thumbnail01.png",
				IsThumbnail: true,
				Images: common.Images{
					{URL: "https://and-period.jp/thumbnail01_240.png", Size: common.ImageSizeSmall},
					{URL: "https://and-period.jp/thumbnail01_675.png", Size: common.ImageSizeMedium},
					{URL: "https://and-period.jp/thumbnail01_900.png", Size: common.ImageSizeLarge},
				},
			},
			{
				URL:         "https://and-period.jp/thumbnail02.png",
				IsThumbnail: false,
				Images: common.Images{
					{URL: "https://and-period.jp/thumbnail02_240.png", Size: common.ImageSizeSmall},
					{URL: "https://and-period.jp/thumbnail02_675.png", Size: common.ImageSizeMedium},
					{URL: "https://and-period.jp/thumbnail02_900.png", Size: common.ImageSizeLarge},
				},
			},
		},
		Price:            400,
		DeliveryType:     sentity.DeliveryTypeNormal,
		Box60Rate:        50,
		Box80Rate:        40,
		Box100Rate:       30,
		OriginPrefecture: "滋賀県",
		OriginCity:       "彦根市",
		CreatedAt:        jst.Date(2022, 1, 1, 0, 0, 0, 0),
		UpdatedAt:        jst.Date(2022, 1, 1, 0, 0, 0, 0),
	}
	products := sentity.Products{product}
	schedule := &sentity.Schedule{
		ID:            "schedule-id",
		CoordinatorID: "coordinator-id",
		ShippingID:    "shipping-id",
		Title:         "スケジュールタイトル",
		Description:   "スケジュールの説明",
		ThumbnailURL:  "https://and-period.jp/thumbnail01.png",
		StartAt:       jst.Date(2022, 1, 1, 0, 0, 0, 0),
		EndAt:         jst.Date(2022, 1, 1, 0, 0, 0, 0),
		CreatedAt:     jst.Date(2022, 1, 1, 0, 0, 0, 0),
		UpdatedAt:     jst.Date(2022, 1, 1, 0, 0, 0, 0),
	}
	live := &sentity.Live{
		LiveProducts: sentity.LiveProducts{
			{
				LiveID:    "live-id",
				ProductID: "product-id",
				CreatedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
				UpdatedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
			},
		},
		ID:          "live-id",
		ScheduleID:  "schedule-id",
		Title:       "配信タイトル",
		Description: "配信の説明",
		Status:      sentity.LiveStatusWaiting,
		Published:   false,
		Canceled:    false,
		ProducerID:  "producer-id",
		StartAt:     jst.Date(2022, 1, 1, 0, 0, 0, 0),
		EndAt:       jst.Date(2022, 1, 1, 0, 0, 0, 0),
		CreatedAt:   jst.Date(2022, 1, 1, 0, 0, 0, 0),
		UpdatedAt:   jst.Date(2022, 1, 1, 0, 0, 0, 0),
	}
	lives := sentity.Lives{live}
	tests := []struct {
		name    string
		setup   func(t *testing.T, mocks *mocks, ctrl *gomock.Controller)
		options []testOption
		req     *request.CreateScheduleRequest
		expect  *testResponse
	}{
		{
			name: "administrator success",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.user.EXPECT().GetCoordinator(gomock.Any(), coordinatorIn).Return(coordinator, nil)
				mocks.user.EXPECT().MultiGetProducers(gomock.Any(), producersIn).Return(producers, nil).Times(2)
				mocks.store.EXPECT().GetShipping(gomock.Any(), shippingIn).Return(shipping, nil)
				mocks.store.EXPECT().MultiGetCategories(gomock.Any(), categoriesIn).Return(categories, nil)
				mocks.store.EXPECT().MultiGetProductTypes(gomock.Any(), productTypesIn).Return(productTypes, nil)
				mocks.store.EXPECT().MultiGetProducts(gomock.Any(), productsIn).Return(products, nil)
				mocks.store.EXPECT().CreateSchedule(gomock.Any(), scheduleIn).Return(schedule, lives, nil)
			},
			options: []testOption{withRole(uentity.AdminRoleAdministrator)},
			req: &request.CreateScheduleRequest{
				CoordinatorID: "coordinator-id",
				ShippingID:    "shipping-id",
				Title:         "スケジュールタイトル",
				Description:   "スケジュールの説明",
				ThumbnailURL:  "https://and-period.jp/thumbnail01.png",
				StartAt:       1640962800,
				EndAt:         1640962800,
				Lives: []*request.CreateScheduleLive{
					{
						Title:       "配信タイトル",
						Description: "配信の説明",
						ProducerID:  "producer-id",
						ProductIDs:  []string{"product-id"},
						StartAt:     1640962800,
						EndAt:       1640962800,
					},
				},
			},
			expect: &testResponse{
				code: http.StatusOK,
				body: &response.ScheduleResponse{
					Schedule: &response.Schedule{
						ID:            "schedule-id",
						CoordinatorID: "coordinator-id",
						ShippingID:    "shipping-id",
						ShippingName:  "デフォルト配送設定",
						Title:         "スケジュールタイトル",
						Description:   "スケジュールの説明",
						ThumbnailURL:  "https://and-period.jp/thumbnail01.png",
						StartAt:       1640962800,
						EndAt:         1640962800,
						CreatedAt:     1640962800,
						UpdatedAt:     1640962800,
					},
					Lives: []*response.Live{
						{
							ID:           "live-id",
							ScheduleID:   "schedule-id",
							Title:        "配信タイトル",
							Description:  "配信の説明",
							Status:       1,
							Published:    false,
							Canceled:     false,
							ProducerID:   "producer-id",
							ProducerName: "&. 管理者",
							Products: []*response.Product{
								{
									ID:              "product-id",
									TypeID:          "product-type-id",
									TypeName:        "じゃがいも",
									TypeIconURL:     "https://and-period.jp/icon.png",
									CategoryID:      "category-id",
									CategoryName:    "野菜",
									ProducerID:      "producer-id",
									StoreName:       "&.農園",
									Name:            "新鮮なじゃがいも",
									Description:     "新鮮なじゃがいもをお届けします。",
									Public:          true,
									Inventory:       100,
									Weight:          1.3,
									ItemUnit:        "袋",
									ItemDescription: "1袋あたり100gのじゃがいも",
									Media: []*response.ProductMedia{
										{
											URL:         "https://and-period.jp/thumbnail01.png",
											IsThumbnail: true,
											Images: []*response.Image{
												{URL: "https://and-period.jp/thumbnail01_240.png", Size: int32(service.ImageSizeSmall)},
												{URL: "https://and-period.jp/thumbnail01_675.png", Size: int32(service.ImageSizeMedium)},
												{URL: "https://and-period.jp/thumbnail01_900.png", Size: int32(service.ImageSizeLarge)},
											},
										},
										{
											URL:         "https://and-period.jp/thumbnail02.png",
											IsThumbnail: false,
											Images: []*response.Image{
												{URL: "https://and-period.jp/thumbnail02_240.png", Size: int32(service.ImageSizeSmall)},
												{URL: "https://and-period.jp/thumbnail02_675.png", Size: int32(service.ImageSizeMedium)},
												{URL: "https://and-period.jp/thumbnail02_900.png", Size: int32(service.ImageSizeLarge)},
											},
										},
									},
									Price:            400,
									DeliveryType:     1,
									Box60Rate:        50,
									Box80Rate:        40,
									Box100Rate:       30,
									OriginPrefecture: "滋賀県",
									OriginCity:       "彦根市",
									CreatedAt:        1640962800,
									UpdatedAt:        1640962800,
								},
							},
							StartAt:   1640962800,
							EndAt:     1640962800,
							CreatedAt: 1640962800,
							UpdatedAt: 1640962800,
						},
					},
				},
			},
		},
		{
			name: "coordinator success",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.user.EXPECT().GetCoordinator(gomock.Any(), coordinatorIn).Return(coordinator, nil)
				mocks.user.EXPECT().MultiGetProducers(gomock.Any(), producersIn).Return(producers, nil).Times(2)
				mocks.store.EXPECT().GetShipping(gomock.Any(), shippingIn).Return(shipping, nil)
				mocks.store.EXPECT().MultiGetCategories(gomock.Any(), categoriesIn).Return(categories, nil)
				mocks.store.EXPECT().MultiGetProductTypes(gomock.Any(), productTypesIn).Return(productTypes, nil)
				mocks.store.EXPECT().MultiGetProducts(gomock.Any(), productsIn).Return(products, nil)
				mocks.store.EXPECT().CreateSchedule(gomock.Any(), scheduleIn).Return(schedule, lives, nil)
			},
			options: []testOption{withRole(uentity.AdminRoleCoordinator), withAdminID("coordinator-id")},
			req: &request.CreateScheduleRequest{
				CoordinatorID: "coordinator-id",
				ShippingID:    "shipping-id",
				Title:         "スケジュールタイトル",
				Description:   "スケジュールの説明",
				ThumbnailURL:  "https://and-period.jp/thumbnail01.png",
				StartAt:       1640962800,
				EndAt:         1640962800,
				Lives: []*request.CreateScheduleLive{
					{
						Title:       "配信タイトル",
						Description: "配信の説明",
						ProducerID:  "producer-id",
						ProductIDs:  []string{"product-id"},
						StartAt:     1640962800,
						EndAt:       1640962800,
					},
				},
			},
			expect: &testResponse{
				code: http.StatusOK,
				body: &response.ScheduleResponse{
					Schedule: &response.Schedule{
						ID:            "schedule-id",
						CoordinatorID: "coordinator-id",
						ShippingID:    "shipping-id",
						ShippingName:  "デフォルト配送設定",
						Title:         "スケジュールタイトル",
						Description:   "スケジュールの説明",
						ThumbnailURL:  "https://and-period.jp/thumbnail01.png",
						StartAt:       1640962800,
						EndAt:         1640962800,
						CreatedAt:     1640962800,
						UpdatedAt:     1640962800,
					},
					Lives: []*response.Live{
						{
							ID:           "live-id",
							ScheduleID:   "schedule-id",
							Title:        "配信タイトル",
							Description:  "配信の説明",
							Status:       1,
							Published:    false,
							Canceled:     false,
							ProducerID:   "producer-id",
							ProducerName: "&. 管理者",
							Products: []*response.Product{
								{
									ID:              "product-id",
									TypeID:          "product-type-id",
									TypeName:        "じゃがいも",
									TypeIconURL:     "https://and-period.jp/icon.png",
									CategoryID:      "category-id",
									CategoryName:    "野菜",
									ProducerID:      "producer-id",
									StoreName:       "&.農園",
									Name:            "新鮮なじゃがいも",
									Description:     "新鮮なじゃがいもをお届けします。",
									Public:          true,
									Inventory:       100,
									Weight:          1.3,
									ItemUnit:        "袋",
									ItemDescription: "1袋あたり100gのじゃがいも",
									Media: []*response.ProductMedia{
										{
											URL:         "https://and-period.jp/thumbnail01.png",
											IsThumbnail: true,
											Images: []*response.Image{
												{URL: "https://and-period.jp/thumbnail01_240.png", Size: int32(service.ImageSizeSmall)},
												{URL: "https://and-period.jp/thumbnail01_675.png", Size: int32(service.ImageSizeMedium)},
												{URL: "https://and-period.jp/thumbnail01_900.png", Size: int32(service.ImageSizeLarge)},
											},
										},
										{
											URL:         "https://and-period.jp/thumbnail02.png",
											IsThumbnail: false,
											Images: []*response.Image{
												{URL: "https://and-period.jp/thumbnail02_240.png", Size: int32(service.ImageSizeSmall)},
												{URL: "https://and-period.jp/thumbnail02_675.png", Size: int32(service.ImageSizeMedium)},
												{URL: "https://and-period.jp/thumbnail02_900.png", Size: int32(service.ImageSizeLarge)},
											},
										},
									},
									Price:            400,
									DeliveryType:     1,
									Box60Rate:        50,
									Box80Rate:        40,
									Box100Rate:       30,
									OriginPrefecture: "滋賀県",
									OriginCity:       "彦根市",
									CreatedAt:        1640962800,
									UpdatedAt:        1640962800,
								},
							},
							StartAt:   1640962800,
							EndAt:     1640962800,
							CreatedAt: 1640962800,
							UpdatedAt: 1640962800,
						},
					},
				},
			},
		},
		{
			name: "failed to get coordinator",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.user.EXPECT().GetCoordinator(gomock.Any(), coordinatorIn).Return(nil, assert.AnError)
				mocks.user.EXPECT().MultiGetProducers(gomock.Any(), producersIn).Return(producers, nil).Times(2)
				mocks.store.EXPECT().GetShipping(gomock.Any(), shippingIn).Return(shipping, nil)
				mocks.store.EXPECT().MultiGetCategories(gomock.Any(), categoriesIn).Return(categories, nil)
				mocks.store.EXPECT().MultiGetProductTypes(gomock.Any(), productTypesIn).Return(productTypes, nil)
				mocks.store.EXPECT().MultiGetProducts(gomock.Any(), productsIn).Return(products, nil)
			},
			options: []testOption{withRole(uentity.AdminRoleAdministrator)},
			req: &request.CreateScheduleRequest{
				CoordinatorID: "coordinator-id",
				ShippingID:    "shipping-id",
				Title:         "スケジュールタイトル",
				Description:   "スケジュールの説明",
				ThumbnailURL:  "https://and-period.jp/thumbnail01.png",
				StartAt:       1640962800,
				EndAt:         1640962800,
				Lives: []*request.CreateScheduleLive{
					{
						Title:       "配信タイトル",
						Description: "配信の説明",
						ProducerID:  "producer-id",
						ProductIDs:  []string{"product-id"},
						StartAt:     1640962800,
						EndAt:       1640962800,
					},
				},
			},
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
		{
			name: "failed to multi get producers",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.user.EXPECT().GetCoordinator(gomock.Any(), coordinatorIn).Return(coordinator, nil)
				mocks.user.EXPECT().MultiGetProducers(gomock.Any(), producersIn).Return(nil, assert.AnError).Times(2)
				mocks.store.EXPECT().GetShipping(gomock.Any(), shippingIn).Return(shipping, nil)
				mocks.store.EXPECT().MultiGetCategories(gomock.Any(), categoriesIn).Return(categories, nil)
				mocks.store.EXPECT().MultiGetProductTypes(gomock.Any(), productTypesIn).Return(productTypes, nil)
				mocks.store.EXPECT().MultiGetProducts(gomock.Any(), productsIn).Return(products, nil)
			},
			options: []testOption{withRole(uentity.AdminRoleAdministrator)},
			req: &request.CreateScheduleRequest{
				CoordinatorID: "coordinator-id",
				ShippingID:    "shipping-id",
				Title:         "スケジュールタイトル",
				Description:   "スケジュールの説明",
				ThumbnailURL:  "https://and-period.jp/thumbnail01.png",
				StartAt:       1640962800,
				EndAt:         1640962800,
				Lives: []*request.CreateScheduleLive{
					{
						Title:       "配信タイトル",
						Description: "配信の説明",
						ProducerID:  "producer-id",
						ProductIDs:  []string{"product-id"},
						StartAt:     1640962800,
						EndAt:       1640962800,
					},
				},
			},
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
		{
			name: "failed to contain invalid producer id",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.user.EXPECT().GetCoordinator(gomock.Any(), coordinatorIn).Return(coordinator, nil)
				mocks.user.EXPECT().MultiGetProducers(gomock.Any(), producersIn).Return(uentity.Producers{}, nil).Times(2)
				mocks.store.EXPECT().GetShipping(gomock.Any(), shippingIn).Return(shipping, nil)
				mocks.store.EXPECT().MultiGetCategories(gomock.Any(), categoriesIn).Return(categories, nil)
				mocks.store.EXPECT().MultiGetProductTypes(gomock.Any(), productTypesIn).Return(productTypes, nil)
				mocks.store.EXPECT().MultiGetProducts(gomock.Any(), productsIn).Return(products, nil)
			},
			options: []testOption{withRole(uentity.AdminRoleAdministrator)},
			req: &request.CreateScheduleRequest{
				CoordinatorID: "coordinator-id",
				ShippingID:    "shipping-id",
				Title:         "スケジュールタイトル",
				Description:   "スケジュールの説明",
				ThumbnailURL:  "https://and-period.jp/thumbnail01.png",
				StartAt:       1640962800,
				EndAt:         1640962800,
				Lives: []*request.CreateScheduleLive{
					{
						Title:       "配信タイトル",
						Description: "配信の説明",
						ProducerID:  "producer-id",
						ProductIDs:  []string{"product-id"},
						StartAt:     1640962800,
						EndAt:       1640962800,
					},
				},
			},
			expect: &testResponse{
				code: http.StatusBadRequest,
			},
		},
		{
			name: "failed to get shipping",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.user.EXPECT().GetCoordinator(gomock.Any(), coordinatorIn).Return(coordinator, nil)
				mocks.user.EXPECT().MultiGetProducers(gomock.Any(), producersIn).Return(producers, nil).Times(2)
				mocks.store.EXPECT().GetShipping(gomock.Any(), shippingIn).Return(nil, assert.AnError)
				mocks.store.EXPECT().MultiGetCategories(gomock.Any(), categoriesIn).Return(categories, nil)
				mocks.store.EXPECT().MultiGetProductTypes(gomock.Any(), productTypesIn).Return(productTypes, nil)
				mocks.store.EXPECT().MultiGetProducts(gomock.Any(), productsIn).Return(products, nil)
			},
			options: []testOption{withRole(uentity.AdminRoleAdministrator)},
			req: &request.CreateScheduleRequest{
				CoordinatorID: "coordinator-id",
				ShippingID:    "shipping-id",
				Title:         "スケジュールタイトル",
				Description:   "スケジュールの説明",
				ThumbnailURL:  "https://and-period.jp/thumbnail01.png",
				StartAt:       1640962800,
				EndAt:         1640962800,
				Lives: []*request.CreateScheduleLive{
					{
						Title:       "配信タイトル",
						Description: "配信の説明",
						ProducerID:  "producer-id",
						ProductIDs:  []string{"product-id"},
						StartAt:     1640962800,
						EndAt:       1640962800,
					},
				},
			},
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
		{
			name: "failed to multi get products",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.user.EXPECT().GetCoordinator(gomock.Any(), coordinatorIn).Return(coordinator, nil)
				mocks.user.EXPECT().MultiGetProducers(gomock.Any(), producersIn).Return(producers, nil)
				mocks.store.EXPECT().GetShipping(gomock.Any(), shippingIn).Return(shipping, nil)
				mocks.store.EXPECT().MultiGetProducts(gomock.Any(), productsIn).Return(nil, assert.AnError)
			},
			options: []testOption{withRole(uentity.AdminRoleAdministrator)},
			req: &request.CreateScheduleRequest{
				CoordinatorID: "coordinator-id",
				ShippingID:    "shipping-id",
				Title:         "スケジュールタイトル",
				Description:   "スケジュールの説明",
				ThumbnailURL:  "https://and-period.jp/thumbnail01.png",
				StartAt:       1640962800,
				EndAt:         1640962800,
				Lives: []*request.CreateScheduleLive{
					{
						Title:       "配信タイトル",
						Description: "配信の説明",
						ProducerID:  "producer-id",
						ProductIDs:  []string{"product-id"},
						StartAt:     1640962800,
						EndAt:       1640962800,
					},
				},
			},
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
		{
			name: "failed to contain invalid product id",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.user.EXPECT().GetCoordinator(gomock.Any(), coordinatorIn).Return(coordinator, nil)
				mocks.user.EXPECT().MultiGetProducers(gomock.Any(), producersIn).Return(producers, nil)
				mocks.store.EXPECT().GetShipping(gomock.Any(), shippingIn).Return(shipping, nil)
				mocks.store.EXPECT().MultiGetProducts(gomock.Any(), productsIn).Return(sentity.Products{}, nil)
			},
			options: []testOption{withRole(uentity.AdminRoleAdministrator)},
			req: &request.CreateScheduleRequest{
				CoordinatorID: "coordinator-id",
				ShippingID:    "shipping-id",
				Title:         "スケジュールタイトル",
				Description:   "スケジュールの説明",
				ThumbnailURL:  "https://and-period.jp/thumbnail01.png",
				StartAt:       1640962800,
				EndAt:         1640962800,
				Lives: []*request.CreateScheduleLive{
					{
						Title:       "配信タイトル",
						Description: "配信の説明",
						ProducerID:  "producer-id",
						ProductIDs:  []string{"product-id"},
						StartAt:     1640962800,
						EndAt:       1640962800,
					},
				},
			},
			expect: &testResponse{
				code: http.StatusBadRequest,
			},
		},
		{
			name: "failed to create schedule",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.user.EXPECT().GetCoordinator(gomock.Any(), coordinatorIn).Return(coordinator, nil)
				mocks.user.EXPECT().MultiGetProducers(gomock.Any(), producersIn).Return(producers, nil).Times(2)
				mocks.store.EXPECT().GetShipping(gomock.Any(), shippingIn).Return(shipping, nil)
				mocks.store.EXPECT().MultiGetCategories(gomock.Any(), categoriesIn).Return(categories, nil)
				mocks.store.EXPECT().MultiGetProductTypes(gomock.Any(), productTypesIn).Return(productTypes, nil)
				mocks.store.EXPECT().MultiGetProducts(gomock.Any(), productsIn).Return(products, nil)
				mocks.store.EXPECT().CreateSchedule(gomock.Any(), scheduleIn).Return(nil, nil, assert.AnError)
			},
			options: []testOption{withRole(uentity.AdminRoleAdministrator)},
			req: &request.CreateScheduleRequest{
				CoordinatorID: "coordinator-id",
				ShippingID:    "shipping-id",
				Title:         "スケジュールタイトル",
				Description:   "スケジュールの説明",
				ThumbnailURL:  "https://and-period.jp/thumbnail01.png",
				StartAt:       1640962800,
				EndAt:         1640962800,
				Lives: []*request.CreateScheduleLive{
					{
						Title:       "配信タイトル",
						Description: "配信の説明",
						ProducerID:  "producer-id",
						ProductIDs:  []string{"product-id"},
						StartAt:     1640962800,
						EndAt:       1640962800,
					},
				},
			},
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			const path = "/v1/schedules"
			testPost(t, tt.setup, tt.expect, path, tt.req, tt.options...)
		})
	}
}
