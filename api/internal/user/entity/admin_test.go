package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func AdminRole_Validate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		role   AdminRole
		expect error
	}{
		{
			name:   "administrator",
			role:   AdminRoleAdministrator,
			expect: nil,
		},
		{
			name:   "developer",
			role:   AdminRoleDeveloper,
			expect: nil,
		},
		{
			name:   "operator",
			role:   AdminRoleOperator,
			expect: nil,
		},
		{
			name:   "unknown",
			role:   AdminRoleUnknown,
			expect: errInvalidAdminRole,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.role.Validate()
			assert.ErrorIs(t, err, tt.expect)
		})
	}
}
