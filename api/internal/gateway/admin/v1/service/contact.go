package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/messenger/entity"
)

type Contact struct {
	response.Contact
}

func NewContact(contact *entity.Contact) *Contact {
	return &Contact{
		Contact: response.Contact{
			ID:          contact.ID,
			Title:       contact.Title,
			CategoryID:  contact.CategoryID,
			Content:     contact.Content,
			Username:    contact.Username,
			UserID:      contact.UserID,
			Email:       contact.Username,
			PhoneNumber: contact.PhoneNumber,
			Status:      response.ContactStatus(contact.Status),
			ResponderID: contact.ResponderID,
			CreatedAt:   contact.CreatedAt.Unix(),
			UpdatedAt:   contact.UpdatedAt.Unix(),
		},
	}
}

func (c *Contact) Response() *response.Contact {
	return &c.Contact
}
