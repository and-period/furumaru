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
	"github.com/and-period/furumaru/api/internal/user"
	"github.com/and-period/furumaru/api/internal/user/entity"
	uentity "github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestFilterAccessProducer(t *testing.T) {
	t.Parallel()

	in := &user.GetProducerInput{
		ProducerID: "producer-id",
	}
	producer := &uentity.Producer{
		Admin: uentity.Admin{
			ID:            "producer-id",
			Role:          entity.AdminRoleProducer,
			Status:        entity.AdminStatusActivated,
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

	tests := []struct {
		name    string
		setup   func(t *testing.T, mocks *mocks, ctrl *gomock.Controller)
		options []testOption
		expect  int
	}{
		{
			name:    "administrator success",
			setup:   func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {},
			options: []testOption{withRole(uentity.AdminRoleAdministrator)},
			expect:  http.StatusOK,
		},
		{
			name: "coordinator success",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.user.EXPECT().GetProducer(gomock.Any(), in).Return(producer, nil)
			},
			options: []testOption{withRole(uentity.AdminRoleCoordinator), withAdminID("coordinator-id")},
			expect:  http.StatusOK,
		},
		{
			name: "coordinator forbidden",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.user.EXPECT().GetProducer(gomock.Any(), in).Return(producer, nil)
			},
			options: []testOption{withRole(uentity.AdminRoleCoordinator)},
			expect:  http.StatusForbidden,
		},
		{
			name: "coordinator failed to get producer",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.user.EXPECT().GetProducer(gomock.Any(), in).Return(nil, assert.AnError)
			},
			options: []testOption{withRole(uentity.AdminRoleCoordinator), withAdminID("coordinator-id")},
			expect:  http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			const route, path = "/producers/:producerId", "/producers/producer-id"
			testMiddleware(t, tt.setup, route, path, tt.expect, func(h *handler) gin.HandlerFunc {
				return h.filterAccessProducer
			}, tt.options...)
		})
	}
}

