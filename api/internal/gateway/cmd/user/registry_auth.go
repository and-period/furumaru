package user

import (
	"github.com/and-period/furumaru/api/pkg/cognito"
)

func (a *app) injectAuth(p *params) {
	// Amazon Cognitoの設定
	userAuthParams := &cognito.Params{
		UserPoolID:  a.CognitoUserPoolID,
		AppClientID: a.CognitoUserClientID,
		AuthDomain:  a.CognitoUserAuthDomain,
	}
	p.userAuth = cognito.NewClient(p.aws, userAuthParams)
}
