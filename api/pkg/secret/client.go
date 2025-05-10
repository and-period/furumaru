//go:generate go tool mockgen -source=$GOFILE -package mock_$GOPACKAGE -destination=./../../mock/pkg/$GOPACKAGE/$GOFILE
package secret

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

var errConvertString = errors.New("secret: failed to convert to string")

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
	results := make(map[string]interface{})
	if err := json.Unmarshal([]byte(secrets), &results); err != nil {
		return nil, err
	}
	res := make(map[string]string, len(results))
	for key, value := range results {
		str, err := toString(value)
		if err != nil {
			return nil, err
		}
		res[key] = str
	}
	return res, nil
}

func toString(in interface{}) (string, error) {
	switch v := in.(type) {
	case int:
		return strconv.FormatInt(int64(v), 10), nil
	case int32:
		return strconv.FormatInt(int64(v), 10), nil
	case int64:
		return strconv.FormatInt(v, 10), nil
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 32), nil
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64), nil
	case bool:
		return strconv.FormatBool(v), nil
	case string:
		return v, nil
	default:
		return "", errConvertString
	}
}
