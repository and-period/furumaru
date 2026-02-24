package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/stretchr/testify/assert"
	lib "github.com/stripe/stripe-go/v82"
	"go.uber.org/mock/gomock"
)

func TestPaymentIntentAuthorized(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		setup  func(m *mocks)
		event  *lib.Event
		expect int
	}{
		{
			name: "success",
			setup: func(m *mocks) {
				m.store.EXPECT().
					NotifyPaymentAuthorized(gomock.Any(), &store.NotifyPaymentAuthorizedInput{
						NotifyPaymentPayload: store.NotifyPaymentPayload{
							OrderID:   "order-id",
							PaymentID: "pi_test",
							IssuedAt:  time.Unix(1700000000, 0),
							Status:    entity.PaymentStatusAuthorized,
						},
					}).
					Return(nil)
			},
			event: newPaymentIntentEvent(t, eventTypePaymentIntentAmountCapturableUpdated, &lib.PaymentIntent{
				ID:       "pi_test",
				Metadata: map[string]string{"order_id": "order-id"},
			}),
			expect: http.StatusNoContent,
		},
		{
			name:  "invalid event data",
			setup: func(_ *mocks) {},
			event: &lib.Event{
				Type:    lib.EventType(eventTypePaymentIntentAmountCapturableUpdated),
				Created: 1700000000,
				Data:    &lib.EventData{Raw: []byte(`invalid-json`)},
			},
			expect: http.StatusBadRequest,
		},
		{
			name: "store error",
			setup: func(m *mocks) {
				m.store.EXPECT().
					NotifyPaymentAuthorized(gomock.Any(), gomock.Any()).
					Return(assert.AnError)
			},
			event: newPaymentIntentEvent(t, eventTypePaymentIntentAmountCapturableUpdated, &lib.PaymentIntent{
				ID:       "pi_test",
				Metadata: map[string]string{"order_id": "order-id"},
			}),
			expect: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			m := newMocks(ctrl)
			h := newHandler(m)
			tt.setup(m)

			w := httptest.NewRecorder()
			ctx, _ := testGinContext(t, w)
			ctx.Request = httptest.NewRequest(http.MethodPost, "/", nil)
			h.paymentIntentAuthorized(ctx, tt.event)
			assert.Equal(t, tt.expect, ctx.Writer.Status())
		})
	}
}

func TestPaymentIntentSucceeded(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		setup  func(m *mocks)
		event  *lib.Event
		expect int
	}{
		{
			name: "success",
			setup: func(m *mocks) {
				m.store.EXPECT().
					NotifyPaymentCaptured(gomock.Any(), &store.NotifyPaymentCapturedInput{
						NotifyPaymentPayload: store.NotifyPaymentPayload{
							OrderID:   "order-id",
							PaymentID: "pi_test",
							IssuedAt:  time.Unix(1700000000, 0),
							Status:    entity.PaymentStatusCaptured,
						},
					}).
					Return(nil)
			},
			event: newPaymentIntentEvent(t, eventTypePaymentIntentSucceeded, &lib.PaymentIntent{
				ID:       "pi_test",
				Metadata: map[string]string{"order_id": "order-id"},
			}),
			expect: http.StatusNoContent,
		},
		{
			name:  "invalid event data",
			setup: func(_ *mocks) {},
			event: &lib.Event{
				Type:    lib.EventType(eventTypePaymentIntentSucceeded),
				Created: 1700000000,
				Data:    &lib.EventData{Raw: []byte(`invalid-json`)},
			},
			expect: http.StatusBadRequest,
		},
		{
			name: "store error",
			setup: func(m *mocks) {
				m.store.EXPECT().
					NotifyPaymentCaptured(gomock.Any(), gomock.Any()).
					Return(assert.AnError)
			},
			event: newPaymentIntentEvent(t, eventTypePaymentIntentSucceeded, &lib.PaymentIntent{
				ID:       "pi_test",
				Metadata: map[string]string{"order_id": "order-id"},
			}),
			expect: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			m := newMocks(ctrl)
			h := newHandler(m)
			tt.setup(m)

			w := httptest.NewRecorder()
			ctx, _ := testGinContext(t, w)
			ctx.Request = httptest.NewRequest(http.MethodPost, "/", nil)
			h.paymentIntentSucceeded(ctx, tt.event)
			assert.Equal(t, tt.expect, ctx.Writer.Status())
		})
	}
}

