package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/messenger/entity"
)

type Contact struct {
	response.Contact
}

type Contacts []*Contact

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

func NewContacts(contacts entity.Contacts) Contacts {
	res := make(Contacts, len(contacts))
	for i := range contacts {
		res[i] = NewContact(contacts[i])
	}
	return res
}

func (cs Contacts) Response() []*response.Contact {
	res := make([]*response.Contact, len(cs))
	for i := range cs {
		res[i] = cs[i].Response()
	}
	return res
}
