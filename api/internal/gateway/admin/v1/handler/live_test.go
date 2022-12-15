package handler

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/and-period/furumaru/api/internal/common"
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

func TestGetLive(t *testing.T) {
	t.Parallel()

	producerIn := &user.GetProducerInput{
		ProducerID: "producer-id",
	}
	producersIn := &user.MultiGetProducersInput{
		ProducerIDs: []string{"producer-id"},
	}
	producer := &uentity.Producer{
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
	}
	producers := uentity.Producers{producer}
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
	liveIn := &store.GetLiveInput{
		LiveID: "live-id",
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
		ID:             "live-id",
		ScheduleID:     "schedule-id",
		Title:          "配信タイトル",
		Description:    "配信の説明",
		Status:         sentity.LiveStatusWaiting,
		Published:      false,
		Canceled:       false,
		ProducerID:     "producer-id",
		ChannelArn:     "channel-arn",
		StreamKeyArn:   "streamKey-arn",
		StartAt:        jst.Date(2022, 1, 1, 0, 0, 0, 0),
		EndAt:          jst.Date(2022, 1, 1, 0, 0, 0, 0),
		CreatedAt:      jst.Date(2022, 1, 1, 0, 0, 0, 0),
		UpdatedAt:      jst.Date(2022, 1, 1, 0, 0, 0, 0),
		ChannelName:    "配信チャンネル",
		IngestEndpoint: "ingest-endpoint",
		StreamKey:      "streamKey-value",
		StreamID:       "stream-id",
		PlaybackURL:    "playback-url",
		ViewerCount:    100,
	}
	tests := []struct {
		name       string
		setup      func(t *testing.T, mocks *mocks, ctrl *gomock.Controller)
		scheduleID string
		liveID     string
		expect     *testResponse
	}{
		{
			name: "success",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.store.EXPECT().GetLive(gomock.Any(), liveIn).Return(live, nil)
				mocks.user.EXPECT().GetProducer(gomock.Any(), producerIn).Return(producer, nil)
				mocks.user.EXPECT().MultiGetProducers(gomock.Any(), producersIn).Return(producers, nil)
				mocks.store.EXPECT().MultiGetCategories(gomock.Any(), categoriesIn).Return(categories, nil)
				mocks.store.EXPECT().MultiGetProductTypes(gomock.Any(), productTypesIn).Return(productTypes, nil)
				mocks.store.EXPECT().MultiGetProducts(gomock.Any(), productsIn).Return(products, nil)
			},
			scheduleID: "schedule-id",
			liveID:     "live-id",
			expect: &testResponse{
				code: http.StatusOK,
				body: &response.LiveResponse{
					Live: &response.Live{
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
						ChannelArn:     "channel-arn",
						StreamKeyArn:   "streamKey-arn",
						StartAt:        1640962800,
						EndAt:          1640962800,
						CreatedAt:      1640962800,
						UpdatedAt:      1640962800,
						ChannelName:    "配信チャンネル",
						IngestEndpoint: "ingest-endpoint",
						StreamKey:      "streamKey-value",
						StreamID:       "stream-id",
						PlaybackURL:    "playback-url",
						ViewerCount:    100,
					},
				},
			},
		},
		{
			name: "failed to get live",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.store.EXPECT().GetLive(gomock.Any(), liveIn).Return(nil, assert.AnError)
			},
			liveID: "live-id",
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
		{
			name: "failed to get producer",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.store.EXPECT().GetLive(gomock.Any(), liveIn).Return(live, nil)
				mocks.user.EXPECT().GetProducer(gomock.Any(), producerIn).Return(nil, assert.AnError)
				mocks.user.EXPECT().MultiGetProducers(gomock.Any(), producersIn).Return(producers, nil)
				mocks.store.EXPECT().MultiGetCategories(gomock.Any(), categoriesIn).Return(categories, nil)
				mocks.store.EXPECT().MultiGetProductTypes(gomock.Any(), productTypesIn).Return(productTypes, nil)
				mocks.store.EXPECT().MultiGetProducts(gomock.Any(), productsIn).Return(products, nil)
			},
			liveID: "live-id",
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
		{
			name: "failed to multi get producers",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.store.EXPECT().GetLive(gomock.Any(), liveIn).Return(live, nil)
				mocks.user.EXPECT().GetProducer(gomock.Any(), producerIn).Return(producer, nil)
				mocks.user.EXPECT().MultiGetProducers(gomock.Any(), producersIn).Return(nil, assert.AnError)
				mocks.store.EXPECT().MultiGetCategories(gomock.Any(), categoriesIn).Return(categories, nil)
				mocks.store.EXPECT().MultiGetProductTypes(gomock.Any(), productTypesIn).Return(productTypes, nil)
				mocks.store.EXPECT().MultiGetProducts(gomock.Any(), productsIn).Return(products, nil)
			},
			liveID: "live-id",
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
		{
			name: "failed to get catogory",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.store.EXPECT().GetLive(gomock.Any(), liveIn).Return(live, nil)
				mocks.user.EXPECT().GetProducer(gomock.Any(), producerIn).Return(producer, nil)
				mocks.user.EXPECT().MultiGetProducers(gomock.Any(), producersIn).Return(producers, nil)
				mocks.store.EXPECT().MultiGetCategories(gomock.Any(), categoriesIn).Return(nil, assert.AnError)
				mocks.store.EXPECT().MultiGetProductTypes(gomock.Any(), productTypesIn).Return(productTypes, nil)
				mocks.store.EXPECT().MultiGetProducts(gomock.Any(), productsIn).Return(products, nil)
			},
			liveID: "live-id",
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
		{
			name: "failed to get product type",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.store.EXPECT().GetLive(gomock.Any(), liveIn).Return(live, nil)
				mocks.user.EXPECT().GetProducer(gomock.Any(), producerIn).Return(producer, nil)
				mocks.user.EXPECT().MultiGetProducers(gomock.Any(), producersIn).Return(producers, nil)
				mocks.store.EXPECT().MultiGetProductTypes(gomock.Any(), productTypesIn).Return(nil, assert.AnError)
				mocks.store.EXPECT().MultiGetProducts(gomock.Any(), productsIn).Return(products, nil)
			},
			liveID: "live-id",
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
		{
			name: "failed to multi get products",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.store.EXPECT().GetLive(gomock.Any(), liveIn).Return(live, nil)
				mocks.user.EXPECT().GetProducer(gomock.Any(), producerIn).Return(producer, nil)
				mocks.store.EXPECT().MultiGetProducts(gomock.Any(), productsIn).Return(nil, assert.AnError)
			},
			liveID: "live-id",
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			const format = "/v1/schedules/%s/lives/%s"
			path := fmt.Sprintf(format, tt.scheduleID, tt.liveID)
			testGet(t, tt.setup, tt.expect, path)
		})
	}
}
