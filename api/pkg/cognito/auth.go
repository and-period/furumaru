package cognito

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	cognito "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
)

type AuthResult struct {
	IDToken      string
	AccessToken  string
	RefreshToken string
	ExpiresIn    int32
}

type AuthUser struct {
	Username    string
	Email       string
	PhoneNumber string
}

func (c *client) SignIn(ctx context.Context, username, password string) (*AuthResult, error) {
	in := &cognito.InitiateAuthInput{
		ClientId: c.appClientID,
		AuthFlow: types.AuthFlowTypeUserPasswordAuth,
		AuthParameters: map[string]string{
			"USERNAME": username,
			"PASSWORD": password,
		},
	}
	out, err := c.cognito.InitiateAuth(ctx, in)
	if err != nil {
		return nil, c.authError(err)
	}
	auth := &AuthResult{
		IDToken:      aws.ToString(out.AuthenticationResult.IdToken),
		AccessToken:  aws.ToString(out.AuthenticationResult.AccessToken),
		RefreshToken: aws.ToString(out.AuthenticationResult.RefreshToken),
		ExpiresIn:    aws.ToInt32(&out.AuthenticationResult.ExpiresIn),
	}
	return auth, nil
}

func (c *client) SignOut(ctx context.Context, accessToken string) error {
	in := &cognito.GlobalSignOutInput{
		AccessToken: aws.String(accessToken),
	}
	_, err := c.cognito.GlobalSignOut(ctx, in)
	return c.authError(err)
}

func (c *client) GetUser(ctx context.Context, accessToken string) (*AuthUser, error) {
	in := &cognito.GetUserInput{
		AccessToken: aws.String(accessToken),
	}
	out, err := c.cognito.GetUser(ctx, in)
	if err != nil {
		return nil, c.authError(err)
	}
	var email, phoneNumber string
	for i := range out.UserAttributes {
		if aws.ToString(out.UserAttributes[i].Name) == *emailField {
			email = aws.ToString(out.UserAttributes[i].Value)
			continue
		}
		if aws.ToString(out.UserAttributes[i].Name) == *phoneNumberField {
			phoneNumber = aws.ToString(out.UserAttributes[i].Value)
		}
	}
	auth := &AuthUser{
		Username:    aws.ToString(out.Username),
		Email:       email,
		PhoneNumber: phoneNumber,
	}
	return auth, nil
}

func (c *client) GetUsername(ctx context.Context, accessToken string) (string, error) {
	in := &cognito.GetUserInput{
		AccessToken: aws.String(accessToken),
	}
	out, err := c.cognito.GetUser(ctx, in)
	if err != nil {
		return "", c.authError(err)
	}
	return aws.ToString(out.Username), nil
}

func (c *client) RefreshToken(ctx context.Context, refreshToken string) (*AuthResult, error) {
	in := &cognito.InitiateAuthInput{
		ClientId: c.appClientID,
		AuthFlow: types.AuthFlowTypeRefreshTokenAuth,
		AuthParameters: map[string]string{
			"REFRESH_TOKEN": refreshToken,
		},
	}
	out, err := c.cognito.InitiateAuth(ctx, in)
	if err != nil {
		return nil, c.authError(err)
	}
	auth := &AuthResult{
		IDToken:      aws.ToString(out.AuthenticationResult.IdToken),
		AccessToken:  aws.ToString(out.AuthenticationResult.AccessToken),
		RefreshToken: aws.ToString(out.AuthenticationResult.RefreshToken),
		ExpiresIn:    aws.ToInt32(&out.AuthenticationResult.ExpiresIn),
	}
	return auth, nil
}
