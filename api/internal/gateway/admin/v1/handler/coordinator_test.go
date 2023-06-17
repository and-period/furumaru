package handler

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/and-period/furumaru/api/internal/common"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/request"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/service"
	"github.com/and-period/furumaru/api/internal/media"
	"github.com/and-period/furumaru/api/internal/store"
	sentity "github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/internal/user"
	"github.com/and-period/furumaru/api/internal/user/entity"
	uentity "github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestListCoordinator(t *testing.T) {
	t.Parallel()

	in := &user.ListCoordinatorsInput{
		Limit:  20,
		Offset: 0,
	}
	coordinators := uentity.Coordinators{
		{
			Admin: uentity.Admin{
				ID:            "coordinator-id01",
				Role:          entity.AdminRoleCoordinator,
				Status:        entity.AdminStatusActivated,
				Lastname:      "&.",
				Firstname:     "管理者",
				LastnameKana:  "あんどどっと",
				FirstnameKana: "かんりしゃ",
				Email:         "test-coordinator@and-period.jp",
			},
			AdminID:        "coordinator-id01",
			MarcheName:     "&.マルシェ",
			Username:       "&.農園",
			Profile:        "紹介文です。",
			ProductTypeIDs: []string{"product-type-id"},
			ThumbnailURL:   "https://and-period.jp/thumbnail.png",
			Thumbnails: common.Images{
				{URL: "https://and-period.jp/thumbnail_240.png", Size: common.ImageSizeSmall},
				{URL: "https://and-period.jp/thumbnail_675.png", Size: common.ImageSizeMedium},
				{URL: "https://and-period.jp/thumbnail_900.png", Size: common.ImageSizeLarge},
			},
			HeaderURL: "https://and-period.jp/header.png",
			Headers: common.Images{
				{URL: "https://and-period.jp/header_240.png", Size: common.ImageSizeSmall},
				{URL: "https://and-period.jp/header_675.png", Size: common.ImageSizeMedium},
				{URL: "https://and-period.jp/header_900.png", Size: common.ImageSizeLarge},
			},
			PromotionVideoURL: "https://and-period.jp/promotion.mp4",
			BonusVideoURL:     "https://and-period.jp/bonus.mp4",
			InstagramID:       "instagram-id",
			FacebookID:        "facebook-id",
			PhoneNumber:       "+819012345678",
			PostalCode:        "1000014",
			Prefecture:        "東京都",
			City:              "千代田区",
			CreatedAt:         jst.Date(2022, 1, 1, 0, 0, 0, 0),
			UpdatedAt:         jst.Date(2022, 1, 1, 0, 0, 0, 0),
		},
		{
			Admin: uentity.Admin{
				ID:            "coordinator-id02",
				Role:          entity.AdminRoleCoordinator,
				Status:        entity.AdminStatusActivated,
				Lastname:      "&.",
				Firstname:     "管理者",
				LastnameKana:  "あんどどっと",
				FirstnameKana: "かんりしゃ",
				Email:         "test-coordinator@and-period.jp",
			},
			AdminID:           "coordinator-id02",
			MarcheName:        "&.マルシェ",
			Username:          "&.農園",
			Profile:           "紹介文です。",
			ProductTypeIDs:    []string{"product-type-id"},
			ThumbnailURL:      "https://and-period.jp/thumbnail.png",
			HeaderURL:         "https://and-period.jp/header.png",
			PromotionVideoURL: "https://and-period.jp/promotion.mp4",
			BonusVideoURL:     "https://and-period.jp/bonus.mp4",
			InstagramID:       "instagram-id",
			FacebookID:        "facebook-id",
			PhoneNumber:       "+819012345678",
			PostalCode:        "1000014",
			Prefecture:        "東京都",
			City:              "千代田区",
			CreatedAt:         jst.Date(2022, 1, 1, 0, 0, 0, 0),
			UpdatedAt:         jst.Date(2022, 1, 1, 0, 0, 0, 0),
		},
	}

	tests := []struct {
		name   string
		setup  func(t *testing.T, mocks *mocks, ctrl *gomock.Controller)
		query  string
		expect *testResponse
	}{
		{
			name: "success",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.user.EXPECT().ListCoordinators(gomock.Any(), in).Return(coordinators, int64(2), nil)
			},
			query: "",
			expect: &testResponse{
				code: http.StatusOK,
				body: &response.CoordinatorsResponse{
					Coordinators: []*response.Coordinator{
						{
							ID:             "coordinator-id01",
							Status:         entity.AdminStatusActivated,
							Lastname:       "&.",
							Firstname:      "管理者",
							LastnameKana:   "あんどどっと",
							FirstnameKana:  "かんりしゃ",
							MarcheName:     "&.マルシェ",
							Username:       "&.農園",
							Profile:        "紹介文です。",
							ProductTypeIDs: []string{"product-type-id"},
							ThumbnailURL:   "https://and-period.jp/thumbnail.png",
							Thumbnails: []*response.Image{
								{URL: "https://and-period.jp/thumbnail_240.png", Size: int32(service.ImageSizeSmall)},
								{URL: "https://and-period.jp/thumbnail_675.png", Size: int32(service.ImageSizeMedium)},
								{URL: "https://and-period.jp/thumbnail_900.png", Size: int32(service.ImageSizeLarge)},
							},
							HeaderURL: "https://and-period.jp/header.png",
							Headers: []*response.Image{
								{URL: "https://and-period.jp/header_240.png", Size: int32(service.ImageSizeSmall)},
								{URL: "https://and-period.jp/header_675.png", Size: int32(service.ImageSizeMedium)},
								{URL: "https://and-period.jp/header_900.png", Size: int32(service.ImageSizeLarge)},
							},
							PromotionVideoURL: "https://and-period.jp/promotion.mp4",
							BonusVideoURL:     "https://and-period.jp/bonus.mp4",
							InstagramID:       "instagram-id",
							FacebookID:        "facebook-id",
							Email:             "test-coordinator@and-period.jp",
							PhoneNumber:       "+819012345678",
							PostalCode:        "1000014",
							Prefecture:        "東京都",
							City:              "千代田区",
							CreatedAt:         1640962800,
							UpdatedAt:         1640962800,
						},
						{
							ID:                "coordinator-id02",
							Status:            entity.AdminStatusActivated,
							Lastname:          "&.",
							Firstname:         "管理者",
							LastnameKana:      "あんどどっと",
							FirstnameKana:     "かんりしゃ",
							MarcheName:        "&.マルシェ",
							Username:          "&.農園",
							Profile:           "紹介文です。",
							ProductTypeIDs:    []string{"product-type-id"},
							ThumbnailURL:      "https://and-period.jp/thumbnail.png",
							Thumbnails:        []*response.Image{},
							HeaderURL:         "https://and-period.jp/header.png",
							Headers:           []*response.Image{},
							PromotionVideoURL: "https://and-period.jp/promotion.mp4",
							BonusVideoURL:     "https://and-period.jp/bonus.mp4",
							InstagramID:       "instagram-id",
							FacebookID:        "facebook-id",
							Email:             "test-coordinator@and-period.jp",
							PhoneNumber:       "+819012345678",
							PostalCode:        "1000014",
							Prefecture:        "東京都",
							City:              "千代田区",
							CreatedAt:         1640962800,
							UpdatedAt:         1640962800,
						},
					},
					Total: 2,
				},
			},
		},
		{
			name: "invalid limit",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
			},
			query: "?limit=a",
			expect: &testResponse{
				code: http.StatusBadRequest,
			},
		},
		{
			name: "invalid offset",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
			},
			query: "?offset=a",
			expect: &testResponse{
				code: http.StatusBadRequest,
			},
		},
		{
			name: "failed to get coordinators",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.user.EXPECT().ListCoordinators(gomock.Any(), in).Return(nil, int64(0), assert.AnError)
			},
			query: "",
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			const format = "/v1/coordinators%s"
			path := fmt.Sprintf(format, tt.query)
			testGet(t, tt.setup, tt.expect, path)
		})
	}
}

