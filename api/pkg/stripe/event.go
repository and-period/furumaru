package stripe

import (
	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/webhook"
)

// reference: https://stripe.com/docs/webhooks
func (r *receiver) Receive(payload []byte, signature string) (*stripe.Event, error) {
	event, err := webhook.ConstructEvent(payload, signature, r.webhookKey)
	if err != nil {
		return nil, err
	}
	return &event, nil
}
