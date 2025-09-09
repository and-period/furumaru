package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/types"
	"github.com/and-period/furumaru/api/internal/messenger/entity"
)

type ThreadUserType int32

const (
	ThreadUserTypeUnknown ThreadUserType = iota // 不明
	ThreadUserTypeAdmin                         // 管理者
	ThreadUserTypeUser                          // ユーザー
	ThreadUserTypeGuest                         // ゲスト(ユーザIDなし)
)

func NewThreadUserType(typ entity.ThreadUserType) ThreadUserType {
	switch typ {
	case entity.ThreadUserTypeAdmin:
		return ThreadUserTypeAdmin
	case entity.ThreadUserTypeUser:
		return ThreadUserTypeUser
	case entity.ThreadUserTypeGuest:
		return ThreadUserTypeGuest
	default:
		return ThreadUserTypeUnknown
	}
}

func (t ThreadUserType) StoreEntity() entity.ThreadUserType {
	switch t {
	case ThreadUserTypeAdmin:
		return entity.ThreadUserTypeAdmin
	case ThreadUserTypeUser:
		return entity.ThreadUserTypeUser
	case ThreadUserTypeGuest:
		return entity.ThreadUserTypeGuest
	default:
		return entity.ThreadUserTypeUnknown
	}
}

type Thread struct {
	types.Thread
}

type Threads []*Thread

func (t ThreadUserType) Response() int32 {
	return int32(t)
}

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
