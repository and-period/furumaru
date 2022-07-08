package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/messenger/entity"
)

// ContactStatus - お問い合わせ対応状況
type ContactStatus int32

const (
	ContactStatusUnknown    ContactStatus = 0
	ContactStatusToDo       ContactStatus = 1 // ToDo
	ContactStatusInprogress ContactStatus = 2 // 対応中
	ContactStatusDone       ContactStatus = 3 // 完了
	ContactStatusDiscard    ContactStatus = 4 // 対応不要
)

// ContactPriority - お問い合わせ対応優先度
type ContactPriority int32

const (
	ContactPriorityUnknown ContactPriority = 0
	ContactPriorityLow     ContactPriority = 1 // 低
	ContactPriorityMiddle  ContactPriority = 2 // 中
	ContactPriorityHigh    ContactPriority = 3 // 高
)

type Contact struct {
	response.Contact
}

type Contacts []*Contact

func NewContactStatus(status entity.ContactStatus) ContactStatus {
	switch status {
	case entity.ContactStatusToDo:
		return ContactStatusToDo
	case entity.ContactStatusInprogress:
		return ContactStatusInprogress
	case entity.ContactStatusDone:
		return ContactStatusDone
	case entity.ContactStatusDiscard:
		return ContactStatusDiscard
	default:
		return ContactStatusUnknown
	}
}

func (s ContactStatus) MessengerEntity() entity.ContactStatus {
	switch s {
	case ContactStatusToDo:
		return entity.ContactStatusToDo
	case ContactStatusInprogress:
		return entity.ContactStatusInprogress
	case ContactStatusDone:
		return entity.ContactStatusDone
	case ContactStatusDiscard:
		return entity.ContactStatusDiscard
	default:
		return entity.ContactStatusUnknown
	}
}

func (s ContactStatus) Response() int32 {
	return int32(s)
}

func NewContactPriority(priority entity.ContactPriority) ContactPriority {
	switch priority {
	case entity.ContactPriorityLow:
		return ContactPriorityLow
	case entity.ContactPriorityMiddle:
		return ContactPriorityMiddle
	case entity.ContactPriorityHigh:
		return ContactPriorityHigh
	default:
		return ContactPriorityUnknown
	}
}

func (p ContactPriority) MessengerEntity() entity.ContactPriority {
	switch p {
	case ContactPriorityLow:
		return entity.ContactPriorityLow
	case ContactPriorityMiddle:
		return entity.ContactPriorityMiddle
	case ContactPriorityHigh:
		return entity.ContactPriorityHigh
	default:
		return entity.ContactPriorityUnknown
	}
}

func (p ContactPriority) Response() int32 {
	return int32(p)
}

func NewContact(contact *entity.Contact) *Contact {
	return &Contact{
		Contact: response.Contact{
			ID:          contact.ID,
			Title:       contact.Title,
			Content:     contact.Content,
			Username:    contact.Username,
			Email:       contact.Email,
			PhoneNumber: contact.PhoneNumber,
			Status:      NewContactStatus(contact.Status).Response(),
			Priority:    NewContactPriority(contact.Priority).Response(),
			Note:        contact.Note,
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
