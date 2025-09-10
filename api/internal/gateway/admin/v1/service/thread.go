package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/types"
	"github.com/and-period/furumaru/api/internal/messenger/entity"
)

type ThreadUserType types.ThreadUserType

func NewThreadUserType(typ entity.ThreadUserType) ThreadUserType {
	switch typ {
	case entity.ThreadUserTypeAdmin:
		return ThreadUserType(types.ThreadUserTypeAdmin)
	case entity.ThreadUserTypeUser:
		return ThreadUserType(types.ThreadUserTypeUser)
	case entity.ThreadUserTypeGuest:
		return ThreadUserType(types.ThreadUserTypeGuest)
	default:
		return ThreadUserType(types.ThreadUserTypeUnknown)
	}
}

func (t ThreadUserType) StoreEntity() entity.ThreadUserType {
	switch types.ThreadUserType(t) {
	case types.ThreadUserTypeAdmin:
		return entity.ThreadUserTypeAdmin
	case types.ThreadUserTypeUser:
		return entity.ThreadUserTypeUser
	case types.ThreadUserTypeGuest:
		return entity.ThreadUserTypeGuest
	default:
		return entity.ThreadUserTypeUnknown
	}
}

func (t ThreadUserType) Response() types.ThreadUserType {
	return types.ThreadUserType(t)
}

type Thread struct {
	types.Thread
}

type Threads []*Thread

func NewThread(thread *entity.Thread) *Thread {
	return &Thread{
		Thread: types.Thread{
			ID:        thread.ID,
			ContactID: thread.ContactID,
			UserID:    thread.UserID,
			UserType:  NewThreadUserType(thread.UserType).Response(),
			Content:   thread.Content,
			CreatedAt: thread.CreatedAt.Unix(),
			UpdatedAt: thread.UpdatedAt.Unix(),
		},
	}
}

func (t *Thread) Response() *types.Thread {
	return &t.Thread
}

func NewThreads(threads entity.Threads) Threads {
	res := make(Threads, len(threads))
	for i := range threads {
		res[i] = NewThread(threads[i])
	}
	return res
}

func (ts Threads) Response() []*types.Thread {
	res := make([]*types.Thread, len(ts))
	for i := range ts {
		res[i] = ts[i].Response()
	}
	return res
}
