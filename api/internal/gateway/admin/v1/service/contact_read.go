package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/types"
	"github.com/and-period/furumaru/api/internal/messenger/entity"
)

type ContactRead struct {
	types.ContactRead
}

type ContactReads []*ContactRead

func NewContactRead(contactRead *entity.ContactRead) *ContactRead {
	return &ContactRead{
		ContactRead: types.ContactRead{
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

func (c *ContactRead) Response() *types.ContactRead {
	return &c.ContactRead
}

func NewContactReads(contactReads entity.ContactReads) ContactReads {
	res := make(ContactReads, len(contactReads))
	for i := range contactReads {
		res[i] = NewContactRead(contactReads[i])
	}
	return res
}

func (cs ContactReads) Response() []*types.ContactRead {
	res := make([]*types.ContactRead, len(cs))
	for i := range cs {
		res[i] = cs[i].Response()
	}
	return res
}
