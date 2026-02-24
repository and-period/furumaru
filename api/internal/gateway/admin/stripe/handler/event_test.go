package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	lib "github.com/stripe/stripe-go/v82"
	"go.uber.org/mock/gomock"
)

func TestEvent_InvalidBody(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	m := newMocks(ctrl)
	h := newHandler(m)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/stripe/webhooks", &errorReader{})
	req.Header.Set("Content-Type", "application/json")
	ctx, _ := testGinContext(t, w)
	ctx.Request = req

	h.Event(ctx)
	assert.Equal(t, http.StatusBadRequest, ctx.Writer.Status())
}

func TestEvent_InvalidSignature(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	m := newMocks(ctrl)
	h := newHandler(m)

	body := []byte(`{"id": "evt_test"}`)
	m.receiver.EXPECT().
		Receive(body, "invalid-sig").
		Return(nil, assert.AnError)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/stripe/webhooks", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Stripe-Signature", "invalid-sig")
	ctx, _ := testGinContext(t, w)
	ctx.Request = req

	h.Event(ctx)
	assert.Equal(t, http.StatusUnauthorized, ctx.Writer.Status())
}

func TestEvent_UnexpectedEvent(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	m := newMocks(ctrl)
	h := newHandler(m)

	event := &lib.Event{
		Type:    "unknown.event",
		Created: 1700000000,
		Data:    &lib.EventData{Raw: []byte(`{}`)},
	}
	body := []byte(`{"type": "unknown.event"}`)
	m.receiver.EXPECT().
		Receive(body, "valid-sig").
		Return(event, nil)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/stripe/webhooks", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Stripe-Signature", "valid-sig")
	ctx, _ := testGinContext(t, w)
	ctx.Request = req

	h.Event(ctx)
	assert.Equal(t, http.StatusNoContent, ctx.Writer.Status())
}

// errorReader implements io.Reader that always returns an error.
type errorReader struct{}

func (r *errorReader) Read(_ []byte) (int, error) {
	return 0, assert.AnError
}