func TestListProducers(t *testing.T) {
	t.Parallel()

	in := &user.ListProducersInput{
		Limit:  20,
		Offset: 0,
	}
	producers := uentity.Producers{
		{
			Admin: uentity.Admin{
				ID:            "producer-id01",
				Role:          entity.AdminRoleProducer,
				Status:        entity.AdminStatusActivated,
				Lastname:      "&.",
				Firstname:     "管理者",
				LastnameKana:  "あんどどっと",
				FirstnameKana: "かんりしゃ",
				Email:         "test-producer@and-period.jp",
			},
			AdminID:       "producer-id01",
			CoordinatorID: "coordinator-id",
			StoreName:     "&.農園",
			ThumbnailURL:  "https://and-period.jp/thumbnail.png",
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
			PhoneNumber: "+819012345678",
			PostalCode:  "1000014",
			Prefecture:  "東京都",
			City:        "千代田区",
			CreatedAt:   jst.Date(2022, 1, 1, 0, 0, 0, 0),
			UpdatedAt:   jst.Date(2022, 1, 1, 0, 0, 0, 0),
		},
		{
			Admin: uentity.Admin{
				ID:            "producer-id02",
				Role:          entity.AdminRoleProducer,
				Status:        entity.AdminStatusActivated,
				Lastname:      "&.",
				Firstname:     "管理者",
				LastnameKana:  "あんどどっと",
				FirstnameKana: "かんりしゃ",
				Email:         "test-producer@and-period.jp",
			},
			AdminID:       "producer-id02",
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

	tests := []struct {
		name    string
		setup   func(t *testing.T, mocks *mocks, ctrl *gomock.Controller)
		options []testOption
		query   string
		expect  *testResponse
	}{
		{
			name: "success administrator",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.user.EXPECT().ListProducers(gomock.Any(), in).Return(producers, int64(2), nil)
			},
			options: []testOption{withRole(uentity.AdminRoleAdministrator)},
			query:   "",
			expect: &testResponse{
				code: http.StatusOK,
				body: &response.ProducersResponse{
					Producers: []*response.Producer{
						{
							ID:            "producer-id01",
							Status:        entity.AdminStatusActivated,
							CoordinatorID: "coordinator-id",
							Lastname:      "&.",
							Firstname:     "管理者",
							LastnameKana:  "あんどどっと",
							FirstnameKana: "かんりしゃ",
							StoreName:     "&.農園",
							ThumbnailURL:  "https://and-period.jp/thumbnail.png",
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
							Email:       "test-producer@and-period.jp",
							PhoneNumber: "+819012345678",
							PostalCode:  "1000014",
							Prefecture:  "東京都",
							City:        "千代田区",
							CreatedAt:   1640962800,
							UpdatedAt:   1640962800,
						},
						{
							ID:            "producer-id02",
							Status:        entity.AdminStatusActivated,
							CoordinatorID: "coordinator-id",
							Lastname:      "&.",
							Firstname:     "管理者",
							LastnameKana:  "あんどどっと",
							FirstnameKana: "かんりしゃ",
							StoreName:     "&.農園",
							ThumbnailURL:  "https://and-period.jp/thumbnail.png",
							Thumbnails:    []*response.Image{},
							HeaderURL:     "https://and-period.jp/header.png",
							Headers:       []*response.Image{},
							Email:         "test-producer@and-period.jp",
							PhoneNumber:   "+819012345678",
							PostalCode:    "1000014",
							Prefecture:    "東京都",
							City:          "千代田区",
							CreatedAt:     1640962800,
							UpdatedAt:     1640962800,
						},
					},
					Total: 2,
				},
			},
		},
		{
			name: "success coordinator",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				in := &user.ListProducersInput{
					CoordinatorID: idmock,
					Limit:         20,
					Offset:        0,
				}
				mocks.user.EXPECT().ListProducers(gomock.Any(), in).Return(entity.Producers{}, int64(0), nil)
			},
			options: []testOption{withRole(uentity.AdminRoleCoordinator)},
			query:   "",
			expect: &testResponse{
				code: http.StatusOK,
				body: &response.ProducersResponse{
					Producers: []*response.Producer{},
					Total:     0,
				},
			},
		},
		{
			name: "success only unrelated producers",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				in := &user.ListProducersInput{
					Limit:         20,
					Offset:        0,
					OnlyUnrelated: true,
				}
				mocks.user.EXPECT().ListProducers(gomock.Any(), in).Return(entity.Producers{}, int64(0), nil)
			},
			options: []testOption{withRole(uentity.AdminRoleCoordinator)},
			query:   "?filters=unrelated",
			expect: &testResponse{
				code: http.StatusOK,
				body: &response.ProducersResponse{
					Producers: []*response.Producer{},
					Total:     0,
				},
			},
		},
		{
			name:  "invalid limit",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {},
			query: "?limit=a",
			expect: &testResponse{
				code: http.StatusBadRequest,
			},
		},
		{
			name:  "invalid offset",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {},
			query: "?offset=a",
			expect: &testResponse{
				code: http.StatusBadRequest,
			},
		},
		{
			name: "failed to get producers",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.user.EXPECT().ListProducers(gomock.Any(), in).Return(nil, int64(0), assert.AnError)
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
			const format = "/v1/producers%s"
			path := fmt.Sprintf(format, tt.query)
			testGet(t, tt.setup, tt.expect, path, tt.options...)
		})
	}
}

