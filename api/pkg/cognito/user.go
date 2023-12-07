package cognito

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	cognito "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
)

type SignUpParams struct {
	Username    string
	Email       string
	PhoneNumber string
	Password    string
}

type ConfirmForgotPasswordParams struct {
	Username    string
	VerifyCode  string
	NewPassword string
}

type ChangeEmailParams struct {
	AccessToken string
	Username    string
	OldEmail    string
	NewEmail    string
}

type ConfirmChangeEmailParams struct {
	AccessToken string
	Username    string
	VerifyCode  string
}

type ChangePasswordParams struct {
	AccessToken string
	OldPassword string
	NewPassword string
}

func (c *client) SignUp(ctx context.Context, params *SignUpParams) error {
	in := &cognito.SignUpInput{
		ClientId: c.appClientID,
		Username: aws.String(params.Username),
		Password: aws.String(params.Password),
		UserAttributes: []types.AttributeType{
			{
				Name:  emailField,
				Value: aws.String(params.Email),
			},
			{
				Name:  phoneNumberField,
				Value: aws.String(params.PhoneNumber),
			},
		},
	}
	_, err := c.cognito.SignUp(ctx, in)
	return c.authError(err)
}

func (c *client) ConfirmSignUp(ctx context.Context, username, verifyCode string) error {
	confirmIn := &cognito.ConfirmSignUpInput{
		ClientId:         c.appClientID,
		Username:         aws.String(username),
		ConfirmationCode: aws.String(verifyCode),
	}
	_, err := c.cognito.ConfirmSignUp(ctx, confirmIn)
	if err != nil {
		return c.authError(err)
	}
	updateIn := &cognito.AdminUpdateUserAttributesInput{
		UserPoolId: c.userPoolID,
		Username:   aws.String(username),
		UserAttributes: []types.AttributeType{
			{
				Name:  emailVerifiedField,
				Value: aws.String("true"),
			},
			{
				Name:  phoneNumberVerifiedField,
				Value: aws.String("true"),
			},
		},
	}
	_, err = c.cognito.AdminUpdateUserAttributes(ctx, updateIn)
	return c.authError(err)
}

func (c *client) ResendSignUpCode(ctx context.Context, username string) error {
	in := &cognito.ResendConfirmationCodeInput{
		ClientId: c.appClientID,
		Username: aws.String(username),
	}
	_, err := c.cognito.ResendConfirmationCode(ctx, in)
	return c.authError(err)
}

func (c *client) ForgotPassword(ctx context.Context, username string) error {
	in := &cognito.ForgotPasswordInput{
		ClientId: c.appClientID,
		Username: aws.String(username),
	}
	_, err := c.cognito.ForgotPassword(ctx, in)
	return c.authError(err)
}

func (c *client) ConfirmForgotPassword(ctx context.Context, params *ConfirmForgotPasswordParams) error {
	in := &cognito.ConfirmForgotPasswordInput{
		ClientId:         c.appClientID,
		Username:         aws.String(params.Username),
		ConfirmationCode: aws.String(params.VerifyCode),
		Password:         aws.String(params.NewPassword),
	}
	_, err := c.cognito.ConfirmForgotPassword(ctx, in)
	return c.authError(err)
}

func (c *client) ChangeEmail(ctx context.Context, params *ChangeEmailParams) error {
	// 新しいメールアドレスへの変更と運用サポート対応用に、古いメールアドレスは一時的にカスタムフィールドへ追加
	changeEmailIn := &cognito.UpdateUserAttributesInput{
		AccessToken: aws.String(params.AccessToken),
		UserAttributes: []types.AttributeType{
			{
				Name:  emailField,
				Value: aws.String(params.NewEmail),
			},
		},
	}
	_, err := c.cognito.UpdateUserAttributes(ctx, changeEmailIn)
	if err != nil {
		return c.authError(err)
	}
	// 検証コードを確認するまでの間、今までのメールアドレスが使えなくなってしまうための対応
	acceptingEmailIn := &cognito.AdminUpdateUserAttributesInput{
		UserPoolId: c.userPoolID,
		Username:   aws.String(params.Username),
		UserAttributes: []types.AttributeType{
			{
				Name:  emailField,
				Value: aws.String(params.OldEmail),
			},
			{
				Name:  emailRequestedField,
				Value: aws.String(params.NewEmail),
			},
			{
				Name:  emailVerifiedField,
				Value: aws.String("true"),
			},
		},
	}
	_, err = c.cognito.AdminUpdateUserAttributes(ctx, acceptingEmailIn)
	return c.authError(err)
}

func (c *client) ConfirmChangeEmail(ctx context.Context, params *ConfirmChangeEmailParams) (string, error) {
	username := aws.String(params.Username)
	// 新しいメールアドレス情報の取得
	userIn := &cognito.AdminGetUserInput{
		UserPoolId: c.userPoolID,
		Username:   username,
	}
	out, err := c.cognito.AdminGetUser(ctx, userIn)
	if err != nil {
		return "", c.authError(err)
	}
	var email string
	for i := range out.UserAttributes {
		if aws.ToString(out.UserAttributes[i].Name) == *emailRequestedField {
			email = aws.ToString(out.UserAttributes[i].Value)
			break
		}
	}
	if email == "" {
		return "", fmt.Errorf("%w: %s", ErrNotFound, errNotFoundEmail.Error())
	}
	// コードの検証
	verifyIn := &cognito.VerifyUserAttributeInput{
		AccessToken:   aws.String(params.AccessToken),
		AttributeName: emailField,
		Code:          aws.String(params.VerifyCode),
	}
	_, err = c.cognito.VerifyUserAttribute(ctx, verifyIn)
	if err != nil {
		return "", c.authError(err)
	}
	// メールアドレス情報の更新
	updateIn := &cognito.AdminUpdateUserAttributesInput{
		UserPoolId: c.userPoolID,
		Username:   username,
		UserAttributes: []types.AttributeType{
			{
				Name:  emailField,
				Value: aws.String(email),
			},
			{
				Name:  emailVerifiedField,
				Value: aws.String("true"),
			},
			{
				Name:  emailRequestedField,
				Value: aws.String(""),
			},
		},
	}
	_, err = c.cognito.AdminUpdateUserAttributes(ctx, updateIn)
	return email, c.authError(err)
}

func (c *client) ChangePassword(ctx context.Context, params *ChangePasswordParams) error {
	in := &cognito.ChangePasswordInput{
		AccessToken:      aws.String(params.AccessToken),
		PreviousPassword: aws.String(params.OldPassword),
		ProposedPassword: aws.String(params.NewPassword),
	}
	_, err := c.cognito.ChangePassword(ctx, in)
	return c.authError(err)
}

func (c *client) DeleteUser(ctx context.Context, username string) error {
	in := &cognito.AdminDeleteUserInput{
		UserPoolId: c.userPoolID,
		Username:   aws.String(username),
	}
	_, err := c.cognito.AdminDeleteUser(ctx, in)
	return c.authError(err)
}
