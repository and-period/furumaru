package cognito

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	cognito "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"go.uber.org/zap"
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

type GenerateAuthURLParams struct {
	State        string
	Nonce        string
	ProviderType ProviderType
	RedirectURI  string
}

type GetAccessTokenParams struct {
	Code        string
	RedirectURI string
}

type LinkProviderParams struct {
	Username     string
	ProviderType ProviderType
	AccountID    string
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

func (c *client) GenerateAuthURL(ctx context.Context, params *GenerateAuthURLParams) (string, error) {
	const format = "https://%s/oauth2/authorize"
	authURL, err := url.Parse(fmt.Sprintf(format, c.authDomain))
	if err != nil {
		return "", fmt.Errorf("cognito: failed to parse request uri: %w", err)
	}

	values := url.Values{}
	values.Add("response_type", "code")                          // 応答形式（固定:code）
	values.Add("client_id", aws.ToString(c.appClientID))         // クライアントID
	values.Add("redirect_uri", params.RedirectURI)               // 応答先URI（WTリダイレクト先URL）
	values.Add("identity_provider", string(params.ProviderType)) // 認証プロバイダー
	values.Add("state", params.State)                            // セキュア文字列（CSRF/XSRF対策）
	values.Add("nonce", params.Nonce)                            // セキュア文字列（リプレイアタック対策）

	authURL.RawQuery = values.Encode()

	// スコープIDはスペース(%20)で連結する必要があるため、別途エスケープ処理をして結合する
	const scope = "openid email aws.cognito.signin.user.admin"
	return fmt.Sprintf("%s&scope=%s", authURL.String(), url.PathEscape(scope)), nil
}

type getAccessTokenResponse struct {
	AccessToken  string `json:"access_token"`  // アクセストークン
	RefreshToken string `json:"refresh_token"` // リフレッシュトークン
	IDToken      string `json:"id_token"`      // IDトークン
	TokenType    string `json:"token_type"`    // トークン形式
	ExpiresIn    int32  `json:"expires_in"`    // トークン有効期限
}

func (c *client) GetAccessToken(ctx context.Context, params *GetAccessTokenParams) (*AuthResult, error) {
	const format = "https://%s/oauth2/token"
	authURL, err := url.Parse(fmt.Sprintf(format, c.authDomain))
	if err != nil {
		return nil, fmt.Errorf("cognito: failed to parse request uri: %w", err)
	}

	values := &url.Values{}
	values.Add("grant_type", "authorization_code") // 権限形式
	values.Add("code", params.Code)                // 認可コード
	values.Add("redirect_uri", params.RedirectURI) // 応答先URI
	values.Add("client_id", aws.ToString(c.appClientID))
	body := strings.NewReader(values.Encode())

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, authURL.String(), body)
	if err != nil {
		return nil, fmt.Errorf("cognito: failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("cognito: failed to request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		c.logger.Error("Failed to get access token", zap.Int("status", res.StatusCode), zap.String("body", string(body)))
		return nil, fmt.Errorf("cognito: failed to get access token: status code=%d", res.StatusCode)
	}

	out := &getAccessTokenResponse{}
	if err := json.NewDecoder(res.Body).Decode(out); err != nil {
		return nil, fmt.Errorf("cognito: failed to decode get access token response: %w", err)
	}

	auth := &AuthResult{
		IDToken:      out.IDToken,
		AccessToken:  out.AccessToken,
		RefreshToken: out.RefreshToken,
		ExpiresIn:    out.ExpiresIn,
	}
	return auth, nil
}

func (c *client) LinkProvider(ctx context.Context, params *LinkProviderParams) error {
	linkIn := &cognito.AdminLinkProviderForUserInput{
		UserPoolId: c.userPoolID,
		DestinationUser: &types.ProviderUserIdentifierType{
			ProviderName:           aws.String(string(ProviderTypeCognito)),
			ProviderAttributeValue: aws.String(params.Username),
		},
		SourceUser: &types.ProviderUserIdentifierType{
			ProviderName:           aws.String(string(params.ProviderType)),
			ProviderAttributeName:  aws.String("Cognito_Subject"),
			ProviderAttributeValue: aws.String(params.AccountID),
		},
	}
	if _, err := c.cognito.AdminLinkProviderForUser(ctx, linkIn); err != nil {
		return c.authError(err)
	}
	attrIn := &cognito.AdminUpdateUserAttributesInput{
		UserPoolId: c.userPoolID,
		Username:   aws.String(params.Username),
		UserAttributes: []types.AttributeType{
			{
				Name:  emailVerifiedField,
				Value: aws.String("true"),
			},
		},
	}
	if _, err := c.cognito.AdminUpdateUserAttributes(ctx, attrIn); err != nil {
		return c.authError(err)
	}
	return nil
}