func TestPaymentIntentFailed(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		setup  func(m *mocks)
		event  *lib.Event
		expect int
	}{
		{
			name: "success",
			setup: func(m *mocks) {
				m.store.EXPECT().
					NotifyPaymentFailed(gomock.Any(), &store.NotifyPaymentFailedInput{
						NotifyPaymentPayload: store.NotifyPaymentPayload{
							OrderID:   "order-id",
							PaymentID: "pi_test",
							IssuedAt:  time.Unix(1700000000, 0),
							Status:    entity.PaymentStatusFailed,
						},
					}).
					Return(nil)
			},
			event: newPaymentIntentEvent(t, eventTypePaymentIntentPaymentFailed, &lib.PaymentIntent{
				ID:       "pi_test",
				Metadata: map[string]string{"order_id": "order-id"},
			}),
			expect: http.StatusNoContent,
		},
		{
			name:  "invalid event data",
			setup: func(_ *mocks) {},
			event: &lib.Event{
				Type:    lib.EventType(eventTypePaymentIntentPaymentFailed),
				Created: 1700000000,
				Data:    &lib.EventData{Raw: []byte(`invalid-json`)},
			},
			expect: http.StatusBadRequest,
		},
		{
			name: "store error",
			setup: func(m *mocks) {
				m.store.EXPECT().
					NotifyPaymentFailed(gomock.Any(), gomock.Any()).
					Return(assert.AnError)
			},
			event: newPaymentIntentEvent(t, eventTypePaymentIntentPaymentFailed, &lib.PaymentIntent{
				ID:       "pi_test",
				Metadata: map[string]string{"order_id": "order-id"},
			}),
			expect: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			m := newMocks(ctrl)
			h := newHandler(m)
			tt.setup(m)

			w := httptest.NewRecorder()
			ctx, _ := testGinContext(t, w)
			ctx.Request = httptest.NewRequest(http.MethodPost, "/", nil)
			h.paymentIntentFailed(ctx, tt.event)
			assert.Equal(t, tt.expect, ctx.Writer.Status())
		})
	}
}

func TestPaymentIntentCanceled(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		setup  func(m *mocks)
		event  *lib.Event
		expect int
	}{
		{
			name: "success",
			setup: func(m *mocks) {
				m.store.EXPECT().
					NotifyPaymentRefunded(gomock.Any(), &store.NotifyPaymentRefundedInput{
						NotifyPaymentPayload: store.NotifyPaymentPayload{
							OrderID:   "order-id",
							PaymentID: "pi_test",
							IssuedAt:  time.Unix(1700000000, 0),
							Status:    entity.PaymentStatusCanceled,
						},
						Type:   entity.RefundTypeCanceled,
						Reason: string(lib.PaymentIntentCancellationReasonRequestedByCustomer),
					}).
					Return(nil)
			},
			event: newPaymentIntentEvent(t, eventTypePaymentIntentCanceled, &lib.PaymentIntent{
				ID:                 "pi_test",
				Metadata:           map[string]string{"order_id": "order-id"},
				CancellationReason: lib.PaymentIntentCancellationReasonRequestedByCustomer,
			}),
			expect: http.StatusNoContent,
		},
		{
			name: "success without cancellation reason",
			setup: func(m *mocks) {
				m.store.EXPECT().
					NotifyPaymentRefunded(gomock.Any(), &store.NotifyPaymentRefundedInput{
						NotifyPaymentPayload: store.NotifyPaymentPayload{
							OrderID:   "order-id",
							PaymentID: "pi_test",
							IssuedAt:  time.Unix(1700000000, 0),
							Status:    entity.PaymentStatusCanceled,
						},
						Type: entity.RefundTypeCanceled,
					}).
					Return(nil)
			},
			event: newPaymentIntentEvent(t, eventTypePaymentIntentCanceled, &lib.PaymentIntent{
				ID:       "pi_test",
				Metadata: map[string]string{"order_id": "order-id"},
			}),
			expect: http.StatusNoContent,
		},
		{
			name:  "invalid event data",
			setup: func(_ *mocks) {},
			event: &lib.Event{
				Type:    lib.EventType(eventTypePaymentIntentCanceled),
				Created: 1700000000,
				Data:    &lib.EventData{Raw: []byte(`invalid-json`)},
			},
			expect: http.StatusBadRequest,
		},
		{
			name: "store error",
			setup: func(m *mocks) {
				m.store.EXPECT().
					NotifyPaymentRefunded(gomock.Any(), gomock.Any()).
					Return(assert.AnError)
			},
			event: newPaymentIntentEvent(t, eventTypePaymentIntentCanceled, &lib.PaymentIntent{
				ID:       "pi_test",
				Metadata: map[string]string{"order_id": "order-id"},
			}),
			expect: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			m := newMocks(ctrl)
			h := newHandler(m)
			tt.setup(m)

			w := httptest.NewRecorder()
			ctx, _ := testGinContext(t, w)
			ctx.Request = httptest.NewRequest(http.MethodPost, "/", nil)
			h.paymentIntentCanceled(ctx, tt.event)
			assert.Equal(t, tt.expect, ctx.Writer.Status())
		})
	}
}

