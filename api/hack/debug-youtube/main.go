package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/and-period/furumaru/api/pkg/secret"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"go.uber.org/zap"
	"golang.org/x/oauth2/google/externalaccount"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

func main() {
	start := time.Now()
	fmt.Println("Start..")
	if err := run(); err != nil {
		panic(err)
	}
	fmt.Printf("Done: %s\n", time.Since(start))
}

// Config is the configuration for the external account token source.
type Config struct {
	// Audience is the Secure Token Service (STS) audience which contains the resource name for the workload
	// identity pool or the workforce pool and the provider identifier in that pool. Required.
	Audience string `json:"audience"`
	// SubjectTokenType is the STS token type based on the Oauth2.0 token exchange spec.
	// Expected values include:
	// “urn:ietf:params:oauth:token-type:jwt”
	// “urn:ietf:params:oauth:token-type:id-token”
	// “urn:ietf:params:oauth:token-type:saml2”
	// “urn:ietf:params:aws:token-type:aws4_request”
	// Required.
	SubjectTokenType string `json:"subject_token_type"`
	// TokenURL is the STS token exchange endpoint. If not provided, will default to
	// https://sts.UNIVERSE_DOMAIN/v1/token, with UNIVERSE_DOMAIN set to the
	// default service domain googleapis.com unless UniverseDomain is set.
	// Optional.
	TokenURL string `json:"token_url"`
	// TokenInfoURL is the token_info endpoint used to retrieve the account related information (
	// user attributes like account identifier, eg. email, username, uid, etc). This is
	// needed for gCloud session account identification. Optional.
	TokenInfoURL string `json:"token_info_url"`
	// ServiceAccountImpersonationURL is the URL for the service account impersonation request. This is only
	// required for workload identity pools when APIs to be accessed have not integrated with UberMint. Optional.
	ServiceAccountImpersonationURL string `json:"service_account_impersonation_url"`
	// ServiceAccountImpersonationLifetimeSeconds is the number of seconds the service account impersonation
	// token will be valid for. If not provided, it will default to 3600. Optional.
	ServiceAccountImpersonationLifetimeSeconds int `json:"service_account_impersonation_lifetime_seconds"`
	// ClientSecret is currently only required if token_info endpoint also
	// needs to be called with the generated GCP access token. When provided, STS will be
	// called with additional basic authentication using ClientId as username and ClientSecret as password. Optional.
	ClientSecret string `json:"client_secret"`
	// ClientID is only required in conjunction with ClientSecret, as described above. Optional.
	ClientID string `json:"client_id"`
	// CredentialSource contains the necessary information to retrieve the token itself, as well
	// as some environmental information. One of SubjectTokenSupplier, AWSSecurityCredentialSupplier or
	// CredentialSource must be provided. Optional.
	CredentialSource *externalaccount.CredentialSource `json:"credential_source"`
	// QuotaProjectID is injected by gCloud. If the value is non-empty, the Auth libraries
	// will set the x-goog-user-project header which overrides the project associated with the credentials. Optional.
	QuotaProjectID string `json:"quota_project_id"`
	// Scopes contains the desired scopes for the returned access token. Optional.
	Scopes []string `json:"scopes"`
	// WorkforcePoolUserProject is the workforce pool user project number when the credential
	// corresponds to a workforce pool and not a workload identity pool.
	// The underlying principal must still have serviceusage.services.use IAM
	// permission to use the project for billing/quota. Optional.
	WorkforcePoolUserProject string `json:"workforce_pool_user_project"`
	// SubjectTokenSupplier is an optional token supplier for OIDC/SAML credentials.
	// One of SubjectTokenSupplier, AWSSecurityCredentialSupplier or CredentialSource must be provided. Optional.
	SubjectTokenSupplier externalaccount.SubjectTokenSupplier `json:"subject_token_supplier"`
	// AwsSecurityCredentialsSupplier is an AWS Security Credential supplier for AWS credentials.
	// One of SubjectTokenSupplier, AWSSecurityCredentialSupplier or CredentialSource must be provided. Optional.
	AwsSecurityCredentialsSupplier externalaccount.AwsSecurityCredentialsSupplier `json:"aws_security_credentials_supplier"`
	// UniverseDomain is the default service domain for a given Cloud universe.
	// This value will be used in the default STS token URL. The default value
	// is "googleapis.com". It will not be used if TokenURL is set. Optional.
	UniverseDomain string
}

func run() error {
	logger, err := zap.NewDevelopment()
	if err != nil {
		return err
	}
	ctx := context.Background()

	awscfg, err := awsconfig.LoadDefaultConfig(ctx, awsconfig.WithRegion("ap-northeast-1"))
	if err != nil {
		return fmt.Errorf("cmd: failed to load aws config: %w", err)
	}
	secret := secret.NewClient(awscfg)
	secrets, err := secret.Get(ctx, "furumaru-gcp-workload-identity-stg")
	if err != nil {
		return fmt.Errorf("cmd: failed to get secret: %w", err)
	}
	logger.Debug("Succeeded to get secret", zap.Any("secrets", secrets))

	config := &Config{}
	if err := json.Unmarshal([]byte(secrets["credential"]), config); err != nil {
		return fmt.Errorf("cmd: failed to unmarshal config: %w", err)
	}
	logger.Debug("Succeeded to unmarshal config", zap.Any("config", config))

	token, err := externalaccount.NewTokenSource(ctx, externalaccount.Config(*config))
	if err != nil {
		return fmt.Errorf("cmd: failed to get token: %w", err)
	}
	service, err := youtube.NewService(ctx, option.WithTokenSource(token))
	if err != nil {
		return fmt.Errorf("cmd: failed to get service: %w", err)
	}
	logger.Debug("Succeeded to get service", zap.Any("service", service))

	res, err := service.Channels.List([]string{"id", "snippet"}).Mine(true).Context(ctx).Do()
	if err != nil {
		return fmt.Errorf("cmd: failed to get channels: %w", err)
	}
	logger.Info("Succeeded to get channels", zap.Any("res", res))
	return nil
}
