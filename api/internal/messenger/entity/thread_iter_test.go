package entity

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestThreads_All(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		threads Threads
	}{
		{
			name: "success",
			threads: Threads{
				{
					ID:        "thread-id01",
					ContactID: "contact-id",
					UserType:  ThreadUserTypeAdmin,
					UserID:    "admin-id",
					Content:   "content1",
				},
				{
					ID:        "thread-id02",
					ContactID: "contact-id",
					UserType:  ThreadUserTypeUser,
					UserID:    "user-id",
					Content:   "content2",
				},
			},
		},
		{
			name:    "empty",
			threads: Threads{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var indices []int
			var ids []string
			for i, th := range tt.threads.All() {
				indices = append(indices, i)
				ids = append(ids, th.ID)
			}
			for i, th := range tt.threads {
				if i < len(indices) {
					assert.Equal(t, i, indices[i])
					assert.Equal(t, th.ID, ids[i])
				}
			}
			assert.Len(t, indices, len(tt.threads))
		})
	}
}

func TestThreads_All_EarlyBreak(t *testing.T) {
	t.Parallel()
	threads := Threads{
		{ID: "thread-id01"},
		{ID: "thread-id02"},
		{ID: "thread-id03"},
	}
	var count int
	for range threads.All() {
		count++
		if count == 2 {
			break
		}
	}
	assert.Equal(t, 2, count)
}

func TestThreads_IterIDs(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		threads Threads
		expect  []string
	}{
		{
			name: "success",
			threads: Threads{
				{ID: "thread-id01", ContactID: "contact-id", UserType: ThreadUserTypeAdmin, UserID: "admin-id"},
				{ID: "thread-id02", ContactID: "contact-id", UserType: ThreadUserTypeUser, UserID: "user-id"},
			},
			expect: []string{"thread-id01", "thread-id02"},
		},
		{
			name:    "empty",
			threads: Threads{},
			expect:  nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := slices.Collect(tt.threads.IterIDs())
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestThreads_IterUserIDs(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		threads Threads
		expect  []string
	}{
		{
			name: "success with mixed types",
			threads: Threads{
				{ID: "thread-id01", ContactID: "contact-id", UserType: ThreadUserTypeAdmin, UserID: "admin-id"},
				{ID: "thread-id02", ContactID: "contact-id", UserType: ThreadUserTypeUser, UserID: "user-id01"},
				{ID: "thread-id03", ContactID: "contact-id", UserType: ThreadUserTypeGuest, UserID: ""},
				{ID: "thread-id04", ContactID: "contact-id", UserType: ThreadUserTypeUser, UserID: "user-id02"},
			},
			expect: []string{"user-id01", "user-id02"},
		},
		{
			name: "no user type threads",
			threads: Threads{
				{ID: "thread-id01", ContactID: "contact-id", UserType: ThreadUserTypeAdmin, UserID: "admin-id"},
			},
			expect: nil,
		},
		{
			name:    "empty",
			threads: Threads{},
			expect:  nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := slices.Collect(tt.threads.IterUserIDs())
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestThreads_IterAdminIDs(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		threads Threads
		expect  []string
	}{
		{
			name: "success with mixed types",
			threads: Threads{
				{ID: "thread-id01", ContactID: "contact-id", UserType: ThreadUserTypeAdmin, UserID: "admin-id01"},
				{ID: "thread-id02", ContactID: "contact-id", UserType: ThreadUserTypeUser, UserID: "user-id"},
				{ID: "thread-id03", ContactID: "contact-id", UserType: ThreadUserTypeAdmin, UserID: "admin-id02"},
			},
			expect: []string{"admin-id01", "admin-id02"},
		},
		{
			name: "no admin type threads",
			threads: Threads{
				{ID: "thread-id01", ContactID: "contact-id", UserType: ThreadUserTypeUser, UserID: "user-id"},
			},
			expect: nil,
		},
		{
			name:    "empty",
			threads: Threads{},
			expect:  nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := slices.Collect(tt.threads.IterAdminIDs())
			assert.Equal(t, tt.expect, actual)
		})
	}
}
