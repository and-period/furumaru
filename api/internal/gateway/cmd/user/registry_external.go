package user

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/and-period/furumaru/api/internal/gateway/user/facility/auth"
	"github.com/and-period/furumaru/api/internal/store/komoju"
	kpayment "github.com/and-period/furumaru/api/internal/store/komoju/payment"
	ksession "github.com/and-period/furumaru/api/internal/store/komoju/session"
	"github.com/and-period/furumaru/api/pkg/geolocation"
	"github.com/and-period/furumaru/api/pkg/postalcode"
	"github.com/and-period/furumaru/api/pkg/sentry"
	"github.com/and-period/furumaru/api/pkg/slack"
	"github.com/newrelic/go-agent/v3/newrelic"
)

func (a *app) injectExternal(ctx context.Context, p *params) error {
	// New Relicの設定
	if p.newRelicLicense != "" {
		labels := map[string]string{
			"app":     "furumaru",
			"env":     a.Environment,
			"service": p.serviceName,
			"type":    "backend",
		}
		newrelicApp, err := newrelic.NewApplication(
			newrelic.ConfigAppName(p.serviceName),
			newrelic.ConfigLicense(p.newRelicLicense),
			newrelic.ConfigAppLogMetricsEnabled(true),
			newrelic.ConfigAppLogForwardingEnabled(true),
			newrelic.ConfigCustomInsightsEventsEnabled(true),
			newrelic.ConfigAppLogEnabled(true),
			newrelic.ConfigAppLogForwardingEnabled(true),
			func(cfg *newrelic.Config) {
				cfg.HostDisplayName = p.serviceName
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
		CaptureMode:  komoju.CaptureModeManual,
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
	adminWebURL, err := url.Parse(a.AminWebURL)
	if err != nil {
		return fmt.Errorf("cmd: failed to parse admin web url: %w", err)
	}
	p.adminWebURL = adminWebURL
	userWebURL, err := url.Parse(a.UserWebURL)
	if err != nil {
		return fmt.Errorf("cmd: failed to parse user web url: %w", err)
	}
	p.userWebURL = userWebURL

	// LIFFの設定
	liffVerifier, err := auth.NewLIFFVerifier(ctx)
	if err != nil {
		return fmt.Errorf("cmd: failed to new liff verifier: %w", err)
	}
	p.liffVerifier = liffVerifier

	// JWTの設定
	jwtVerifierParams := &auth.JWTVerifierParams{
		Cache:      p.cache,
		Issuer:     a.JWTIssuer,
		PrivateKey: []byte(p.jwtSecret),
	}
	jwtVerifier, err := auth.NewJWTVerifier(jwtVerifierParams)
	if err != nil {
		return fmt.Errorf("cmd: failed to new jwt verifier: %w", err)
	}
	p.jwtVerifier = jwtVerifier
	jwtGeneratorParams := &auth.JWTGeneratorParams{
		Cache:      p.cache,
		Issuer:     a.JWTIssuer,
		PrivateKey: []byte(p.jwtSecret),
	}
	jwtGenerator, err := auth.NewJWTGenerator(jwtGeneratorParams)
	if err != nil {
		return fmt.Errorf("cmd: failed to new jwt generator: %w", err)
	}
	p.jwtGenerator = jwtGenerator

	return nil
}