func TestGetCoordinator(t *testing.T) {
	t.Parallel()

	in := &user.GetCoordinatorInput{
		CoordinatorID: "coordinator-id",
	}
	coordinator := &uentity.Coordinator{
		Admin: uentity.Admin{
			ID:            "coordinator-id",
			Role:          entity.AdminRoleCoordinator,
			Status:        entity.AdminStatusActivated,
			Lastname:      "&.",
			Firstname:     "管理者",
			LastnameKana:  "あんどどっと",
			FirstnameKana: "かんりしゃ",
			Email:         "test-coordinator@and-period.jp",
		},
		AdminID:        "coordinator-id",
		MarcheName:     "&.マルシェ",
		Username:       "&.農園",
		Profile:        "紹介文です。",
		ProductTypeIDs: []string{"product-type-id"},
		ThumbnailURL:   "https://and-period.jp/thumbnail.png",
		Thumbnails: common.Images{
			{URL: "https://and-period.jp/thumbnail_240.png", Size: common.ImageSizeSmall},
			{URL: "https://and-period.jp/thumbnail_675.png", Size: common.ImageSizeMedium},
			{URL: "https://and-period.jp/thumbnail_900.png", Size: common.ImageSizeLarge},
		},
		HeaderURL: "https://and-period.jp/header.png",
		Headers: common.Images{
			{URL: "https://and-period.jp/header_240.png", Size: common.ImageSizeSmall},
			{URL: "https://and-period.jp/header_675.png", Size: common.ImageSizeMedium},
			{URL: "https://and-period.jp/header_900.png", Size: common.ImageSizeLarge},
		},
		PromotionVideoURL: "https://and-period.jp/promotion.mp4",
		BonusVideoURL:     "https://and-period.jp/bonus.mp4",
		InstagramID:       "instagram-id",
		FacebookID:        "facebook-id",
		PhoneNumber:       "+819012345678",
		PostalCode:        "1000014",
		Prefecture:        "東京都",
		City:              "千代田区",
		CreatedAt:         jst.Date(2022, 1, 1, 0, 0, 0, 0),
		UpdatedAt:         jst.Date(2022, 1, 1, 0, 0, 0, 0),
	}

	tests := []struct {
		name          string
		setup         func(t *testing.T, mocks *mocks, ctrl *gomock.Controller)
		coordinatorID string
		expect        *testResponse
	}{
		{
			name: "success",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.user.EXPECT().GetCoordinator(gomock.Any(), in).Return(coordinator, nil)
			},
			coordinatorID: "coordinator-id",
			expect: &testResponse{
				code: http.StatusOK,
				body: &response.CoordinatorResponse{
					Coordinator: &response.Coordinator{
						ID:             "coordinator-id",
						Status:         entity.AdminStatusActivated,
						Lastname:       "&.",
						Firstname:      "管理者",
						LastnameKana:   "あんどどっと",
						FirstnameKana:  "かんりしゃ",
						MarcheName:     "&.マルシェ",
						Username:       "&.農園",
						Profile:        "紹介文です。",
						ProductTypeIDs: []string{"product-type-id"},
						ThumbnailURL:   "https://and-period.jp/thumbnail.png",
						Thumbnails: []*response.Image{
							{URL: "https://and-period.jp/thumbnail_240.png", Size: int32(service.ImageSizeSmall)},
							{URL: "https://and-period.jp/thumbnail_675.png", Size: int32(service.ImageSizeMedium)},
							{URL: "https://and-period.jp/thumbnail_900.png", Size: int32(service.ImageSizeLarge)},
						},
						HeaderURL: "https://and-period.jp/header.png",
						Headers: []*response.Image{
							{URL: "https://and-period.jp/header_240.png", Size: int32(service.ImageSizeSmall)},
							{URL: "https://and-period.jp/header_675.png", Size: int32(service.ImageSizeMedium)},
							{URL: "https://and-period.jp/header_900.png", Size: int32(service.ImageSizeLarge)},
						},
						PromotionVideoURL: "https://and-period.jp/promotion.mp4",
						BonusVideoURL:     "https://and-period.jp/bonus.mp4",
						InstagramID:       "instagram-id",
						FacebookID:        "facebook-id",
						Email:             "test-coordinator@and-period.jp",
						PhoneNumber:       "+819012345678",
						PostalCode:        "1000014",
						Prefecture:        "東京都",
						City:              "千代田区",
						CreatedAt:         1640962800,
						UpdatedAt:         1640962800,
					},
				},
			},
		},
		{
			name: "failed to get coordinator",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.user.EXPECT().GetCoordinator(gomock.Any(), in).Return(nil, assert.AnError)
			},
			coordinatorID: "coordinator-id",
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			const format = "/v1/coordinators/%s"
			path := fmt.Sprintf(format, tt.coordinatorID)
			testGet(t, tt.setup, tt.expect, path)
		})
	}
}

