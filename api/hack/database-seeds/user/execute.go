package user

import (
	"context"
	"fmt"
	"time"

	"github.com/and-period/furumaru/api/hack/database-seeds/common"
	"github.com/and-period/furumaru/api/hack/database-seeds/master"
	"github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/mysql"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const database = "users"

type app struct {
	db     *mysql.Client
	logger *zap.Logger
	now    func() time.Time
	srcDir string
}

func NewClient(params *common.Params) (common.Client, error) {
	db, err := common.NewDBClient(params, database)
	if err != nil {
		return nil, err
	}
	return &app{
		db:     db,
		logger: params.Logger,
		now:    jst.Now,
		srcDir: params.SrcDir,
	}, nil
}

func (a *app) Execute(ctx context.Context) error {
	a.logger.Info("Executing users database seeds...")
	if err := a.executeAdminPolicies(ctx); err != nil {
		return err
	}
	a.logger.Info("Completed admin policies table")
	if err := a.executeAdminRoles(ctx); err != nil {
		return err
	}
	a.logger.Info("Completed admin roles table")
	if err := a.executeAdminGroups(ctx); err != nil {
		return err
	}
	a.logger.Info("Completed admin groups table")
	a.logger.Info("Completed users database seeds")
	return nil
}

func (a *app) executeAdminPolicies(ctx context.Context) error {
	return a.db.Transaction(ctx, func(tx *gorm.DB) error {
		for _, policy := range master.AdminPolicies {
			now := a.now()

			policy.CreatedAt = now
			policy.UpdatedAt = now

			updates := map[string]interface{}{
				"name":        policy.Name,
				"description": policy.Description,
				"priority":    policy.Priority,
				"path":        policy.Path,
				"method":      policy.Method,
				"action":      policy.Action,
				"updated_at":  policy.UpdatedAt,
			}
			stmt := tx.WithContext(ctx).Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "id"}},
				DoUpdates: clause.Assignments(updates),
			})
			if err := stmt.Create(&policy).Error; err != nil {
				return fmt.Errorf("failed to create admin policy. policy=%+v: %w", policy, err)
			}
		}
		return nil
	})
}

func (a *app) executeAdminRoles(ctx context.Context) error {
	return a.db.Transaction(ctx, func(tx *gorm.DB) error {
		for _, role := range master.AdminRoles {
			/**
			 * 管理者ロールの作成
			 */
			now := a.now()

			role.CreatedAt = now
			role.UpdatedAt = now

			updates := map[string]interface{}{
				"name":        role.Name,
				"description": role.Description,
				"updated_at":  role.UpdatedAt,
			}
			stmt := tx.WithContext(ctx).Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "id"}},
				DoUpdates: clause.Assignments(updates),
			})
			if err := stmt.Create(&role.AdminRole).Error; err != nil {
				return fmt.Errorf("failed to create admin role. role=%+v: %w", role, err)
			}

			if len(role.PolicyIDs) == 0 {
				continue
			}

			/**
			 * 管理者ロールと管理者ポリシーの紐付け
			 */
			stmt = tx.WithContext(ctx).
				Where("role_id = ?", role.ID).
				Where("policy_id NOT IN (?)", role.PolicyIDs)

			if err := stmt.Delete(&entity.AdminRolePolicy{}).Error; err != nil {
				return fmt.Errorf("failed to delete admin role policy. roleID=%s: %w", role.ID, err)
			}

			for _, policyID := range role.PolicyIDs {
				now := a.now()

				rolePolicy := &entity.AdminRolePolicy{
					RoleID:    role.ID,
					PolicyID:  policyID,
					CreatedAt: now,
					UpdatedAt: now,
				}
				stmt := tx.WithContext(ctx).Clauses(clause.OnConflict{
					Columns:   []clause.Column{{Name: "role_id"}, {Name: "policy_id"}},
					DoNothing: true,
				})
				if err := stmt.Create(&rolePolicy).Error; err != nil {
					return fmt.Errorf("failed to create admin role policy. rolePolicy=%+v: %w", rolePolicy, err)
				}
			}
		}
		return nil
	})
}

func (a *app) executeAdminGroups(ctx context.Context) error {
	return a.db.Transaction(ctx, func(tx *gorm.DB) error {
		for _, group := range master.AdminGroups {
			now := a.now()

			group.CreatedAt = now
			group.UpdatedAt = now

			updates := map[string]interface{}{
				"type":        group.Type,
				"name":        group.Name,
				"description": group.Description,
				"updated_at":  group.UpdatedAt,
				"deleted_at":  nil,
			}
			stmt := tx.WithContext(ctx).Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "id"}},
				DoUpdates: clause.Assignments(updates),
			})
			if err := stmt.Create(&group.AdminGroup).Error; err != nil {
				return fmt.Errorf("failed to create admin group. group=%+v: %w", group, err)
			}

			if len(group.RoleIDs) == 0 {
				continue
			}

			stmt = tx.WithContext(ctx).
				Where("group_id = ?", group.ID).
				Where("role_id NOT IN (?)", group.RoleIDs)

			if err := stmt.Delete(&entity.AdminGroupRole{}).Error; err != nil {
				return fmt.Errorf("failed to delete admin group role. groupID=%s: %w", group.ID, err)
			}

			for _, roleID := range group.RoleIDs {
				now := a.now()

				groupRole := &entity.AdminGroupRole{
					GroupID:   group.ID,
					RoleID:    roleID,
					CreatedAt: now,
					UpdatedAt: now,
				}
				updates := map[string]interface{}{
					"deleted_at": nil,
				}
				stmt := tx.WithContext(ctx).Clauses(clause.OnConflict{
					Columns:   []clause.Column{{Name: "group_id"}, {Name: "role_id"}},
					DoUpdates: clause.Assignments(updates),
				})
				if err := stmt.Create(&groupRole).Error; err != nil {
					return fmt.Errorf("failed to create admin group role. groupRole=%+v: %w", groupRole, err)
				}
			}
		}
		return nil
	})
}
