package cognito

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	cognito "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
)

type AdminCreateUserParams struct {
	Username string
	Email    string
	Password string
}

type AdminChangeEmailParams struct {
	Username string
	Email    string
}

type AdminChangePasswordParams struct {
	Username  string
	Password  string
	Permanent bool
}

func (c *client) AdminCreateUser(ctx context.Context, params *AdminCreateUserParams) error {
	in := &cognito.AdminCreateUserInput{
		UserPoolId: c.userPoolID,
		Username:   aws.String(params.Username),
		UserAttributes: []types.AttributeType{
			{
				Name:  emailField,
				Value: aws.String(params.Email),
			},
		},
		MessageAction: types.MessageActionTypeResend,
	}
	if params.Password == "" {
		// 一時的なパスワードを付与し、メール通知 (初回ログイン時にパスワード変更要求)
		_, err := c.cognito.AdminCreateUser(ctx, in)
		return c.authError(err)
	}
	// 恒久的なパスワードを付与 (未通知、かつ初回ログイン時のパスワード変更要求も不要)
	attr := types.AttributeType{
		Name:  emailVerifiedField,
		Value: aws.String("true"),
	}
	in.TemporaryPassword = aws.String(params.Password)
	in.MessageAction = types.MessageActionTypeSuppress
	in.UserAttributes = append(in.UserAttributes, attr)
	if _, err := c.cognito.AdminCreateUser(ctx, in); err != nil {
		return c.authError(err)
	}
	passIn := &AdminChangePasswordParams{
		Username:  params.Username,
		Password:  params.Password,
		Permanent: true,
	}
	err := c.AdminChangePassword(ctx, passIn)
	return c.authError(err)
}

func (c *client) AdminChangeEmail(ctx context.Context, params *AdminChangeEmailParams) error {
	in := &cognito.AdminUpdateUserAttributesInput{
		UserPoolId: c.userPoolID,
		Username:   aws.String(params.Username),
		UserAttributes: []types.AttributeType{
			{
				Name:  emailField,
				Value: aws.String(params.Email),
			},
			{
				Name:  emailVerifiedField,
				Value: aws.String("true"),
			},
		},
	}
	_, err := c.cognito.AdminUpdateUserAttributes(ctx, in)
	return c.authError(err)
}

func (c *client) AdminVerifyEmail(ctx context.Context, username string) error {
	in := &cognito.AdminUpdateUserAttributesInput{
		UserPoolId: c.userPoolID,
		Username:   aws.String(username),
		UserAttributes: []types.AttributeType{
			{
				Name:  emailVerifiedField,
				Value: aws.String("true"),
			},
		},
	}
	_, err := c.cognito.AdminUpdateUserAttributes(ctx, in)
	return c.authError(err)
}

func (c *client) AdminChangePassword(ctx context.Context, params *AdminChangePasswordParams) error {
	in := &cognito.AdminSetUserPasswordInput{
		UserPoolId: c.userPoolID,
		Username:   aws.String(params.Username),
		Password:   aws.String(params.Password),
		Permanent:  *aws.Bool(params.Permanent),
	}
	_, err := c.cognito.AdminSetUserPassword(ctx, in)
	return c.authError(err)
}
