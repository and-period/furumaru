package admin

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/and-period/furumaru/api/internal/store/komoju"
	kpayment "github.com/and-period/furumaru/api/internal/store/komoju/payment"
	ksession "github.com/and-period/furumaru/api/internal/store/komoju/session"
	"github.com/and-period/furumaru/api/pkg/geolocation"
	"github.com/and-period/furumaru/api/pkg/postalcode"
	"github.com/and-period/furumaru/api/pkg/sentry"
	"github.com/and-period/furumaru/api/pkg/slack"
	"github.com/and-period/furumaru/api/pkg/youtube"
	"github.com/newrelic/go-agent/v3/newrelic"
)

func (a *app) injectExternal(p *params) error {
	// New Relicの設定
	if p.newRelicLicense != "" {
		appName := fmt.Sprintf("%s-%s", a.AppName, a.Environment)
		labels := map[string]string{
			"app":     "furumaru",
			"env":     a.Environment,
			"service": a.AppName,
			"type":    "backend",
		}
		newrelicApp, err := newrelic.NewApplication(
			newrelic.ConfigAppName(appName),
			newrelic.ConfigLicense(p.newRelicLicense),
			newrelic.ConfigAppLogMetricsEnabled(true),
			newrelic.ConfigAppLogForwardingEnabled(true),
			newrelic.ConfigCustomInsightsEventsEnabled(true),
			newrelic.ConfigAppLogEnabled(true),
			newrelic.ConfigAppLogForwardingEnabled(true),
			func(cfg *newrelic.Config) {
				cfg.HostDisplayName = appName
				cfg.Labels = labels
			},
		)
		if err != nil {
			return fmt.Errorf("cmd: failed to create newrelic client: %w", err)
		}
		p.newRelic = newrelicApp
	}

	// Sentryの設定
	if p.sentryDsn != "" {
		sentryApp, err := sentry.NewClient(
			sentry.WithServerName(a.AppName),
			sentry.WithEnvironment(a.Environment),
			sentry.WithDSN(p.sentryDsn),
			sentry.WithBind(true),
			sentry.WithTracesSampleRate(a.TraceSampleRate),
		)
		if err != nil {
			return fmt.Errorf("cmd: failed to create sentry client: %w", err)
		}
		p.sentry = sentryApp
	}

	// Slackの設定
	if p.slackToken != "" {
		slackParams := &slack.Params{
			Token:     p.slackToken,
			ChannelID: p.slackChannelID,
		}
		p.slack = slack.NewClient(slackParams)
	}

	// KOMOJUの設定
	kpaymentParams := &kpayment.Params{
		Host:         a.KomojuHost,
		ClientID:     p.komojuClientID,
		ClientSecret: p.komojuClientPassword,
	}
	ksessionParams := &ksession.Params{
		Host:         a.KomojuHost,
		ClientID:     p.komojuClientID,
		ClientSecret: p.komojuClientPassword,
	}
	komojuOpts := []komoju.Option{
		komoju.WithDebugMode(p.debugMode),
	}
	komojuParams := &komoju.Params{
		Payment: kpayment.NewClient(&http.Client{}, kpaymentParams, komojuOpts...),
		Session: ksession.NewClient(&http.Client{}, ksessionParams, komojuOpts...),
	}
	p.komoju = komoju.NewKomoju(komojuParams)

	// PostalCodeの設定
	p.postalCode = postalcode.NewClient(&http.Client{})

	// Geolocationの設定
	geolocationParams := &geolocation.Params{
		APIKey: p.googleMapsPlatformAPIKey,
	}
	geolocation, err := geolocation.NewClient(geolocationParams)
	if err != nil {
		return fmt.Errorf("cmd: failed to create geolocation client: %w", err)
	}
	p.geolocation = geolocation

	// WebURLの設定
	adminWebURL, err := url.Parse(a.AdminWebURL)
	if err != nil {
		return fmt.Errorf("cmd: failed to parse admin web url: %w", err)
	}
	p.adminWebURL = adminWebURL
	userWebURL, err := url.Parse(a.UserWebURL)
	if err != nil {
		return fmt.Errorf("cmd: failed to parse user web url: %w", err)
	}
	p.userWebURL = userWebURL
	assetsURL, err := url.Parse(a.AssetsURL)
	if err != nil {
		return fmt.Errorf("cmd: failed to parse assets url: %w", err)
	}
	p.assetsURL = assetsURL

	// Youtubeの設定
	youtubeParams := &youtube.Params{
		ClientID:        p.googleClientID,
		ClientSecret:    p.googleClientSecret,
		AuthCallbackURL: a.YoutubeAuthCallbackURL,
	}
	p.youtube = youtube.NewClient(youtubeParams)

	return nil
}
