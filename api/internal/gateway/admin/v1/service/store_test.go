package service

import (
	"testing"

	"github.com/and-period/marche/api/internal/gateway/admin/v1/response"
	"github.com/and-period/marche/api/internal/store/entity"
	"github.com/and-period/marche/api/pkg/jst"
	"github.com/stretchr/testify/assert"
)

func TestStore(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		store  *entity.Store
		staffs Staffs
		expect *Store
	}{
		{
			name: "success",
			store: &entity.Store{
				ID:           1,
				Name:         "&.農園",
				ThumbnailURL: "https://and-period.jp",
				CreatedAt:    jst.Date(2022, 1, 1, 0, 0, 0, 0),
				UpdatedAt:    jst.Date(2022, 1, 1, 0, 0, 0, 0),
			},
			staffs: Staffs{
				{
					Staff: &response.Staff{
						ID:    "kSByoE6FetnPs5Byk3a9Zx",
						Name:  "&. スタッフ1",
						Email: "test-user01@and-period.jp",
						Role:  1,
					},
				},
				{
					Staff: &response.Staff{
						ID:    "kSByoE6FetnPs5Byk3a9Za",
						Name:  "&. スタッフ2",
						Email: "test-user02@and-period.jp",
						Role:  2,
					},
				},
			},
			expect: &Store{
				Store: &response.Store{
					ID:           1,
					Name:         "&.農園",
					ThumbnailURL: "https://and-period.jp",
					Staffs: []*response.Staff{
						{
							ID:    "kSByoE6FetnPs5Byk3a9Zx",
							Name:  "&. スタッフ1",
							Email: "test-user01@and-period.jp",
							Role:  1,
						},
						{
							ID:    "kSByoE6FetnPs5Byk3a9Za",
							Name:  "&. スタッフ2",
							Email: "test-user02@and-period.jp",
							Role:  2,
						},
					},
					CreatedAt: 1640962800,
					UpdatedAt: 1640962800,
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, NewStore(tt.store, tt.staffs))
		})
	}
}

func TestStore_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		store  *Store
		expect *response.Store
	}{
		{
			name: "success",
			store: &Store{
				Store: &response.Store{
					ID:           1,
					Name:         "&.農園",
					ThumbnailURL: "https://and-period.jp",
					Staffs: []*response.Staff{
						{
							ID:    "kSByoE6FetnPs5Byk3a9Zx",
							Name:  "&. スタッフ1",
							Email: "test-user01@and-period.jp",
							Role:  1,
						},
						{
							ID:    "kSByoE6FetnPs5Byk3a9Za",
							Name:  "&. スタッフ2",
							Email: "test-user02@and-period.jp",
							Role:  2,
						},
					},
					CreatedAt: 1640962800,
					UpdatedAt: 1640962800,
				},
			},
			expect: &response.Store{
				ID:           1,
				Name:         "&.農園",
				ThumbnailURL: "https://and-period.jp",
				Staffs: []*response.Staff{
					{
						ID:    "kSByoE6FetnPs5Byk3a9Zx",
						Name:  "&. スタッフ1",
						Email: "test-user01@and-period.jp",
						Role:  1,
					},
					{
						ID:    "kSByoE6FetnPs5Byk3a9Za",
						Name:  "&. スタッフ2",
						Email: "test-user02@and-period.jp",
						Role:  2,
					},
				},
				CreatedAt: 1640962800,
				UpdatedAt: 1640962800,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.store.Response())
		})
	}
}

func TestStores(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		stores entity.Stores
		expect Stores
	}{
		{
			name: "success",
			stores: entity.Stores{
				{
					ID:           1,
					Name:         "&.農園",
					ThumbnailURL: "https://and-period.jp",
					CreatedAt:    jst.Date(2022, 1, 1, 0, 0, 0, 0),
					UpdatedAt:    jst.Date(2022, 1, 1, 0, 0, 0, 0),
				},
				{
					ID:           2,
					Name:         "&.水産",
					ThumbnailURL: "https://and-period.jp",
					CreatedAt:    jst.Date(2022, 1, 1, 0, 0, 0, 0),
					UpdatedAt:    jst.Date(2022, 1, 1, 0, 0, 0, 0),
				},
			},
			expect: Stores{
				{
					Store: &response.Store{
						ID:           1,
						Name:         "&.農園",
						ThumbnailURL: "https://and-period.jp",
						Staffs:       []*response.Staff{},
						CreatedAt:    1640962800,
						UpdatedAt:    1640962800,
					},
				},
				{
					Store: &response.Store{
						ID:           2,
						Name:         "&.水産",
						ThumbnailURL: "https://and-period.jp",
						Staffs:       []*response.Staff{},
						CreatedAt:    1640962800,
						UpdatedAt:    1640962800,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, NewStores(tt.stores))
		})
	}
}

func TestStores_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		stores Stores
		expect []*response.Store
	}{
		{
			name: "success",
			stores: Stores{
				{
					Store: &response.Store{
						ID:           1,
						Name:         "&.農園",
						ThumbnailURL: "https://and-period.jp",
						Staffs:       []*response.Staff{},
						CreatedAt:    1640962800,
						UpdatedAt:    1640962800,
					},
				},
				{
					Store: &response.Store{
						ID:           2,
						Name:         "&.水産",
						ThumbnailURL: "https://and-period.jp",
						Staffs:       []*response.Staff{},
						CreatedAt:    1640962800,
						UpdatedAt:    1640962800,
					},
				},
			},
			expect: []*response.Store{
				{
					ID:           1,
					Name:         "&.農園",
					ThumbnailURL: "https://and-period.jp",
					Staffs:       []*response.Staff{},
					CreatedAt:    1640962800,
					UpdatedAt:    1640962800,
				},
				{
					ID:           2,
					Name:         "&.水産",
					ThumbnailURL: "https://and-period.jp",
					Staffs:       []*response.Staff{},
					CreatedAt:    1640962800,
					UpdatedAt:    1640962800,
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.stores.Response())
		})
	}
}