func TestCreateCoordinator(t *testing.T) {
	t.Parallel()

	productTypesIn := &store.MultiGetProductTypesInput{
		ProductTypeIDs: []string{"product-type-id"},
	}
	productTypes := sentity.ProductTypes{
		{
			ID:         "product-type-id",
			CategoryID: "category-id",
			Name:       "じゃがいも",
		},
	}
	categoriesIn := &store.MultiGetCategoriesInput{
		CategoryIDs: []string{"category-id"},
	}
	categories := sentity.Categories{
		{
			ID:   "category-id",
			Name: "野菜",
		},
	}
	thumbnailIn := &media.UploadFileInput{
		URL: "https://tmp.and-period.jp/thumbnail.png",
	}
	thumbnailURL := "https://and-period.jp/thumbnail.png"
	headerIn := &media.UploadFileInput{
		URL: "https://tmp.and-period.jp/header.png",
	}
	headerURL := "https://and-period.jp/header.png"
	promotionVideoIn := &media.UploadFileInput{
		URL: "https://tmp.and-period.jp/promotion.mp4",
	}
	promotionVideoURL := "https://and-period.jp/promotion.mp4"
	bonusVideoIn := &media.UploadFileInput{
		URL: "https://tmp.and-period.jp/bonus.mp4",
	}
	bonusVideoURL := "https://and-period.jp/bonus.mp4"
	in := &user.CreateCoordinatorInput{
		Lastname:          "&.",
		Firstname:         "生産者",
		LastnameKana:      "あんどどっと",
		FirstnameKana:     "せいさんしゃ",
		MarcheName:        "&.マルシェ",
		Username:          "&.農園",
		Profile:           "紹介文です。",
		ProductTypeIDs:    []string{"product-type-id"},
		ThumbnailURL:      thumbnailURL,
		HeaderURL:         headerURL,
		PromotionVideoURL: "https://and-period.jp/promotion.mp4",
		BonusVideoURL:     "https://and-period.jp/bonus.mp4",
		InstagramID:       "instagram-id",
		FacebookID:        "facebook-id",
		Email:             "test-coordinator@and-period.jp",
		PhoneNumber:       "+819012345678",
		PostalCode:        "1000014",
		Prefecture:        "東京都",
		City:              "千代田区",
		AddressLine1:      "永田町1-7-1",
		AddressLine2:      "",
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
		AdminID:           "coordinator-id",
		MarcheName:        "&.マルシェ",
		Username:          "&.農園",
		Profile:           "紹介文です。",
		ProductTypeIDs:    []string{"product-type-id"},
		ThumbnailURL:      thumbnailURL,
		HeaderURL:         headerURL,
		PromotionVideoURL: "https://and-period.jp/promotion.mp4",
		BonusVideoURL:     "https://and-period.jp/bonus.mp4",
		InstagramID:       "instagram-id",
		FacebookID:        "facebook-id",
		PhoneNumber:       "+819012345678",
		PostalCode:        "1000014",
		Prefecture:        "東京都",
		City:              "千代田区",
		AddressLine1:      "永田町1-7-1",
		AddressLine2:      "",
		CreatedAt:         jst.Date(2022, 1, 1, 0, 0, 0, 0),
		UpdatedAt:         jst.Date(2022, 1, 1, 0, 0, 0, 0),
	}

	tests := []struct {
		name   string
		setup  func(t *testing.T, mocks *mocks, ctrl *gomock.Controller)
		req    *request.CreateCoordinatorRequest
		expect *testResponse
	}{
		{
			name: "success",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.store.EXPECT().MultiGetProductTypes(gomock.Any(), productTypesIn).Return(productTypes, nil)
				mocks.store.EXPECT().MultiGetCategories(gomock.Any(), categoriesIn).Return(categories, nil)
				mocks.media.EXPECT().UploadCoordinatorThumbnail(gomock.Any(), thumbnailIn).Return(thumbnailURL, nil)
				mocks.media.EXPECT().UploadCoordinatorHeader(gomock.Any(), headerIn).Return(headerURL, nil)
				mocks.media.EXPECT().UploadCoordinatorPromotionVideo(gomock.Any(), promotionVideoIn).Return(promotionVideoURL, nil)
				mocks.media.EXPECT().UploadCoordinatorBonusVideo(gomock.Any(), bonusVideoIn).Return(bonusVideoURL, nil)
				mocks.user.EXPECT().CreateCoordinator(gomock.Any(), in).Return(coordinator, nil)
			},
			req: &request.CreateCoordinatorRequest{
				Lastname:          "&.",
				Firstname:         "生産者",
				LastnameKana:      "あんどどっと",
				FirstnameKana:     "せいさんしゃ",
				MarcheName:        "&.マルシェ",
				Username:          "&.農園",
				Profile:           "紹介文です。",
				ProductTypeIDs:    []string{"product-type-id"},
				ThumbnailURL:      "https://tmp.and-period.jp/thumbnail.png",
				HeaderURL:         "https://tmp.and-period.jp/header.png",
				PromotionVideoURL: "https://tmp.and-period.jp/promotion.mp4",
				BonusVideoURL:     "https://tmp.and-period.jp/bonus.mp4",
				InstagramID:       "instagram-id",
				FacebookID:        "facebook-id",
				Email:             "test-coordinator@and-period.jp",
				PhoneNumber:       "+819012345678",
				PostalCode:        "1000014",
				Prefecture:        "東京都",
				City:              "千代田区",
				AddressLine1:      "永田町1-7-1",
				AddressLine2:      "",
			},
			expect: &testResponse{
				code: http.StatusOK,
				body: &response.CoordinatorResponse{
					Coordinator: &response.Coordinator{
						ID:                "coordinator-id",
						Lastname:          "&.",
						Firstname:         "管理者",
						LastnameKana:      "あんどどっと",
						FirstnameKana:     "かんりしゃ",
						MarcheName:        "&.マルシェ",
						Username:          "&.農園",
						Profile:           "紹介文です。",
						ProductTypeIDs:    []string{"product-type-id"},
						ThumbnailURL:      "https://and-period.jp/thumbnail.png",
						Thumbnails:        []*response.Image{},
						HeaderURL:         "https://and-period.jp/header.png",
						Headers:           []*response.Image{},
						PromotionVideoURL: "https://and-period.jp/promotion.mp4",
						BonusVideoURL:     "https://and-period.jp/bonus.mp4",
						InstagramID:       "instagram-id",
						FacebookID:        "facebook-id",
						Email:             "test-coordinator@and-period.jp",
						PhoneNumber:       "+819012345678",
						PostalCode:        "1000014",
						Prefecture:        "東京都",
						City:              "千代田区",
						AddressLine1:      "永田町1-7-1",
						AddressLine2:      "",
						CreatedAt:         1640962800,
						UpdatedAt:         1640962800,
					},
				},
			},
		},
		{
			name: "failed to upload coordinator thumbnail",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.store.EXPECT().MultiGetProductTypes(gomock.Any(), productTypesIn).Return(productTypes, nil)
				mocks.store.EXPECT().MultiGetCategories(gomock.Any(), categoriesIn).Return(categories, nil)
				mocks.media.EXPECT().UploadCoordinatorThumbnail(gomock.Any(), thumbnailIn).Return("", assert.AnError)
				mocks.media.EXPECT().UploadCoordinatorHeader(gomock.Any(), headerIn).Return(headerURL, nil)
				mocks.media.EXPECT().UploadCoordinatorPromotionVideo(gomock.Any(), promotionVideoIn).Return(promotionVideoURL, nil)
				mocks.media.EXPECT().UploadCoordinatorBonusVideo(gomock.Any(), bonusVideoIn).Return(bonusVideoURL, nil)
			},
			req: &request.CreateCoordinatorRequest{
				Lastname:          "&.",
				Firstname:         "生産者",
				LastnameKana:      "あんどどっと",
				FirstnameKana:     "せいさんしゃ",
				MarcheName:        "&.マルシェ",
				Username:          "&.農園",
				Profile:           "紹介文です。",
				ProductTypeIDs:    []string{"product-type-id"},
				ThumbnailURL:      "https://tmp.and-period.jp/thumbnail.png",
				HeaderURL:         "https://tmp.and-period.jp/header.png",
				PromotionVideoURL: "https://tmp.and-period.jp/promotion.mp4",
				BonusVideoURL:     "https://tmp.and-period.jp/bonus.mp4",
				InstagramID:       "instagram-id",
				FacebookID:        "facebook-id",
				Email:             "test-coordinator@and-period.jp",
				PhoneNumber:       "+819012345678",
				PostalCode:        "1000014",
				Prefecture:        "東京都",
				City:              "千代田区",
				AddressLine1:      "永田町1-7-1",
				AddressLine2:      "",
			},
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
		{
			name: "failed to upload coordinator thumbnail",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.store.EXPECT().MultiGetProductTypes(gomock.Any(), productTypesIn).Return(productTypes, nil)
				mocks.store.EXPECT().MultiGetCategories(gomock.Any(), categoriesIn).Return(categories, nil)
				mocks.media.EXPECT().UploadCoordinatorThumbnail(gomock.Any(), thumbnailIn).Return(thumbnailURL, nil)
				mocks.media.EXPECT().UploadCoordinatorHeader(gomock.Any(), headerIn).Return("", assert.AnError)
				mocks.media.EXPECT().UploadCoordinatorPromotionVideo(gomock.Any(), promotionVideoIn).Return(promotionVideoURL, nil)
				mocks.media.EXPECT().UploadCoordinatorBonusVideo(gomock.Any(), bonusVideoIn).Return(bonusVideoURL, nil)
			},
			req: &request.CreateCoordinatorRequest{
				Lastname:          "&.",
				Firstname:         "生産者",
				LastnameKana:      "あんどどっと",
				FirstnameKana:     "せいさんしゃ",
				MarcheName:        "&.マルシェ",
				Username:          "&.農園",
				Profile:           "紹介文です。",
				ProductTypeIDs:    []string{"product-type-id"},
				ThumbnailURL:      "https://tmp.and-period.jp/thumbnail.png",
				HeaderURL:         "https://tmp.and-period.jp/header.png",
				PromotionVideoURL: "https://tmp.and-period.jp/promotion.mp4",
				BonusVideoURL:     "https://tmp.and-period.jp/bonus.mp4",
				InstagramID:       "instagram-id",
				FacebookID:        "facebook-id",
				Email:             "test-coordinator@and-period.jp",
				PhoneNumber:       "+819012345678",
				PostalCode:        "1000014",
				Prefecture:        "東京都",
				City:              "千代田区",
				AddressLine1:      "永田町1-7-1",
				AddressLine2:      "",
			},
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
		{
			name: "failed to create coordinator",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				in := *in
				in.ThumbnailURL, in.HeaderURL = "", ""
				in.PromotionVideoURL, in.BonusVideoURL = "", ""
				mocks.store.EXPECT().MultiGetProductTypes(gomock.Any(), productTypesIn).Return(productTypes, nil)
				mocks.store.EXPECT().MultiGetCategories(gomock.Any(), categoriesIn).Return(categories, nil)
				mocks.user.EXPECT().CreateCoordinator(gomock.Any(), &in).Return(nil, assert.AnError)
			},
			req: &request.CreateCoordinatorRequest{
				Lastname:          "&.",
				Firstname:         "生産者",
				LastnameKana:      "あんどどっと",
				FirstnameKana:     "せいさんしゃ",
				MarcheName:        "&.マルシェ",
				Username:          "&.農園",
				Profile:           "紹介文です。",
				ProductTypeIDs:    []string{"product-type-id"},
				ThumbnailURL:      "",
				HeaderURL:         "",
				PromotionVideoURL: "",
				BonusVideoURL:     "",
				InstagramID:       "instagram-id",
				FacebookID:        "facebook-id",
				Email:             "test-coordinator@and-period.jp",
				PhoneNumber:       "+819012345678",
				PostalCode:        "1000014",
				Prefecture:        "東京都",
				City:              "千代田区",
				AddressLine1:      "永田町1-7-1",
				AddressLine2:      "",
			},
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			const path = "/v1/coordinators"
			testPost(t, tt.setup, tt.expect, path, tt.req)
		})
	}
}

