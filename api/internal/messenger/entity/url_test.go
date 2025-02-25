package entity

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAdminURLMaker(t *testing.T) {
	t.Parallel()
	webURL, err := url.Parse("http://example.com")
	require.NoError(t, err)
	maker := NewAdminURLMaker(webURL)
	res := maker.SignIn()
	assert.Equal(t, "http://example.com/signin", res)
	res = maker.Contact("contact-id")
	assert.Equal(t, "http://example.com/contacts/contact-id", res)
	res = maker.Notification("notification-id")
	assert.Equal(t, "http://example.com/notifications/notification-id", res)
	res = maker.Order("order-id")
	assert.Equal(t, "http://example.com/orders/order-id", res)
}

func TestUserURLMaker(t *testing.T) {
	t.Parallel()
	webURL, err := url.Parse("http://example.com")
	require.NoError(t, err)
	maker := NewUserURLMaker(webURL)
	res := maker.SignIn()
	assert.Equal(t, "http://example.com/signin", res)
	res = maker.Live("schedule-id")
	assert.Equal(t, "http://example.com/live/schedule-id", res)
	res = maker.ProductReview("product-id")
	assert.Equal(t, "http://example.com/reviews/products/product-id", res)
	res = maker.ExperienceReview("experience-id")
	assert.Equal(t, "http://example.com/reviews/experiences/experience-id", res)
}