func TestGetProducer(t *testing.T) {
	t.Parallel()

	in := &user.GetProducerInput{
		ProducerID: "producer-id",
	}
	producer := &uentity.Producer{
		Admin: uentity.Admin{
			ID:            "producer-id",
			Role:          entity.AdminRoleProducer,
			Status:        entity.AdminStatusActivated,
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
		PhoneNumber: "+819012345678",
		PostalCode:  "1000014",
		Prefecture:  "東京都",
		City:        "千代田区",
		CreatedAt:   jst.Date(2022, 1, 1, 0, 0, 0, 0),
		UpdatedAt:   jst.Date(2022, 1, 1, 0, 0, 0, 0),
	}

	tests := []struct {
		name       string
		setup      func(t *testing.T, mocks *mocks, ctrl *gomock.Controller)
		producerID string
		expect     *testResponse
	}{
		{
			name: "success",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.user.EXPECT().GetProducer(gomock.Any(), in).Return(producer, nil)
			},
			producerID: "producer-id",
			expect: &testResponse{
				code: http.StatusOK,
				body: &response.ProducerResponse{
					Producer: &response.Producer{
						ID:            "producer-id",
						Status:        entity.AdminStatusActivated,
						CoordinatorID: "coordinator-id",
						Lastname:      "&.",
						Firstname:     "管理者",
						LastnameKana:  "あんどどっと",
						FirstnameKana: "かんりしゃ",
						StoreName:     "&.農園",
						ThumbnailURL:  "https://and-period.jp/thumbnail.png",
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
						Email:       "test-producer@and-period.jp",
						PhoneNumber: "+819012345678",
						PostalCode:  "1000014",
						Prefecture:  "東京都",
						City:        "千代田区",
						CreatedAt:   1640962800,
						UpdatedAt:   1640962800,
					},
				},
			},
		},
		{
			name: "failed to get producer",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.user.EXPECT().GetProducer(gomock.Any(), in).Return(nil, assert.AnError)
			},
			producerID: "producer-id",
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			const format = "/v1/producers/%s"
			path := fmt.Sprintf(format, tt.producerID)
			testGet(t, tt.setup, tt.expect, path)
		})
	}
}

