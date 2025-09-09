package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/types"
	"github.com/and-period/furumaru/api/internal/messenger/entity"
)

type Contact struct {
	types.Contact
}

type Contacts []*Contact

func NewContact(contact *entity.Contact) *Contact {
	return &Contact{
		Contact: types.Contact{
			ID:          contact.ID,
			Title:       contact.Title,
			CategoryID:  contact.CategoryID,
			Content:     contact.Content,
			Username:    contact.Username,
			UserID:      contact.UserID,
			Email:       contact.Username,
			PhoneNumber: contact.PhoneNumber,
			Status:      types.ContactStatus(contact.Status),
			ResponderID: contact.ResponderID,
			CreatedAt:   contact.CreatedAt.Unix(),
			UpdatedAt:   contact.UpdatedAt.Unix(),
		},
	}
}

func (c *Contact) Response() *types.Contact {
	return &c.Contact
}

func NewContacts(contacts entity.Contacts) Contacts {
	res := make(Contacts, len(contacts))
	for i := range contacts {
		res[i] = NewContact(contacts[i])
	}
	return res
}

func (cs Contacts) Response() []*types.Contact {
	res := make([]*types.Contact, len(cs))
	for i := range cs {
		res[i] = cs[i].Response()
	}
	return res
}
