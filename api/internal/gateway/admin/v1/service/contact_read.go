package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/messenger/entity"
)

type ContactRead struct {
	response.ContactRead
}

type ContactReads []*ContactRead

func NewContactRead(contactRead *entity.ContactRead) *ContactRead {
	return &ContactRead{
		ContactRead: response.ContactRead{
			ID:        contactRead.ID,
			ContactID: contactRead.ContactID,
			UserID:    contactRead.UserID,
			UserType:  int32(contactRead.UserType),
			Read:      contactRead.Read,
			CreatedAt: contactRead.CreatedAt.Unix(),
			UpdatedAt: contactRead.UpdatedAt.Unix(),
		},
	}
}

func (c *ContactRead) Response() *response.ContactRead {
	return &c.ContactRead
}

func NewContactReads(contactReads entity.ContactReads) ContactReads {
	res := make(ContactReads, len(contactReads))
	for i := range contactReads {
		res[i] = NewContactRead(contactReads[i])
	}
	return res
}

func (cs ContactReads) Response() []*response.ContactRead {
	res := make([]*response.ContactRead, len(cs))
	for i := range cs {
		res[i] = cs[i].Response()
	}
	return res
}