func TestCreateProducer(t *testing.T) {
	t.Parallel()

	thumbnailIn := &media.UploadFileInput{
		URL: "https://tmp.and-period.jp/thumbnail.png",
	}
	thumbnailURL := "https://and-period.jp/thumbnail.png"
	headerIn := &media.UploadFileInput{
		URL: "https://tmp.and-period.jp/header.png",
	}
	headerURL := "https://and-period.jp/header.png"
	in := &user.CreateProducerInput{
		Lastname:      "&.",
		Firstname:     "生産者",
		LastnameKana:  "あんどどっと",
		FirstnameKana: "せいさんしゃ",
		StoreName:     "&.農園",
		ThumbnailURL:  thumbnailURL,
		HeaderURL:     headerURL,
		Email:         "test-producer@and-period.jp",
		PhoneNumber:   "+819012345678",
		PostalCode:    "1000014",
		Prefecture:    "東京都",
		City:          "千代田区",
		AddressLine1:  "永田町1-7-1",
		AddressLine2:  "",
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
		ThumbnailURL:  thumbnailURL,
		HeaderURL:     headerURL,
		PhoneNumber:   "+819012345678",
		PostalCode:    "1000014",
		Prefecture:    "東京都",
		City:          "千代田区",
		AddressLine1:  "永田町1-7-1",
		AddressLine2:  "",
		CreatedAt:     jst.Date(2022, 1, 1, 0, 0, 0, 0),
		UpdatedAt:     jst.Date(2022, 1, 1, 0, 0, 0, 0),
	}

	tests := []struct {
		name   string
		setup  func(t *testing.T, mocks *mocks, ctrl *gomock.Controller)
		req    *request.CreateProducerRequest
		expect *testResponse
	}{
		{
			name: "success",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.media.EXPECT().UploadProducerThumbnail(gomock.Any(), thumbnailIn).Return(thumbnailURL, nil)
				mocks.media.EXPECT().UploadProducerHeader(gomock.Any(), headerIn).Return(headerURL, nil)
				mocks.user.EXPECT().CreateProducer(gomock.Any(), in).Return(producer, nil)
			},
			req: &request.CreateProducerRequest{
				Lastname:      "&.",
				Firstname:     "生産者",
				LastnameKana:  "あんどどっと",
				FirstnameKana: "せいさんしゃ",
				StoreName:     "&.農園",
				ThumbnailURL:  "https://tmp.and-period.jp/thumbnail.png",
				HeaderURL:     "https://tmp.and-period.jp/header.png",
				Email:         "test-producer@and-period.jp",
				PhoneNumber:   "+819012345678",
				PostalCode:    "1000014",
				Prefecture:    "東京都",
				City:          "千代田区",
				AddressLine1:  "永田町1-7-1",
				AddressLine2:  "",
			},
			expect: &testResponse{
				code: http.StatusOK,
				body: &response.ProducerResponse{
					Producer: &response.Producer{
						ID:            "producer-id",
						CoordinatorID: "coordinator-id",
						Lastname:      "&.",
						Firstname:     "管理者",
						LastnameKana:  "あんどどっと",
						FirstnameKana: "かんりしゃ",
						StoreName:     "&.農園",
						ThumbnailURL:  "https://and-period.jp/thumbnail.png",
						Thumbnails:    []*response.Image{},
						HeaderURL:     "https://and-period.jp/header.png",
						Headers:       []*response.Image{},
						Email:         "test-producer@and-period.jp",
						PhoneNumber:   "+819012345678",
						PostalCode:    "1000014",
						Prefecture:    "東京都",
						City:          "千代田区",
						AddressLine1:  "永田町1-7-1",
						AddressLine2:  "",
						CreatedAt:     1640962800,
						UpdatedAt:     1640962800,
					},
				},
			},
		},
		{
			name: "failed to upload producer thumbnail",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.media.EXPECT().UploadProducerThumbnail(gomock.Any(), thumbnailIn).Return("", assert.AnError)
				mocks.media.EXPECT().UploadProducerHeader(gomock.Any(), headerIn).Return(headerURL, nil)
			},
			req: &request.CreateProducerRequest{
				Lastname:      "&.",
				Firstname:     "生産者",
				LastnameKana:  "あんどどっと",
				FirstnameKana: "せいさんしゃ",
				StoreName:     "&.農園",
				ThumbnailURL:  "https://tmp.and-period.jp/thumbnail.png",
				HeaderURL:     "https://tmp.and-period.jp/header.png",
				Email:         "test-producer@and-period.jp",
				PhoneNumber:   "+819012345678",
				PostalCode:    "1000014",
				Prefecture:    "東京都",
				City:          "千代田区",
				AddressLine1:  "永田町1-7-1",
				AddressLine2:  "",
			},
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
		{
			name: "failed to upload producer header",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.media.EXPECT().UploadProducerThumbnail(gomock.Any(), thumbnailIn).Return(thumbnailURL, nil)
				mocks.media.EXPECT().UploadProducerHeader(gomock.Any(), headerIn).Return("", assert.AnError)
			},
			req: &request.CreateProducerRequest{
				Lastname:      "&.",
				Firstname:     "生産者",
				LastnameKana:  "あんどどっと",
				FirstnameKana: "せいさんしゃ",
				StoreName:     "&.農園",
				ThumbnailURL:  "https://tmp.and-period.jp/thumbnail.png",
				HeaderURL:     "https://tmp.and-period.jp/header.png",
				Email:         "test-producer@and-period.jp",
				PhoneNumber:   "+819012345678",
				PostalCode:    "1000014",
				Prefecture:    "東京都",
				City:          "千代田区",
				AddressLine1:  "永田町1-7-1",
				AddressLine2:  "",
			},
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
		{
			name: "failed to create producer",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				in := *in
				in.ThumbnailURL = ""
				in.HeaderURL = ""
				mocks.user.EXPECT().CreateProducer(gomock.Any(), &in).Return(nil, assert.AnError)
			},
			req: &request.CreateProducerRequest{
				Lastname:      "&.",
				Firstname:     "生産者",
				LastnameKana:  "あんどどっと",
				FirstnameKana: "せいさんしゃ",
				StoreName:     "&.農園",
				ThumbnailURL:  "",
				HeaderURL:     "",
				Email:         "test-producer@and-period.jp",
				PhoneNumber:   "+819012345678",
				PostalCode:    "1000014",
				Prefecture:    "東京都",
				City:          "千代田区",
				AddressLine1:  "永田町1-7-1",
				AddressLine2:  "",
			},
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			const path = "/v1/producers"
			testPost(t, tt.setup, tt.expect, path, tt.req)
		})
	}
}

