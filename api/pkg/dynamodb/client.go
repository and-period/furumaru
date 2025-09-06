//go:generate go tool mockgen -source=$GOFILE -package mock_$GOPACKAGE -destination=./../../mock/pkg/$GOPACKAGE/$GOFILE
package dynamodb

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/and-period/furumaru/api/pkg/log"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/retry"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

const defaultDelimiter = "-"

var (
	ErrInvalidArgument   = errors.New("dynamodb: invalid argument")
	ErrNotFound          = errors.New("dynamodb: not found")
	ErrAlreadyExists     = errors.New("dynamodb: already exists")
	ErrInternal          = errors.New("dynamodb: internal")
	ErrAborted           = errors.New("dynamodb: aborted")
	ErrCanceled          = errors.New("dynamodb: canceled")
	ErrResourceExhausted = errors.New("dynamodb: resource exhausted")
	ErrOutOfRange        = errors.New("dynamodb: out of range")
	ErrUnknown           = errors.New("dynamodb: unknown")
	ErrTimeout           = errors.New("dynamodb: timeout")
)

type Client interface {
	Count(ctx context.Context, entity Entity) (int64, error)
	Get(ctx context.Context, entity Entity) error
	Insert(ctx context.Context, entity Entity) error
	Delete(ctx context.Context, entity Entity) error
	Scan(ctx context.Context, entities Entities, filter map[string]any) error
	BatchDelete(ctx context.Context, entities Entities) error
}

type Entity interface {
	TableName() string
	PrimaryKey() map[string]interface{}
}

type Entities interface {
	TableName() string
	Items() []Entity
}

type Params struct {
	TablePrefix string
	TableSuffix string
}

type client struct {
	db          *dynamodb.Client
	tablePrefix string
	tableSuffix string
	delimiter   string
}

type options struct {
	maxRetries int
	interval   time.Duration
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

func WithDelimiter(delim string) Option {
	return func(opts *options) {
		opts.delimiter = delim
	}
}

func NewClient(cfg aws.Config, params *Params, opts ...Option) Client {
	dopts := &options{
		maxRetries: retry.DefaultMaxAttempts,
		interval:   retry.DefaultMaxBackoff,
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
		return 0, c.dbError(err)
	}
	return int64(out.Count), nil
}

func (c *client) Get(ctx context.Context, e Entity) error {
	key, err := c.keys(e.PrimaryKey())
	if err != nil {
		return c.dbError(err)
	}
	in := &dynamodb.GetItemInput{
		TableName: c.tableName(e),
		Key:       key,
	}
	out, err := c.db.GetItem(ctx, in)
	if err != nil {
		return c.dbError(err)
	}
	err = attributevalue.UnmarshalMap(out.Item, &e)
	return c.dbError(err)
}

func (c *client) Insert(ctx context.Context, e Entity) error {
	item, err := attributevalue.MarshalMap(e)
	if err != nil {
		return c.dbError(err)
	}
	in := &dynamodb.PutItemInput{
		TableName: c.tableName(e),
		Item:      item,
	}
	_, err = c.db.PutItem(ctx, in)
	return c.dbError(err)
}

func (c *client) Delete(ctx context.Context, e Entity) error {
	key, err := c.keys(e.PrimaryKey())
	if err != nil {
		return c.dbError(err)
	}
	in := &dynamodb.DeleteItemInput{
		TableName: c.tableName(e),
		Key:       key,
	}
	_, err = c.db.DeleteItem(ctx, in)
	return c.dbError(err)
}

func (c *client) Scan(ctx context.Context, entities Entities, filter map[string]any) error {
	if len(filter) == 0 {
		return fmt.Errorf("%w: scan without filter is not allowed", ErrInvalidArgument)
	}

	conditions := make([]string, 0, len(filter))
	expressionAttributeNames := make(map[string]string, len(filter))
	expressionAttributeValues := make(map[string]types.AttributeValue, len(filter))
	for k, v := range filter {
		nameKey := fmt.Sprintf("#%s", k)
		valueKey := fmt.Sprintf(":%s", k)

		expressionAttributeNames[nameKey] = k
		av, err := attributevalue.Marshal(v)
		if err != nil {
			return c.dbError(err)
		}
		expressionAttributeValues[valueKey] = av

		conditions = append(conditions, fmt.Sprintf("%s = %s", nameKey, valueKey))
	}

	in := &dynamodb.ScanInput{
		TableName:                 c.tableName(entities),
		FilterExpression:          aws.String(strings.Join(conditions, " AND ")),
		ExpressionAttributeNames:  expressionAttributeNames,
		ExpressionAttributeValues: expressionAttributeValues,
	}
	result, err := c.db.Scan(ctx, in)
	if err != nil {
		return c.dbError(err)
	}

	err = attributevalue.UnmarshalListOfMaps(result.Items, entities)
	return c.dbError(err)
}

func (c *client) BatchDelete(ctx context.Context, entities Entities) error {
	const batchSize = 25 // DynamoDBのBatchWriteItemは最大25アイテムまで

	table := c.tableName(entities)
	items := entities.Items()

	for i := 0; i < len(items); i += batchSize {
		end := i + batchSize
		if end > len(items) {
			end = len(items)
		}

		reqs := make([]types.WriteRequest, 0, end-i)
		for _, item := range items[i:end] {
			key, err := c.keys(item.PrimaryKey())
			if err != nil {
				return c.dbError(err)
			}
			req := types.WriteRequest{
				DeleteRequest: &types.DeleteRequest{
					Key: key,
				},
			}
			reqs = append(reqs, req)
		}

		in := &dynamodb.BatchWriteItemInput{
			RequestItems: map[string][]types.WriteRequest{
				aws.ToString(table): reqs,
			},
		}

		if _, err := c.db.BatchWriteItem(ctx, in); err != nil {
			return c.dbError(err)
		}
	}

	return nil
}

func (c *client) tableName(e interface{ TableName() string }) *string {
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
			return nil, c.dbError(err)
		}
		res[k] = item
	}
	return res, nil
}