func TestChargeRefunded(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		setup  func(m *mocks)
		event  *lib.Event
		expect int
	}{
		{
			name: "success",
			setup: func(m *mocks) {
				m.store.EXPECT().
					NotifyPaymentRefunded(gomock.Any(), &store.NotifyPaymentRefundedInput{
						NotifyPaymentPayload: store.NotifyPaymentPayload{
							OrderID:   "order-id",
							PaymentID: "pi_test",
							IssuedAt:  time.Unix(1700000000, 0),
							Status:    entity.PaymentStatusRefunded,
						},
						Type:   entity.RefundTypeRefunded,
						Reason: string(lib.RefundReasonRequestedByCustomer),
						Total:  1000,
					}).
					Return(nil)
			},
			event: newChargeEvent(t, &lib.Charge{
				ID:       "ch_test",
				Metadata: map[string]string{"order_id": "order-id"},
				PaymentIntent: &lib.PaymentIntent{
					ID: "pi_test",
				},
				Refunds: &lib.RefundList{
					Data: []*lib.Refund{
						{
							Amount: 1000,
							Reason: lib.RefundReasonRequestedByCustomer,
						},
					},
				},
			}),
			expect: http.StatusNoContent,
		},
		{
			name: "success with multiple refunds",
			setup: func(m *mocks) {
				m.store.EXPECT().
					NotifyPaymentRefunded(gomock.Any(), &store.NotifyPaymentRefundedInput{
						NotifyPaymentPayload: store.NotifyPaymentPayload{
							OrderID:   "order-id",
							PaymentID: "pi_test",
							IssuedAt:  time.Unix(1700000000, 0),
							Status:    entity.PaymentStatusRefunded,
						},
						Type:   entity.RefundTypeRefunded,
						Reason: string(lib.RefundReasonRequestedByCustomer) + "\n" + string(lib.RefundReasonDuplicate),
						Total:  1500,
					}).
					Return(nil)
			},
			event: newChargeEvent(t, &lib.Charge{
				ID:       "ch_test",
				Metadata: map[string]string{"order_id": "order-id"},
				PaymentIntent: &lib.PaymentIntent{
					ID: "pi_test",
				},
				Refunds: &lib.RefundList{
					Data: []*lib.Refund{
						{
							Amount: 1000,
							Reason: lib.RefundReasonRequestedByCustomer,
						},
						{
							Amount: 500,
							Reason: lib.RefundReasonDuplicate,
						},
					},
				},
			}),
			expect: http.StatusNoContent,
		},
		{
			name: "success without refund details",
			setup: func(m *mocks) {
				m.store.EXPECT().
					NotifyPaymentRefunded(gomock.Any(), &store.NotifyPaymentRefundedInput{
						NotifyPaymentPayload: store.NotifyPaymentPayload{
							OrderID:   "order-id",
							PaymentID: "pi_test",
							IssuedAt:  time.Unix(1700000000, 0),
							Status:    entity.PaymentStatusRefunded,
						},
						Type: entity.RefundTypeRefunded,
					}).
					Return(nil)
			},
			event: newChargeEvent(t, &lib.Charge{
				ID:       "ch_test",
				Metadata: map[string]string{"order_id": "order-id"},
				PaymentIntent: &lib.PaymentIntent{
					ID: "pi_test",
				},
			}),
			expect: http.StatusNoContent,
		},
		{
			name:  "invalid event data",
			setup: func(_ *mocks) {},
			event: &lib.Event{
				Type:    lib.EventType(eventTypeChargeRefunded),
				Created: 1700000000,
				Data:    &lib.EventData{Raw: []byte(`invalid-json`)},
			},
			expect: http.StatusBadRequest,
		},
		{
			name: "store error",
			setup: func(m *mocks) {
				m.store.EXPECT().
					NotifyPaymentRefunded(gomock.Any(), gomock.Any()).
					Return(assert.AnError)
			},
			event: newChargeEvent(t, &lib.Charge{
				ID:       "ch_test",
				Metadata: map[string]string{"order_id": "order-id"},
				PaymentIntent: &lib.PaymentIntent{
					ID: "pi_test",
				},
			}),
			expect: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			m := newMocks(ctrl)
			h := newHandler(m)
			tt.setup(m)

			w := httptest.NewRecorder()
			ctx, _ := testGinContext(t, w)
			ctx.Request = httptest.NewRequest(http.MethodPost, "/", nil)
			h.chargeRefunded(ctx, tt.event)
			assert.Equal(t, tt.expect, ctx.Writer.Status())
		})
	}
}