func TestUpdateProducer(t *testing.T) {
	t.Parallel()

	thumbnailIn := &media.UploadFileInput{
		URL: "https://tmp.and-period.jp/thumbnail.png",
	}
	thumbnailURL := "https://and-period.jp/thumbnail.png"
	headerIn := &media.UploadFileInput{
		URL: "https://tmp.and-period.jp/header.png",
	}
	headerURL := "https://and-period.jp/header.png"
	in := &user.UpdateProducerInput{
		ProducerID:    "producer-id",
		Lastname:      "&.",
		Firstname:     "生産者",
		LastnameKana:  "あんどどっと",
		FirstnameKana: "せいさんしゃ",
		StoreName:     "&.農園",
		ThumbnailURL:  "https://and-period.jp/thumbnail.png",
		HeaderURL:     "https://and-period.jp/header.png",
		PhoneNumber:   "+819012345678",
		PostalCode:    "1000014",
		Prefecture:    "東京都",
		City:          "千代田区",
		AddressLine1:  "永田町1-7-1",
		AddressLine2:  "",
	}

	tests := []struct {
		name       string
		setup      func(t *testing.T, mocks *mocks, ctrl *gomock.Controller)
		producerID string
		req        *request.UpdateProducerRequest
		expect     *testResponse
	}{
		{
			name: "success",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.media.EXPECT().UploadProducerThumbnail(gomock.Any(), thumbnailIn).Return(thumbnailURL, nil)
				mocks.media.EXPECT().UploadProducerHeader(gomock.Any(), headerIn).Return(headerURL, nil)
				mocks.user.EXPECT().UpdateProducer(gomock.Any(), in).Return(nil)
			},
			producerID: "producer-id",
			req: &request.UpdateProducerRequest{
				Lastname:      "&.",
				Firstname:     "生産者",
				LastnameKana:  "あんどどっと",
				FirstnameKana: "せいさんしゃ",
				StoreName:     "&.農園",
				ThumbnailURL:  "https://tmp.and-period.jp/thumbnail.png",
				HeaderURL:     "https://tmp.and-period.jp/header.png",
				PhoneNumber:   "+819012345678",
				PostalCode:    "1000014",
				Prefecture:    "東京都",
				City:          "千代田区",
				AddressLine1:  "永田町1-7-1",
				AddressLine2:  "",
			},
			expect: &testResponse{
				code: http.StatusNoContent,
			},
		},
		{
			name: "failed to upload producer thumbnail",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.media.EXPECT().UploadProducerThumbnail(gomock.Any(), thumbnailIn).Return("", assert.AnError)
				mocks.media.EXPECT().UploadProducerHeader(gomock.Any(), headerIn).Return(headerURL, nil)
			},
			producerID: "producer-id",
			req: &request.UpdateProducerRequest{
				Lastname:      "&.",
				Firstname:     "生産者",
				LastnameKana:  "あんどどっと",
				FirstnameKana: "せいさんしゃ",
				StoreName:     "&.農園",
				ThumbnailURL:  "https://tmp.and-period.jp/thumbnail.png",
				HeaderURL:     "https://tmp.and-period.jp/header.png",
				PhoneNumber:   "+819012345678",
				PostalCode:    "1000014",
				Prefecture:    "東京都",
				City:          "千代田区",
				AddressLine1:  "永田町1-7-1",
				AddressLine2:  "",
			},
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
		{
			name: "failed to upload producer header",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.media.EXPECT().UploadProducerThumbnail(gomock.Any(), thumbnailIn).Return(thumbnailURL, nil)
				mocks.media.EXPECT().UploadProducerHeader(gomock.Any(), headerIn).Return("", assert.AnError)
			},
			producerID: "producer-id",
			req: &request.UpdateProducerRequest{
				Lastname:      "&.",
				Firstname:     "生産者",
				LastnameKana:  "あんどどっと",
				FirstnameKana: "せいさんしゃ",
				StoreName:     "&.農園",
				ThumbnailURL:  "https://tmp.and-period.jp/thumbnail.png",
				HeaderURL:     "https://tmp.and-period.jp/header.png",
				PhoneNumber:   "+819012345678",
				PostalCode:    "1000014",
				Prefecture:    "東京都",
				City:          "千代田区",
				AddressLine1:  "永田町1-7-1",
				AddressLine2:  "",
			},
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
		{
			name: "failed to update producer",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				in := *in
				in.ThumbnailURL = ""
				in.HeaderURL = ""
				mocks.user.EXPECT().UpdateProducer(gomock.Any(), &in).Return(assert.AnError)
			},
			producerID: "producer-id",
			req: &request.UpdateProducerRequest{
				Lastname:      "&.",
				Firstname:     "生産者",
				LastnameKana:  "あんどどっと",
				FirstnameKana: "せいさんしゃ",
				StoreName:     "&.農園",
				ThumbnailURL:  "",
				HeaderURL:     "",
				PhoneNumber:   "+819012345678",
				PostalCode:    "1000014",
				Prefecture:    "東京都",
				City:          "千代田区",
				AddressLine1:  "永田町1-7-1",
				AddressLine2:  "",
			},
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			const format = "/v1/producers/%s"
			path := fmt.Sprintf(format, tt.producerID)
			testPatch(t, tt.setup, tt.expect, path, tt.req)
		})
	}
}

