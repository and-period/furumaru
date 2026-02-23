package admin

import (
	"context"

	"golang.org/x/sync/errgroup"
)

func (a *app) getSecret(ctx context.Context, p *params) error {
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		// データベース（TiDB）認証情報の取得
		if a.TiDBSecretName == "" {
			p.tidbHost = a.TiDBHost
			p.tidbPort = a.TiDBPort
			p.tidbUsername = a.TiDBUsername
			p.tidbPassword = a.TiDBPassword
			return nil
		}
		secrets, err := p.secret.Get(ectx, a.TiDBSecretName)
		if err != nil {
			return err
		}
		p.tidbHost = secrets["host"]
		p.tidbPort = secrets["port"]
		p.tidbUsername = secrets["username"]
		p.tidbPassword = secrets["password"]
		return nil
	})
	eg.Go(func() error {
		// Slack認証情報の取得
		if a.SlackSecretName == "" {
			p.slackToken = a.SlackAPIToken
			p.slackChannelID = a.SlackChannelID
			return nil
		}
		secrets, err := p.secret.Get(ectx, a.SlackSecretName)
		if err != nil {
			return err
		}
		p.slackToken = secrets["token"]
		p.slackChannelID = secrets["channelId"]
		return nil
	})
	eg.Go(func() error {
		// New Relic認証情報の取得
		if a.NewRelicSecretName == "" {
			p.newRelicLicense = a.NewRelicLicense
			return nil
		}
		secrets, err := p.secret.Get(ectx, a.NewRelicSecretName)
		if err != nil {
			return err
		}
		p.newRelicLicense = secrets["license"]
		return nil
	})
	eg.Go(func() error {
		// Sentry認証情報の取得
		if a.SentrySecretName == "" {
			p.sentryDsn = a.SentryDsn
			return nil
		}
		secrets, err := p.secret.Get(ectx, a.SentrySecretName)
		if err != nil {
			return err
		}
		p.sentryDsn = secrets["dsn"]
		return nil
	})
	eg.Go(func() error {
		// KOMOJU接続情報の取得
		if a.KomojuSecretName == "" {
			p.komojuClientID = a.KomojuClientID
			p.komojuClientPassword = a.KomojuClientPassword
			return nil
		}
		secrets, err := p.secret.Get(ectx, a.KomojuSecretName)
		if err != nil {
			return err
		}
		p.komojuClientID = secrets["clientId"]
		p.komojuClientPassword = secrets["clientPassword"]
		p.komojuWebhookSecret = secrets["webhookSecret"]
		return nil
	})
	eg.Go(func() error {
		// Google API認証情報の取得
		if a.GoogleSecretName == "" {
			p.googleClientID = a.GoogleClientID
			p.googleClientSecret = a.GoogleClientSecret
			p.googleMapsPlatformAPIKey = a.GoogleMapsPlatformAPIKey
			return nil
		}
		secrets, err := p.secret.Get(ectx, a.GoogleSecretName)
		if err != nil {
			return err
		}
		p.googleClientID = secrets["clientId"]
		p.googleClientSecret = secrets["clientSecret"]
		p.googleMapsPlatformAPIKey = secrets["mapsPlatformAPIKey"]
		return nil
	})
	return eg.Wait()
}
