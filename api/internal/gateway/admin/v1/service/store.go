package service

import (
	"github.com/and-period/marche/api/internal/gateway/admin/v1/response"
	"github.com/and-period/marche/api/internal/store/entity"
)

type Store struct {
	*response.Store
}

type Stores []*Store

func NewStore(store *entity.Store, staffs Staffs) *Store {
	return &Store{
		Store: &response.Store{
			ID:           store.ID,
			Name:         store.Name,
			ThumbnailURL: store.ThumbnailURL,
			Staffs:       staffs.Response(),
			CreatedAt:    store.CreatedAt.Unix(),
			UpdatedAt:    store.UpdatedAt.Unix(),
		},
	}
}

func (s *Store) Response() *response.Store {
	return s.Store
}

func NewStores(stores entity.Stores) Stores {
	res := make(Stores, len(stores))
	for i := range stores {
		res[i] = NewStore(stores[i], nil)
	}
	return res
}

func (ss Stores) Response() []*response.Store {
	res := make([]*response.Store, len(ss))
	for i := range ss {
		res[i] = ss[i].Response()
	}
	return res
}