func TestUpdateProducerEmail(t *testing.T) {
	t.Parallel()

	in := &user.UpdateProducerEmailInput{
		ProducerID: "producer-id",
		Email:      "test-producer@and-period.jp",
	}

	tests := []struct {
		name       string
		setup      func(t *testing.T, mocks *mocks, ctrl *gomock.Controller)
		producerID string
		req        *request.UpdateProducerEmailRequest
		expect     *testResponse
	}{
		{
			name: "success",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.user.EXPECT().UpdateProducerEmail(gomock.Any(), in).Return(nil)
			},
			producerID: "producer-id",
			req: &request.UpdateProducerEmailRequest{
				Email: "test-producer@and-period.jp",
			},
			expect: &testResponse{
				code: http.StatusNoContent,
			},
		},
		{
			name: "failed to update producer email",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.user.EXPECT().UpdateProducerEmail(gomock.Any(), in).Return(assert.AnError)
			},
			producerID: "producer-id",
			req: &request.UpdateProducerEmailRequest{
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
			const format = "/v1/producers/%s/email"
			path := fmt.Sprintf(format, tt.producerID)
			testPatch(t, tt.setup, tt.expect, path, tt.req)
		})
	}
}

func TestResetProducerPassword(t *testing.T) {
	t.Parallel()

	in := &user.ResetProducerPasswordInput{
		ProducerID: "producer-id",
	}

	tests := []struct {
		name       string
		setup      func(t *testing.T, mocks *mocks, ctrl *gomock.Controller)
		producerID string
		expect     *testResponse
	}{
		{
			name: "success",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.user.EXPECT().ResetProducerPassword(gomock.Any(), in).Return(nil)
			},
			producerID: "producer-id",
			expect: &testResponse{
				code: http.StatusNoContent,
			},
		},
		{
			name: "failed to reset producer password",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.user.EXPECT().ResetProducerPassword(gomock.Any(), in).Return(assert.AnError)
			},
			producerID: "producer-id",
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			const format = "/v1/producers/%s/password"
			path := fmt.Sprintf(format, tt.producerID)
			testPatch(t, tt.setup, tt.expect, path, nil)
		})
	}
}

func TestDeleteProducer(t *testing.T) {
	t.Parallel()

	in := &user.DeleteProducerInput{
		ProducerID: "producer-id",
	}

	tests := []struct {
		name       string
		setup      func(t *testing.T, mocks *mocks, ctrl *gomock.Controller)
		producerID string
		expect     *testResponse
	}{
		{
			name: "success",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.user.EXPECT().DeleteProducer(gomock.Any(), in).Return(nil)
			},
			producerID: "producer-id",
			expect: &testResponse{
				code: http.StatusNoContent,
			},
		},
		{
			name: "failed to delete producer",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.user.EXPECT().DeleteProducer(gomock.Any(), in).Return(assert.AnError)
			},
			producerID: "producer-id",
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			const format = "/v1/producers/%s"
			path := fmt.Sprintf(format, tt.producerID)
			testDelete(t, tt.setup, tt.expect, path)
		})
	}
}