func TestUpdateCoordinator(t *testing.T) {
	t.Parallel()

	productTypesIn := &store.MultiGetProductTypesInput{
		ProductTypeIDs: []string{"product-type-id"},
	}
	productTypes := sentity.ProductTypes{
		{
			ID:         "product-type-id",
			CategoryID: "category-id",
			Name:       "じゃがいも",
		},
	}
	categoriesIn := &store.MultiGetCategoriesInput{
		CategoryIDs: []string{"category-id"},
	}
	categories := sentity.Categories{
		{
			ID:   "category-id",
			Name: "野菜",
		},
	}
	thumbnailIn := &media.UploadFileInput{
		URL: "https://tmp.and-period.jp/thumbnail.png",
	}
	thumbnailURL := "https://and-period.jp/thumbnail.png"
	headerIn := &media.UploadFileInput{
		URL: "https://tmp.and-period.jp/header.png",
	}
	headerURL := "https://and-period.jp/header.png"
	in := &user.UpdateCoordinatorInput{
		CoordinatorID:  "coordinator-id",
		Lastname:       "&.",
		Firstname:      "生産者",
		LastnameKana:   "あんどどっと",
		FirstnameKana:  "せいさんしゃ",
		MarcheName:     "&.マルシェ",
		Username:       "&.農園",
		Profile:        "紹介文です。",
		ProductTypeIDs: []string{"product-type-id"},
		ThumbnailURL:   "https://and-period.jp/thumbnail.png",
		HeaderURL:      "https://and-period.jp/header.png",
		PhoneNumber:    "+819012345678",
		PostalCode:     "1000014",
		Prefecture:     "東京都",
		City:           "千代田区",
		AddressLine1:   "永田町1-7-1",
		AddressLine2:   "",
	}

	tests := []struct {
		name          string
		setup         func(t *testing.T, mocks *mocks, ctrl *gomock.Controller)
		coordinatorID string
		req           *request.UpdateCoordinatorRequest
		expect        *testResponse
	}{
		{
			name: "success",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.store.EXPECT().MultiGetProductTypes(gomock.Any(), productTypesIn).Return(productTypes, nil)
				mocks.store.EXPECT().MultiGetCategories(gomock.Any(), categoriesIn).Return(categories, nil)
				mocks.media.EXPECT().UploadCoordinatorThumbnail(gomock.Any(), thumbnailIn).Return(thumbnailURL, nil)
				mocks.media.EXPECT().UploadCoordinatorHeader(gomock.Any(), headerIn).Return(headerURL, nil)
				mocks.user.EXPECT().UpdateCoordinator(gomock.Any(), in).Return(nil)
			},
			coordinatorID: "coordinator-id",
			req: &request.UpdateCoordinatorRequest{
				Lastname:       "&.",
				Firstname:      "生産者",
				LastnameKana:   "あんどどっと",
				FirstnameKana:  "せいさんしゃ",
				MarcheName:     "&.マルシェ",
				Username:       "&.農園",
				Profile:        "紹介文です。",
				ProductTypeIDs: []string{"product-type-id"},
				ThumbnailURL:   "https://tmp.and-period.jp/thumbnail.png",
				HeaderURL:      "https://tmp.and-period.jp/header.png",
				PhoneNumber:    "+819012345678",
				PostalCode:     "1000014",
				Prefecture:     "東京都",
				City:           "千代田区",
				AddressLine1:   "永田町1-7-1",
				AddressLine2:   "",
			},
			expect: &testResponse{
				code: http.StatusNoContent,
			},
		},
		{
			name: "failed to upload coordinator thumbnail",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.store.EXPECT().MultiGetProductTypes(gomock.Any(), productTypesIn).Return(productTypes, nil)
				mocks.store.EXPECT().MultiGetCategories(gomock.Any(), categoriesIn).Return(categories, nil)
				mocks.media.EXPECT().UploadCoordinatorThumbnail(gomock.Any(), thumbnailIn).Return("", assert.AnError)
				mocks.media.EXPECT().UploadCoordinatorHeader(gomock.Any(), headerIn).Return(headerURL, nil)
			},
			coordinatorID: "coordinator-id",
			req: &request.UpdateCoordinatorRequest{
				Lastname:       "&.",
				Firstname:      "生産者",
				LastnameKana:   "あんどどっと",
				FirstnameKana:  "せいさんしゃ",
				MarcheName:     "&.マルシェ",
				Username:       "&.農園",
				Profile:        "紹介文です。",
				ProductTypeIDs: []string{"product-type-id"},
				ThumbnailURL:   "https://tmp.and-period.jp/thumbnail.png",
				HeaderURL:      "https://tmp.and-period.jp/header.png",
				PhoneNumber:    "+819012345678",
				PostalCode:     "1000014",
				Prefecture:     "東京都",
				City:           "千代田区",
				AddressLine1:   "永田町1-7-1",
				AddressLine2:   "",
			},
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
		{
			name: "failed to upload coordinator header",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.store.EXPECT().MultiGetProductTypes(gomock.Any(), productTypesIn).Return(productTypes, nil)
				mocks.store.EXPECT().MultiGetCategories(gomock.Any(), categoriesIn).Return(categories, nil)
				mocks.media.EXPECT().UploadCoordinatorThumbnail(gomock.Any(), thumbnailIn).Return(thumbnailURL, nil)
				mocks.media.EXPECT().UploadCoordinatorHeader(gomock.Any(), headerIn).Return("", assert.AnError)
			},
			coordinatorID: "coordinator-id",
			req: &request.UpdateCoordinatorRequest{
				Lastname:       "&.",
				Firstname:      "生産者",
				LastnameKana:   "あんどどっと",
				FirstnameKana:  "せいさんしゃ",
				MarcheName:     "&.マルシェ",
				Username:       "&.農園",
				Profile:        "紹介文です。",
				ProductTypeIDs: []string{"product-type-id"},
				ThumbnailURL:   "https://tmp.and-period.jp/thumbnail.png",
				HeaderURL:      "https://tmp.and-period.jp/header.png",
				PhoneNumber:    "+819012345678",
				PostalCode:     "1000014",
				Prefecture:     "東京都",
				City:           "千代田区",
				AddressLine1:   "永田町1-7-1",
				AddressLine2:   "",
			},
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
		{
			name: "failed to update coordinator",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				in := *in
				in.ThumbnailURL = ""
				in.HeaderURL = ""
				mocks.store.EXPECT().MultiGetProductTypes(gomock.Any(), productTypesIn).Return(productTypes, nil)
				mocks.store.EXPECT().MultiGetCategories(gomock.Any(), categoriesIn).Return(categories, nil)
				mocks.user.EXPECT().UpdateCoordinator(gomock.Any(), &in).Return(assert.AnError)
			},
			coordinatorID: "coordinator-id",
			req: &request.UpdateCoordinatorRequest{
				Lastname:       "&.",
				Firstname:      "生産者",
				LastnameKana:   "あんどどっと",
				FirstnameKana:  "せいさんしゃ",
				MarcheName:     "&.マルシェ",
				Username:       "&.農園",
				Profile:        "紹介文です。",
				ProductTypeIDs: []string{"product-type-id"},
				ThumbnailURL:   "",
				HeaderURL:      "",
				PhoneNumber:    "+819012345678",
				PostalCode:     "1000014",
				Prefecture:     "東京都",
				City:           "千代田区",
				AddressLine1:   "永田町1-7-1",
				AddressLine2:   "",
			},
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			const format = "/v1/coordinators/%s"
			path := fmt.Sprintf(format, tt.coordinatorID)
			testPatch(t, tt.setup, tt.expect, path, tt.req)
		})
	}
}

