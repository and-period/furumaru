package admin

import (
	"fmt"
	"time"

	khandler "github.com/and-period/furumaru/api/internal/gateway/admin/komoju/handler"
	shandler "github.com/and-period/furumaru/api/internal/gateway/admin/stripe/handler"
	v1 "github.com/and-period/furumaru/api/internal/gateway/admin/v1/handler"
	"github.com/and-period/furumaru/api/internal/media"
	mediadb "github.com/and-period/furumaru/api/internal/media/database/tidb"
	mediasrv "github.com/and-period/furumaru/api/internal/media/service"
	"github.com/and-period/furumaru/api/internal/messenger"
	messengerdb "github.com/and-period/furumaru/api/internal/messenger/database/tidb"
	messengersrv "github.com/and-period/furumaru/api/internal/messenger/service"
	"github.com/and-period/furumaru/api/internal/store"
	storedb "github.com/and-period/furumaru/api/internal/store/database/tidb"
	storesrv "github.com/and-period/furumaru/api/internal/store/service"
	"github.com/and-period/furumaru/api/internal/user"
	userdb "github.com/and-period/furumaru/api/internal/user/database/tidb"
	uentity "github.com/and-period/furumaru/api/internal/user/entity"
	usersrv "github.com/and-period/furumaru/api/internal/user/service"
	"github.com/and-period/furumaru/api/pkg/mysql"
	pkgstripe "github.com/and-period/furumaru/api/pkg/stripe"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rafaelhl/gorm-newrelic-telemetry-plugin/telemetry"
)

func (a *app) injectServices(p *params) error {
	// Serviceの設定
	mediaService, err := a.newMediaService(p)
	if err != nil {
		return fmt.Errorf("cmd: failed to create media service: %w", err)
	}
	messengerService, err := a.newMessengerService(p)
	if err != nil {
		return fmt.Errorf("cmd: failed to create messenger service: %w", err)
	}
	userService, err := a.newUserService(p, mediaService, messengerService)
	if err != nil {
		return fmt.Errorf("cmd: failed to create user service: %w", err)
	}
	storeService, err := a.newStoreService(p, userService, mediaService, messengerService)
	if err != nil {
		return fmt.Errorf("cmd: failed to create store service: %w", err)
	}

	// Handlerの設定
	v1Params := &v1.Params{
		WaitGroup: p.waitGroup,
		User:      userService,
		Store:     storeService,
		Messenger: messengerService,
		Media:     mediaService,
	}
	khandlerParams := &khandler.Params{
		WaitGroup:     p.waitGroup,
		Store:         storeService,
		WebhookSecret: p.komojuWebhookSecret,
	}
	a.v1 = v1.NewHandler(v1Params,
		v1.WithEnvironment(a.Environment),
		v1.WithSentry(p.sentry),
	)
	a.komoju = khandler.NewHandler(khandlerParams,
		khandler.WithEnvironment(a.Environment),
		khandler.WithSentry(p.sentry),
	)
	if p.stripeWebhookSecret != "" {
		receiver := pkgstripe.NewReceiver(&pkgstripe.Params{
			WebhookKey: p.stripeWebhookSecret,
		})
		shandlerParams := &shandler.Params{
			WaitGroup: p.waitGroup,
			Store:     storeService,
			Receiver:  receiver,
		}
		a.stripe = shandler.NewHandler(shandlerParams,
			shandler.WithEnvironment(a.Environment),
			shandler.WithSentry(p.sentry),
		)
	}
	return nil
}

func (a *app) newTiDB(dbname string, p *params) (*mysql.Client, error) {
	params := &mysql.Params{
		Host:     p.tidbHost,
		Port:     p.tidbPort,
		Database: dbname,
		Username: p.tidbUsername,
		Password: p.tidbPassword,
	}
	location, err := time.LoadLocation(a.DBTimeZone)
	if err != nil {
		return nil, err
	}
	cli, err := mysql.NewTiDBClient(
		params,
		mysql.WithNow(p.now),
		mysql.WithLocation(location),
	)
	if err != nil {
		return nil, err
	}
	if err := cli.DB.Use(telemetry.NewNrTracer(dbname, p.tidbHost, string(newrelic.DatastoreMySQL))); err != nil {
		return nil, err
	}
	return cli, nil
}

