package stripe

import (
	"github.com/stripe/stripe-go/v87"
	"github.com/stripe/stripe-go/v87/webhook"
)

// reference: https://stripe.com/docs/webhooks
func (r *receiver) Receive(payload []byte, signature string) (*stripe.Event, error) {
	event, err := webhook.ConstructEventWithOptions(payload, signature, r.webhookKey, webhook.ConstructEventOptions{
		IgnoreAPIVersionMismatch: true,
	})
	if err != nil {
		return nil, err
	}
	return &event, nil
}
