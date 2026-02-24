package handler

import (
	"encoding/json"
	"net/http/httptest"
	"sync"
	"testing"

	mock_store "github.com/and-period/furumaru/api/mock/store"
	mock_stripe "github.com/and-period/furumaru/api/mock/pkg/stripe"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	lib "github.com/stripe/stripe-go/v82"
	"go.uber.org/mock/gomock"
)

type mocks struct {
	store    *mock_store.MockService
	receiver *mock_stripe.MockReceiver
}

func newMocks(ctrl *gomock.Controller) *mocks {
	return &mocks{
		store:    mock_store.NewMockService(ctrl),
		receiver: mock_stripe.NewMockReceiver(ctrl),
	}
}

func newHandler(m *mocks) *handler {
	params := &Params{
		WaitGroup: &sync.WaitGroup{},
		Store:     m.store,
		Receiver:  m.receiver,
	}
	return NewHandler(params).(*handler)
}

func testGinContext(t *testing.T, w *httptest.ResponseRecorder) (*gin.Context, *gin.Engine) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	ctx, engine := gin.CreateTestContext(w)
	return ctx, engine
}

func TestNewHandler(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	m := newMocks(ctrl)
	h := newHandler(m)
	assert.NotNil(t, h)
}

func TestRoutes(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	m := newMocks(ctrl)
	h := newHandler(m)
	rg := gin.New().Group("")
	h.Routes(rg)
}

func newPaymentIntentEvent(t *testing.T, eventType string, pi *lib.PaymentIntent) *lib.Event {
	t.Helper()
	raw, err := json.Marshal(pi)
	if err != nil {
		t.Fatal(err)
	}
	return &lib.Event{
		Type:    lib.EventType(eventType),
		Created: 1700000000,
		Data:    &lib.EventData{Raw: raw},
	}
}

func newChargeEvent(t *testing.T, charge *lib.Charge) *lib.Event {
	t.Helper()
	raw, err := json.Marshal(charge)
	if err != nil {
		t.Fatal(err)
	}
	return &lib.Event{
		Type:    lib.EventType(eventTypeChargeRefunded),
		Created: 1700000000,
		Data:    &lib.EventData{Raw: raw},
	}
}
