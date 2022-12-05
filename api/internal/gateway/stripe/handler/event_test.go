package handler

import (
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stripe/stripe-go/v73"
)

func TestEvent(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		setup  func(t *testing.T, mocks *mocks, ctrl *gomock.Controller)
		req    *stripe.Event
		expect *testResponse
	}{
		{
			name: "success to payment failed",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				event := &stripe.Event{Type: "payment_intent.payment_failed"}
				mocks.receiver.EXPECT().Receive(gomock.Any(), signaturemock).Return(event, nil)
			},
			req: &stripe.Event{},
			expect: &testResponse{
				code: http.StatusNoContent,
			},
		},
		{
			name: "not implemented event",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				event := &stripe.Event{Type: "unknown"}
				mocks.receiver.EXPECT().Receive(gomock.Any(), signaturemock).Return(event, nil)
			},
			req: &stripe.Event{},
			expect: &testResponse{
				code: http.StatusNoContent,
			},
		},
		{
			name: "invalid request",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.receiver.EXPECT().Receive(gomock.Any(), signaturemock).Return(nil, assert.AnError)
			},
			req: &stripe.Event{},
			expect: &testResponse{
				code: http.StatusBadRequest,
			},
		},
		{
			name: "unknown dispatch event",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.receiver.EXPECT().Receive(gomock.Any(), signaturemock).Return(nil, nil)
			},
			req: &stripe.Event{},
			expect: &testResponse{
				code: http.StatusBadRequest,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			const path = "/stripe/webhooks"
			testPost(t, tt.setup, tt.expect, path, tt.req)
		})
	}
}
