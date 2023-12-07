//go:generate mockgen -source=$GOFILE -package mock_$GOPACKAGE -destination=./../../mock/pkg/$GOPACKAGE/$GOFILE
package cognito

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/retry"
	cognito "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"go.uber.org/zap"
)

type Client interface {
	// #############################################
	// 認証関連
	// #############################################
	// サインイン
	SignIn(ctx context.Context, username, password string) (*AuthResult, error)
	// サインアウト (アクセストークン使用)
	SignOut(ctx context.Context, accessToken string) error
	// ユーザー情報取得 (アクセストークン使用)
	GetUser(ctx context.Context, accessToken string) (*AuthUser, error)
	// ユーザーID取得 (アクセストークン使用)
	GetUsername(ctx context.Context, accessToken string) (string, error)
	// トークンの更新 (更新トークン使用)
	RefreshToken(ctx context.Context, refreshToken string) (*AuthResult, error)

	// #############################################
	// ユーザー関連
	// #############################################
	// ユーザー登録
	SignUp(ctx context.Context, params *SignUpParams) error
	// ユーザー登録 (コード検証)
	ConfirmSignUp(ctx context.Context, username, verifyCode string) error
	// ユーザー登録 (コードの再送)
	ResendSignUpCode(ctx context.Context, username string) error
	// パスワードリセット
	ForgotPassword(ctx context.Context, username string) error
	// パスワードリセット (コード検証)
	ConfirmForgotPassword(ctx context.Context, params *ConfirmForgotPasswordParams) error
	// メールアドレス更新
	ChangeEmail(ctx context.Context, params *ChangeEmailParams) error
	// メールアドレス変更 (コード検証)
	ConfirmChangeEmail(ctx context.Context, params *ConfirmChangeEmailParams) (string, error)
	// パスワード更新
	ChangePassword(ctx context.Context, params *ChangePasswordParams) error
	// ユーザー削除
	DeleteUser(ctx context.Context, username string) error

	// #############################################
	// ユーザー関連 (管理者)
	// #############################################
	// ユーザー登録
	AdminCreateUser(ctx context.Context, params *AdminCreateUserParams) error
	// メールアドレス更新
	AdminChangeEmail(ctx context.Context, params *AdminChangeEmailParams) error
	// パスワード更新
	AdminChangePassword(ctx context.Context, params *AdminChangePasswordParams) error
}

var (
	emailField               = aws.String("email")
	emailVerifiedField       = aws.String("email_verified")
	emailRequestedField      = aws.String("custom:requested_email")
	phoneNumberField         = aws.String("phone_number")
	phoneNumberVerifiedField = aws.String("phone_number_verified")
)

var (
	ErrInvalidArgument   = errors.New("cognito: invalid argument")
	ErrUnauthenticated   = errors.New("cognito: unauthenticated")
	ErrNotFound          = errors.New("cognito: not found")
	ErrAlreadyExists     = errors.New("cognito: already exists")
	ErrInternal          = errors.New("cognito: internal")
	ErrCanceled          = errors.New("cognito: canceled")
	ErrResourceExhausted = errors.New("cognito: resource exhausted")
	ErrUnknown           = errors.New("cognito: unknown")
	ErrTimeout           = errors.New("cognito: timeout")
	errNotFoundEmail     = errors.New("cognito: not found requested email")
)

type Params struct {
	UserPoolID      string
	AppClientID     string
	AppClientSecret string
}

type client struct {
	cognito         *cognito.Client
	logger          *zap.Logger
	userPoolID      *string
	appClientID     *string
	appClientSecret *string
}

type options struct {
	maxRetries int
	interval   time.Duration
	logger     *zap.Logger
}

type Option func(*options)

func WithMaxRetries(maxRetries int) Option {
	return func(opts *options) {
		opts.maxRetries = maxRetries
	}
}

func WithInterval(interval time.Duration) Option {
	return func(opts *options) {
		opts.interval = interval
	}
}

func WithLogger(logger *zap.Logger) Option {
	return func(opts *options) {
		opts.logger = logger
	}
}

func NewClient(cfg aws.Config, params *Params, opts ...Option) Client {
	dopts := &options{
		maxRetries: retry.DefaultMaxAttempts,
		interval:   retry.DefaultMaxBackoff,
		logger:     zap.NewNop(),
	}
	for i := range opts {
		opts[i](dopts)
	}
	cli := cognito.NewFromConfig(cfg, func(o *cognito.Options) {
		o.Retryer = retry.NewStandard(func(o *retry.StandardOptions) {
			o.MaxAttempts = dopts.maxRetries
			o.MaxBackoff = dopts.interval
		})
	})
	return &client{
		cognito:         cli,
		userPoolID:      aws.String(params.UserPoolID),
		appClientID:     aws.String(params.AppClientID),
		appClientSecret: aws.String(params.AppClientSecret),
		logger:          dopts.logger,
	}
}

func (c *client) authError(err error) error {
	if err == nil {
		return nil
	}
	c.logger.Debug("Failed to cognito api", zap.Error(err))

	switch {
	case errors.Is(err, context.Canceled):
		return fmt.Errorf("%w: %s", ErrCanceled, err.Error())
	case errors.Is(err, context.DeadlineExceeded):
		return fmt.Errorf("%w: %s", ErrTimeout, err.Error())
	}

	var (
		aee *types.AliasExistsException
		cfe *types.CodeDeliveryFailureException
		cme *types.CodeMismatchException
		ece *types.ExpiredCodeException
		iee *types.InternalErrorException
		ipe *types.InvalidParameterException
		lee *types.LimitExceededException
		nae *types.NotAuthorizedException
		pre *types.PasswordResetRequiredException
		rne *types.ResourceNotFoundException
		tfe *types.TooManyFailedAttemptsException
		tre *types.TooManyRequestsException
		uce *types.UserNotConfirmedException
		uee *types.UsernameExistsException
		une *types.UserNotFoundException
	)

	switch {
	case errors.As(err, &cme), errors.As(err, &ipe):
		return fmt.Errorf("%w: %s", ErrInvalidArgument, err.Error())
	case errors.As(err, &ece), errors.As(err, &nae), errors.As(err, &pre), errors.As(err, &uce):
		return fmt.Errorf("%w: %s", ErrUnauthenticated, err.Error())
	case errors.As(err, &rne), errors.As(err, &une):
		return fmt.Errorf("%w: %s", ErrNotFound, err.Error())
	case errors.As(err, &aee), errors.As(err, &uee):
		return fmt.Errorf("%w: %s", ErrAlreadyExists, err.Error())
	case errors.As(err, &lee), errors.As(err, &tfe), errors.As(err, &tre):
		return fmt.Errorf("%w: %s", ErrResourceExhausted, err.Error())
	case errors.As(err, &cfe), errors.As(err, &iee):
		return fmt.Errorf("%w: %s", ErrInternal, err.Error())
	default:
		return fmt.Errorf("%w: %s", ErrUnknown, err.Error())
	}
}
