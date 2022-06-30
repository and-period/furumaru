//go:generate mockgen -source=$GOFILE -package mock_$GOPACKAGE -destination=./../../mock/pkg/$GOPACKAGE/$GOFILE
package secret

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

type Client interface {
	// シークレット名を指定して値を取得
	Get(ctx context.Context, name string) (map[string]string, error)
}

type client struct {
	secret *secretsmanager.Client
}

func NewClient(cfg aws.Config) Client {
	cli := secretsmanager.NewFromConfig(cfg)
	return &client{
		secret: cli,
	}
}

func (c *client) Get(ctx context.Context, name string) (map[string]string, error) {
	in := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(name),
		VersionStage: aws.String("AWSCURRENT"),
	}
	out, err := c.secret.GetSecretValue(ctx, in)
	if err != nil {
		return nil, err
	}
	secrets := aws.ToString(out.SecretString)
	res := map[string]string{}
	if err := json.Unmarshal([]byte(secrets), &res); err != nil {
		return nil, err
	}
	return res, nil
}
