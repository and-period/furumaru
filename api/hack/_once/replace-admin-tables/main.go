// admin_authsテーブルからadminsテーブルへデータ移行するためのコマンド
//
// usage: go run ./main.go -db-host=127.0.0.1 -db-port=3316
package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	"github.com/and-period/furumaru/api/internal/exception"
	userdb "github.com/and-period/furumaru/api/internal/user/database"
	"github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/database"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/log"
	"github.com/samber/lo"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const (
	adminTable         = "admins"
	adminAuthTable     = "admin_auths"
	administratorTable = "administrators"
	coordinatorTable   = "coordinators"
	producerTable      = "producers"
)

func main() {
	start := time.Now()
	fmt.Println("Start..")
	if err := run(); err != nil {
		panic(err)
	}
	fmt.Printf("Done: %s\n", time.Since(start))
}

type app struct {
	now    func() time.Time
	logger *zap.Logger
	db     *userdb.Database
}

func run() error {
	app := app{now: jst.Now}
	host := flag.String("db-host", "mysql", "target mysql host")
	port := flag.String("db-port", "3306", "target mysql port")
	username := flag.String("db-username", "root", "target mysql username")
	password := flag.String("db-password", "12345678", "target mysql password")
	flag.Parse()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	logger, err := log.NewLogger()
	if err != nil {
		return err
	}
	defer logger.Sync()
	app.logger = logger

	db, err := app.newDatabase(*host, *port, *username, *password)
	if err != nil {
		return err
	}
	params := &userdb.Params{
		Database: db,
	}
	app.db = userdb.NewDatabase(params)

	_, err = db.Transaction(ctx, func(tx *gorm.DB) (interface{}, error) {
		logger.Info("start fetch admin auths...")
		auths, err := app.fetchAdminAuths(ctx, tx)
		if err != nil {
			return nil, err
		}
		logger.Info("finish fetch admin auths...", zap.Int("total", len(auths)))

		authsMap := auths.GroupByRole()

		admins := make(entity.Admins, 0, len(auths))
		for role, auths := range authsMap {
			var (
				as  entity.Admins
				err error
			)
			switch role {
			case entity.AdminRoleAdministrator:
				as, err = app.newAdminsFromAdministrator(ctx, tx, auths)
			case entity.AdminRoleCoordinator:
				as, err = app.newAdminsFromCoordinator(ctx, tx, auths)
			case entity.AdminRoleProducer:
				as, err = app.newAdminsFromProducer(ctx, tx, auths)
			}
			if err != nil {
				return nil, err
			}
			admins = append(admins, as...)
		}

		return nil, app.upsertAdmins(ctx, tx, admins)
	})
	return err
}

func (a *app) newDatabase(host, port, username, password string) (*database.Client, error) {
	params := &database.Params{
		Socket:   "tcp",
		Host:     host,
		Port:     port,
		Database: "users",
		Username: username,
		Password: password,
	}
	return database.NewClient(params, database.WithLogger(a.logger))
}

func (a *app) fetchAdminAuths(ctx context.Context, tx *gorm.DB) (entity.AdminAuths, error) {
	var auths entity.AdminAuths

	err := tx.WithContext(ctx).
		Table(adminAuthTable).
		Find(&auths).Error
	return auths, err
}

func (a *app) newAdminsFromAdministrator(ctx context.Context, tx *gorm.DB, auths entity.AdminAuths) (entity.Admins, error) {
	administrators, err := a.MultiGetAdministrators(ctx, tx, auths.AdminIDs())
	if err != nil {
		return nil, err
	}
	if err := a.upsertAdministrators(ctx, tx, administrators); err != nil {
		return nil, err
	}
	amap := lo.KeyBy[string, *entity.Administrator](administrators, func(a *entity.Administrator) string {
		return a.ID
	})
	admins := make(entity.Admins, len(auths))
	for i := range auths {
		administrator, ok := amap[auths[i].AdminID]
		if !ok {
			continue
		}
		admins[i] = &entity.Admin{
			ID:            auths[i].AdminID,
			CognitoID:     auths[i].CognitoID,
			Role:          auths[i].Role,
			Lastname:      administrator.Lastname,
			Firstname:     administrator.Firstname,
			LastnameKana:  administrator.LastnameKana,
			FirstnameKana: administrator.FirstnameKana,
			Email:         administrator.Email,
			Device:        auths[i].Device,
			CreatedAt:     auths[i].CreatedAt,
			UpdatedAt:     a.now(),
		}
	}
	return admins, nil
}

