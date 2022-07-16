package handler

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/request"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/golang/mock/gomock"
)

func TestListShippings(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		setup  func(t *testing.T, mocks *mocks, ctrl *gomock.Controller)
		query  string
		expect *testResponse
	}{
		{
			name:  "success",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {},
			query: "",
			expect: &testResponse{
				code: http.StatusOK,
				body: &response.ShippingsResponse{
					Shippings: []*response.Shipping{},
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			const prefix = "/v1/shippings"
			path := fmt.Sprintf("%s%s", prefix, tt.query)
			testGet(t, tt.setup, tt.expect, path)
		})
	}
}

func TestGetShipping(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		setup      func(t *testing.T, mocks *mocks, ctrl *gomock.Controller)
		shippingID string
		expect     *testResponse
	}{
		{
			name:       "success",
			setup:      func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {},
			shippingID: "shipping-id",
			expect: &testResponse{
				code: http.StatusOK,
				body: &response.ShippingResponse{
					Shipping: &response.Shipping{},
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			const prefix = "/v1/shippings"
			path := fmt.Sprintf("%s/%s", prefix, tt.shippingID)
			testGet(t, tt.setup, tt.expect, path)
		})
	}
}

func TestCreateShipping(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		setup  func(t *testing.T, mocks *mocks, ctrl *gomock.Controller)
		req    *request.CreateShippingRequest
		expect *testResponse
	}{
		{
			name:  "success",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {},
			req:   &request.CreateShippingRequest{},
			expect: &testResponse{
				code: http.StatusOK,
				body: &response.ShippingResponse{
					Shipping: &response.Shipping{},
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			const path = "/v1/shippings"
			testPost(t, tt.setup, tt.expect, path, tt.req)
		})
	}
}

func TestUpdateShipping(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		setup      func(t *testing.T, mocks *mocks, ctrl *gomock.Controller)
		shippingID string
		req        *request.UpdateShippingRequest
		expect     *testResponse
	}{
		{
			name:       "success",
			setup:      func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {},
			shippingID: "shipping-id",
			req:        &request.UpdateShippingRequest{},
			expect: &testResponse{
				code: http.StatusNoContent,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			const prefix = "/v1/shippings"
			path := fmt.Sprintf("%s/%s", prefix, tt.shippingID)
			testPatch(t, tt.setup, tt.expect, path, tt.req)
		})
	}
}

func TestDeleteShipping(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		setup      func(t *testing.T, mocks *mocks, ctrl *gomock.Controller)
		shippingID string
		expect     *testResponse
	}{
		{
			name:       "success",
			setup:      func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {},
			shippingID: "shipping-id",
			expect: &testResponse{
				code: http.StatusNoContent,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			const prefix = "/v1/shippings"
			path := fmt.Sprintf("%s/%s", prefix, tt.shippingID)
			testDelete(t, tt.setup, tt.expect, path)
		})
	}
}
