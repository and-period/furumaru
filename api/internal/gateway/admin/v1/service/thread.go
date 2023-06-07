package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/messenger/entity"
)

type Thread struct {
	response.Thread
}

type Threads []*Thread

func NewThread(thread *entity.Thread) *Thread {
	return &Thread{
		Thread: response.Thread{
			ID:        thread.ID,
			ContactID: thread.ContactID,
			UserID:    thread.UserID,
			UserType:  thread.UserType,
			Content:   thread.Content,
			CreatedAt: thread.CreatedAt.Unix(),
			UpdatedAt: thread.UpdatedAt.Unix(),
		},
	}
}

func (t *Thread) Response() *response.Thread {
	return &t.Thread
}

func NewThreads(threads entity.Threads) Threads {
	res := make(Threads, len(threads))
	for i := range threads {
		res[i] = NewThread(threads[i])
	}
	return res
}

func (ts Threads) Response() []*response.Thread {
	res := make([]*response.Thread, len(ts))
	for i := range ts {
		res[i] = ts[i].Response()
	}
	return res
}