func (a *app) newAdminsFromCoordinator(ctx context.Context, tx *gorm.DB, auths entity.AdminAuths) (entity.Admins, error) {
	coordinators, err := a.MultiGetCoordinators(ctx, tx, auths.AdminIDs())
	if err != nil {
		return nil, err
	}
	if err := a.upsertCoordinators(ctx, tx, coordinators); err != nil {
		return nil, err
	}
	cmap := lo.KeyBy[string, *entity.Coordinator](coordinators, func(c *entity.Coordinator) string {
		return c.ID
	})
	admins := make(entity.Admins, len(auths))
	for i := range auths {
		coordinator, ok := cmap[auths[i].AdminID]
		if !ok {
			continue
		}
		admins[i] = &entity.Admin{
			ID:            auths[i].AdminID,
			CognitoID:     auths[i].CognitoID,
			Role:          auths[i].Role,
			Lastname:      coordinator.Lastname,
			Firstname:     coordinator.Firstname,
			LastnameKana:  coordinator.LastnameKana,
			FirstnameKana: coordinator.FirstnameKana,
			Email:         coordinator.Email,
			Device:        auths[i].Device,
			CreatedAt:     auths[i].CreatedAt,
			UpdatedAt:     a.now(),
		}
	}
	return admins, nil
}

func (a *app) newAdminsFromProducer(ctx context.Context, tx *gorm.DB, auths entity.AdminAuths) (entity.Admins, error) {
	producers, err := a.MultiGetProducers(ctx, tx, auths.AdminIDs())
	if err != nil {
		return nil, err
	}
	if err := a.upsertProducers(ctx, tx, producers); err != nil {
		return nil, err
	}
	cmap := lo.KeyBy[string, *entity.Producer](producers, func(c *entity.Producer) string {
		return c.ID
	})
	admins := make(entity.Admins, len(auths))
	for i := range auths {
		producer, ok := cmap[auths[i].AdminID]
		if !ok {
			continue
		}
		admins[i] = &entity.Admin{
			ID:            auths[i].AdminID,
			CognitoID:     auths[i].CognitoID,
			Role:          auths[i].Role,
			Lastname:      producer.Lastname,
			Firstname:     producer.Firstname,
			LastnameKana:  producer.LastnameKana,
			FirstnameKana: producer.FirstnameKana,
			Email:         producer.Email,
			Device:        auths[i].Device,
			CreatedAt:     auths[i].CreatedAt,
			UpdatedAt:     a.now(),
		}
	}
	return admins, nil
}

func (a *app) MultiGetAdministrators(ctx context.Context, tx *gorm.DB, ids []string) (entity.Administrators, error) {
	var admins entity.Administrators

	err := tx.WithContext(ctx).
		Table(administratorTable).
		Where("id IN (?)", ids).
		Find(&admins).Error
	return admins, exception.InternalError(err)
}

func (a *app) MultiGetCoordinators(ctx context.Context, tx *gorm.DB, ids []string) (entity.Coordinators, error) {
	var admins entity.Coordinators

	err := tx.WithContext(ctx).
		Table(coordinatorTable).
		Where("id IN (?)", ids).
		Find(&admins).Error
	return admins, exception.InternalError(err)
}

func (a *app) MultiGetProducers(ctx context.Context, tx *gorm.DB, ids []string) (entity.Producers, error) {
	var admins entity.Producers

	err := tx.WithContext(ctx).
		Table(producerTable).
		Where("id IN (?)", ids).
		Find(&admins).Error
	return admins, exception.InternalError(err)
}

func (a *app) upsertAdmins(ctx context.Context, tx *gorm.DB, as entity.Admins) error {
	return tx.WithContext(ctx).Table(adminTable).Save(&as).Error
}

func (a *app) upsertAdministrators(ctx context.Context, tx *gorm.DB, as entity.Administrators) error {
	for i := range as {
		as[i].AdminID = as[i].ID
	}
	return tx.WithContext(ctx).Table(administratorTable).Save(&as).Error
}

func (a *app) upsertCoordinators(ctx context.Context, tx *gorm.DB, cs entity.Coordinators) error {
	for i := range cs {
		cs[i].AdminID = cs[i].ID
	}
	return tx.WithContext(ctx).Table(coordinatorTable).Save(&cs).Error
}

func (a *app) upsertProducers(ctx context.Context, tx *gorm.DB, ps entity.Producers) error {
	for i := range ps {
		ps[i].AdminID = ps[i].ID
	}
	return tx.WithContext(ctx).Table(producerTable).Save(&ps).Error
}
