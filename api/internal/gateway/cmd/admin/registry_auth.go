package admin

import (
	"github.com/and-period/furumaru/api/pkg/cognito"
)

func (a *app) injectAuth(p *params) {
	// Amazon Cognitoの設定
	adminAuthParams := &cognito.Params{
		UserPoolID:  a.CognitoAdminPoolID,
		AppClientID: a.CognitoAdminClientID,
		AuthDomain:  a.CognitoAdminAuthDomain,
	}
	p.adminAuth = cognito.NewClient(p.aws, adminAuthParams)
	userAuthParams := &cognito.Params{
		UserPoolID:  a.CognitoUserPoolID,
		AppClientID: a.CognitoUserClientID,
	}
	p.userAuth = cognito.NewClient(p.aws, userAuthParams)
}
