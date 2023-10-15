//go:generate mockgen -source=$GOFILE -package mock_$GOPACKAGE -destination=./../../mock/pkg/$GOPACKAGE/$GOFILE
package dynamodb

import (
	"context"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/retry"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"go.uber.org/zap"
)

const defaultDelimiter = "-"

type Client interface {
	Count(ctx context.Context, entity Entity) (int64, error)
	Get(ctx context.Context, primaryKey map[string]interface{}, entity Entity) error
	Insert(ctx context.Context, entity Entity) error
}

type Entity interface {
	TableName() string
	PrimaryKey() map[string]interface{}
}

type Params struct {
	TablePrefix string
	TableSuffix string
}

type client struct {
	db          *dynamodb.Client
	logger      *zap.Logger
	tablePrefix string
	tableSuffix string
	delimiter   string
}

type options struct {
	maxRetries int
	interval   time.Duration
	logger     *zap.Logger
	delimiter  string
}

type Option func(*options)

func WithMaxRetries(maxRetries int) Option {
	return func(opts *options) {
		opts.maxRetries = maxRetries
	}
}

func WithInterval(interval time.Duration) Option {
	return func(opts *options) {
		opts.interval = interval
	}
}

func WithLogger(logger *zap.Logger) Option {
	return func(opts *options) {
		opts.logger = logger
	}
}

func WithDelimiter(delim string) Option {
	return func(opts *options) {
		opts.delimiter = delim
	}
}

func NewClient(cfg aws.Config, params *Params, opts ...Option) Client {
	dopts := &options{
		maxRetries: retry.DefaultMaxAttempts,
		interval:   retry.DefaultMaxBackoff,
		logger:     zap.NewNop(),
		delimiter:  defaultDelimiter,
	}
	for i := range opts {
		opts[i](dopts)
	}
	cli := dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
		o.Retryer = retry.NewStandard(func(o *retry.StandardOptions) {
			o.MaxAttempts = dopts.maxRetries
			o.MaxBackoff = dopts.interval
		})
	})
	return &client{
		db:          cli,
		logger:      dopts.logger,
		tablePrefix: params.TablePrefix,
		tableSuffix: params.TableSuffix,
		delimiter:   dopts.delimiter,
	}
}

func (c *client) Count(ctx context.Context, e Entity) (int64, error) {
	in := &dynamodb.ScanInput{
		TableName: c.tableName(e),
		Select:    types.SelectCount,
	}
	out, err := c.db.Scan(ctx, in)
	if err != nil {
		return 0, err
	}
	return int64(out.Count), nil
}

func (c *client) Get(ctx context.Context, pk map[string]interface{}, e Entity) error {
	key, err := c.keys(pk)
	if err != nil {
		return err
	}
	in := &dynamodb.GetItemInput{
		TableName: c.tableName(e),
		Key:       key,
	}
	out, err := c.db.GetItem(ctx, in)
	if err != nil {
		return err
	}
	err = attributevalue.UnmarshalMap(out.Item, &e)
	return err
}

func (c *client) Insert(ctx context.Context, e Entity) error {
	item, err := attributevalue.MarshalMap(e)
	if err != nil {
		return err
	}
	in := &dynamodb.PutItemInput{
		TableName: c.tableName(e),
		Item:      item,
	}
	_, err = c.db.PutItem(ctx, in)
	return err
}

func (c *client) tableName(e Entity) *string {
	strs := []string{c.tablePrefix, e.TableName()}
	if c.tableSuffix != "" {
		strs = append(strs, c.tableSuffix)
	}
	return aws.String(strings.Join(strs, c.delimiter))
}

func (c *client) keys(keys map[string]interface{}) (map[string]types.AttributeValue, error) {
	res := make(map[string]types.AttributeValue, len(keys))
	for k, v := range keys {
		item, err := attributevalue.Marshal(v)
		if err != nil {
			return nil, err
		}
		res[k] = item
	}
	return res, nil
}