func TestUpdateCoordinatorEmail(t *testing.T) {
	t.Parallel()

	in := &user.UpdateCoordinatorEmailInput{
		CoordinatorID: "coordinator-id",
		Email:         "test-producer@and-period.jp",
	}

	tests := []struct {
		name          string
		setup         func(t *testing.T, mocks *mocks, ctrl *gomock.Controller)
		coordinatorID string
		req           *request.UpdateCoordinatorEmailRequest
		expect        *testResponse
	}{
		{
			name: "success",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.user.EXPECT().UpdateCoordinatorEmail(gomock.Any(), in).Return(nil)
			},
			coordinatorID: "coordinator-id",
			req: &request.UpdateCoordinatorEmailRequest{
				Email: "test-producer@and-period.jp",
			},
			expect: &testResponse{
				code: http.StatusNoContent,
			},
		},
		{
			name: "failed to update coordinator email",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.user.EXPECT().UpdateCoordinatorEmail(gomock.Any(), in).Return(assert.AnError)
			},
			coordinatorID: "coordinator-id",
			req: &request.UpdateCoordinatorEmailRequest{
				Email: "test-producer@and-period.jp",
			},
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			const format = "/v1/coordinators/%s/email"
			path := fmt.Sprintf(format, tt.coordinatorID)
			testPatch(t, tt.setup, tt.expect, path, tt.req)
		})
	}
}

