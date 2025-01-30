package service

import (
	"bytes"
	"context"
	"encoding/csv"

	"github.com/and-period/furumaru/api/internal/user"
	"github.com/and-period/furumaru/api/internal/user/database"
	"github.com/and-period/furumaru/api/internal/user/entity"
	"golang.org/x/sync/errgroup"
)

const adminRoleModel = `
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act, eft

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow)) && !some(where (p.eft == deny))

[matchers]
m = g(r.sub, p.sub) && keyMatch2(r.obj, p.obj) && regexMatch(r.act, p.act)
`

func (s *service) GenerateAdminRole(ctx context.Context, in *user.GenerateAdminRoleInput) (string, string, error) {
	if err := s.validator.Struct(in); err != nil {
		return "", "", internalError(err)
	}

	var (
		policies     entity.AdminPolicies
		rolePolicies entity.AdminRolePolicies
		groupRoles   entity.AdminGroupRoles
	)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		params := &database.ListAdminPoliciesParams{}
		policies, err = s.db.AdminPolicy.List(ectx, params)
		return
	})
	eg.Go(func() (err error) {
		params := &database.ListAdminRolePoliciesParams{}
		rolePolicies, err = s.db.AdminRolePolicy.List(ectx, params)
		return
	})
	eg.Go(func() (err error) {
		params := &database.ListAdminGroupRolesParams{}
		groupRoles, err = s.db.AdminGroupRole.List(ectx, params)
		return
	})
	if err := eg.Wait(); err != nil {
		return "", "", internalError(err)
	}

	buf := &bytes.Buffer{}
	writer := csv.NewWriter(buf)
	if err := policies.Write(writer); err != nil {
		return "", "", internalError(err)
	}
	if err := rolePolicies.Write(writer); err != nil {
		return "", "", internalError(err)
	}
	if err := groupRoles.Write(writer); err != nil {
		return "", "", internalError(err)
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		return "", "", internalError(err)
	}

	return adminRoleModel, buf.String(), nil
}
