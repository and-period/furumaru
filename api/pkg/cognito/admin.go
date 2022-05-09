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
	}
	_, err := c.cognito.AdminCreateUser(ctx, in)
	return authError(err)
}