func TestResetCoordinatorPassword(t *testing.T) {
	t.Parallel()

	in := &user.ResetCoordinatorPasswordInput{
		CoordinatorID: "coordinator-id",
	}

	tests := []struct {
		name          string
		setup         func(t *testing.T, mocks *mocks, ctrl *gomock.Controller)
		coordinatorID string
		expect        *testResponse
	}{
		{
			name: "success",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.user.EXPECT().ResetCoordinatorPassword(gomock.Any(), in).Return(nil)
			},
			coordinatorID: "coordinator-id",
			expect: &testResponse{
				code: http.StatusNoContent,
			},
		},
		{
			name: "failed to reset coordinator password",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.user.EXPECT().ResetCoordinatorPassword(gomock.Any(), in).Return(assert.AnError)
			},
			coordinatorID: "coordinator-id",
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			const format = "/v1/coordinators/%s/password"
			path := fmt.Sprintf(format, tt.coordinatorID)
			testPatch(t, tt.setup, tt.expect, path, nil)
		})
	}
}

func TestDeleteCoordinator(t *testing.T) {
	t.Parallel()

	in := &user.DeleteCoordinatorInput{
		CoordinatorID: "coordinator-id",
	}

	tests := []struct {
		name          string
		setup         func(t *testing.T, mocks *mocks, ctrl *gomock.Controller)
		coordinatorID string
		expect        *testResponse
	}{
		{
			name: "success",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.user.EXPECT().DeleteCoordinator(gomock.Any(), in).Return(nil)
			},
			coordinatorID: "coordinator-id",
			expect: &testResponse{
				code: http.StatusNoContent,
			},
		},
		{
			name: "failed to delete coordinator",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.user.EXPECT().DeleteCoordinator(gomock.Any(), in).Return(assert.AnError)
			},
			coordinatorID: "coordinator-id",
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			const format = "/v1/coordinators/%s"
			path := fmt.Sprintf(format, tt.coordinatorID)
			testDelete(t, tt.setup, tt.expect, path)
		})
	}
}
