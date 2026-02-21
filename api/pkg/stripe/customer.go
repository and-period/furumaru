package stripe

import (
	"context"

	"github.com/stripe/stripe-go/v82"
)

type CreateCustomerParams struct {
	UserID      string
	Name        string
	Email       string
	Description string
	PhoneNumber string
}

// reference: https://stripe.com/docs/api/customers/retrieve
func (c *client) GetCustomer(ctx context.Context, customerID string) (*stripe.Customer, error) {
	params := &stripe.CustomerParams{
		Params: stripe.Params{Context: ctx},
	}
	return c.customer.Get(customerID, params)
}

// reference: https://stripe.com/docs/api/customers/create
func (c *client) CreateCustomer(ctx context.Context, in *CreateCustomerParams) (*stripe.Customer, error) {
	params := &stripe.CustomerParams{
		Params: stripe.Params{
			Context:  ctx,
			Metadata: map[string]string{"userId": in.UserID},
		},
		Name:        stripe.String(in.Name),
		Email:       stripe.String(in.Email),
		Phone:       nullString(in.PhoneNumber),
		Description: nullString(in.Description),
	}
	return c.customer.New(params)
}

// reference: https://stripe.com/docs/api/customers/delete
func (c *client) DeleteCustomer(ctx context.Context, customerID string) error {
	params := &stripe.CustomerParams{
		Params: stripe.Params{Context: ctx},
	}
	_, err := c.customer.Del(customerID, params)
	return err
}