func (c *client) dbError(err error) error {
	if err == nil {
		return nil
	}
	slog.Debug("Failed to dynamodb api", log.Error(err))

	switch {
	case errors.Is(err, context.Canceled):
		return fmt.Errorf("%w: %s", ErrCanceled, err.Error())
	case errors.Is(err, context.DeadlineExceeded):
		return fmt.Errorf("%w: %s", ErrTimeout, err.Error())
	}

	var (
		cfe *types.ConditionalCheckFailedException
		die *types.DuplicateItemException
		pme *types.IdempotentParameterMismatchException
		ine *types.IndexNotFoundException
		see *types.ItemCollectionSizeLimitExceededException
		lee *types.LimitExceededException
		tee *types.ProvisionedThroughputExceededException
		rue *types.ResourceInUseException
		rne *types.ResourceNotFoundException
		tue *types.TableInUseException
		tae *types.TransactionCanceledException
		tce *types.TransactionConflictException
		tpe *types.TransactionInProgressException
	)

	switch {
	case errors.As(err, &cfe), errors.As(err, &pme):
		return fmt.Errorf("%w: %s", ErrInvalidArgument, err.Error())
	case errors.As(err, &rne):
		return fmt.Errorf("%w: %s", ErrNotFound, err.Error())
	case errors.As(err, &die):
		return fmt.Errorf("%w: %s", ErrAlreadyExists, err.Error())
	case errors.As(err, &ine):
		return fmt.Errorf("%w: %s", ErrInternal, err.Error())
	case errors.As(err, &see), errors.As(err, &lee):
		return fmt.Errorf("%w: %s", ErrOutOfRange, err.Error())
	case errors.As(err, &tee):
		return fmt.Errorf("%w: %s", ErrResourceExhausted, err.Error())
	case errors.As(err, &rue), errors.As(err, &tue), errors.As(err, &tae), errors.As(err, &tce), errors.As(err, &tpe):
		return fmt.Errorf("%w: %s", ErrAborted, err.Error())
	default:
		return fmt.Errorf("%w: %s", ErrUnknown, err.Error())
	}
}

func IsRetryable(err error) bool {
	switch {
	case errors.Is(err, ErrCanceled),
		errors.Is(err, ErrTimeout),
		errors.Is(err, ErrInternal),
		errors.Is(err, ErrResourceExhausted),
		errors.Is(err, ErrAborted):
		return true
	default:
		return false
	}
}
