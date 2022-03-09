//go:generate mockgen -source=$GOFILE -package mock_$GOPACKAGE -destination=./../../mock/pkg/$GOPACKAGE/$GOFILE
package cognito

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	cognito "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
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
}

var (
	emailField               = aws.String("email")
	emailVerifiedField       = aws.String("email_verified")
	emailRequestedField      = aws.String("custom:requested_email")
	phoneNumberField         = aws.String("phone_number")
	phoneNumberVerifiedField = aws.String("phone_number_verified")
)

type Params struct {
	UserPoolID      string
	AppClientID     string
	AppClientSecret string
}

type client struct {
	cognito         *cognito.Client
	userPoolID      *string
	appClientID     *string
	appClientSecret *string
}

func NewClient(cfg aws.Config, params *Params) Client {
	return &client{
		cognito:         cognito.NewFromConfig(cfg),
		userPoolID:      aws.String(params.UserPoolID),
		appClientID:     aws.String(params.AppClientID),
		appClientSecret: aws.String(params.AppClientSecret),
	}
}