func (a *app) newMediaService(p *params) (media.Service, error) {
	mysql, err := a.newTiDB("media", p)
	if err != nil {
		return nil, err
	}
	user, err := a.newUserService(p, nil, nil)
	if err != nil {
		return nil, err
	}
	store, err := a.newStoreService(p, nil, nil, nil)
	if err != nil {
		return nil, err
	}
	batchUpdateArchiveCommand := func(broadcastID string) []string {
		// see: ./hack/media-update-archive/main.go
		return []string{
			"./app",
			"-environment", a.Environment,
			"-db-secret-name", a.TiDBSecretName,
			"-sentry-secret-name", a.SentrySecretName,
			"-assets-domain", p.assetsURL.Host,
			"-s3-bucket", a.S3Bucket,
			"-broadcast-id", broadcastID,
		}
	}
	params := &mediasrv.Params{
		WaitGroup:                    p.waitGroup,
		Database:                     mediadb.NewDatabase(mysql),
		Cache:                        p.cache,
		MediaLive:                    p.medialive,
		Youtube:                      p.youtube,
		Storage:                      p.storage,
		Tmp:                          p.tmpStorage,
		Producer:                     p.mediaQueue,
		Batch:                        p.batch,
		BatchUpdateArchiveDefinition: a.BatchMediaUpdateArchiveDefinition,
		BatchUpdateArchiveQueue:      a.BatchMediaUpdateArchiveQueue,
		BatchUpdateArchiveCommand:    batchUpdateArchiveCommand,
		User:                         user,
		Store:                        store,
	}
	return mediasrv.NewService(params)
}

func (a *app) newMessengerService(p *params) (messenger.Service, error) {
	mysql, err := a.newTiDB("messengers", p)
	if err != nil {
		return nil, err
	}
	user, err := a.newUserService(p, nil, nil)
	if err != nil {
		return nil, err
	}
	store, err := a.newStoreService(p, nil, nil, nil)
	if err != nil {
		return nil, err
	}
	params := &messengersrv.Params{
		WaitGroup:   p.waitGroup,
		Producer:    p.messengerQueue,
		AdminWebURL: p.adminWebURL,
		UserWebURL:  p.userWebURL,
		Database:    messengerdb.NewDatabase(mysql),
		User:        user,
		Store:       store,
	}
	return messengersrv.NewService(params), nil
}

func (a *app) newUserService(p *params, media media.Service, messenger messenger.Service) (user.Service, error) {
	mysql, err := a.newTiDB("users", p)
	if err != nil {
		return nil, err
	}
	store, err := a.newStoreService(p, nil, nil, nil)
	if err != nil {
		return nil, err
	}
	groups := map[uentity.AdminType][]string{
		uentity.AdminTypeAdministrator: a.DefaultAdministratorGroupIDs,
		uentity.AdminTypeCoordinator:   a.DefaultCoordinatorGroupIDs,
		uentity.AdminTypeProducer:      a.DefaultProducerGroupIDs,
	}
	params := &usersrv.Params{
		WaitGroup:                  p.waitGroup,
		Database:                   userdb.NewDatabase(mysql),
		Cache:                      p.cache,
		AdminAuth:                  p.adminAuth,
		UserAuth:                   p.userAuth,
		Store:                      store,
		Messenger:                  messenger,
		Media:                      media,
		DefaultAdminGroups:         groups,
		AdminAuthGoogleRedirectURL: a.CognitoAdminGoogleRedirectURL,
		AdminAuthLINERedirectURL:   a.CognitoAdminLINERedirectURL,
	}
	return usersrv.NewService(params), nil
}

func (a *app) newStoreService(
	p *params, user user.Service, media media.Service, messenger messenger.Service,
) (store.Service, error) {
	mysql, err := a.newTiDB("stores", p)
	if err != nil {
		return nil, err
	}
	params := &storesrv.Params{
		WaitGroup:   p.waitGroup,
		Database:    storedb.NewDatabase(mysql),
		Cache:       p.cache,
		User:        user,
		Messenger:   messenger,
		Media:       media,
		PostalCode:  p.postalCode,
		Geolocation: p.geolocation,
		Providers:   p.providers,
	}
	return storesrv.NewService(params), nil
}
