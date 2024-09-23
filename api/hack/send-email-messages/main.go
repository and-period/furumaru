//nolint:forbidigo
package main

import (
	"context"
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/and-period/furumaru/api/pkg/backoff"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/mailer"
	"go.uber.org/zap"
)

var (
	sendgridAPIKey  string
	mailTemplateID  string
	mailFromName    string
	mailFromAddress string
	sourceCSVFile   string
	debug           bool
)

type app struct {
	client mailer.Client
	logger *zap.Logger
	source string
	debug  bool
}

func main() {
	startedAt := jst.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	app, err := setup(ctx)
	if err != nil {
		panic(fmt.Errorf("failed to setup: %w", err))
	}

	if err := app.run(ctx); err != nil {
		panic(fmt.Errorf("failed to run: %w", err))
	}

	endAt := jst.Now()

	const format = "2006-01-02 15:04:05"
	fmt.Printf("Done: %s -> %s (%s)\n",
		jst.Format(startedAt, format), jst.Format(endAt, format),
		endAt.Sub(startedAt).Truncate(time.Second).String(),
	)
}

func setup(_ context.Context) (*app, error) {
	flag.StringVar(&sendgridAPIKey, "sendgrid-api-key", "", "SendGridのAPIキー")
	flag.StringVar(&mailTemplateID, "mail-template-id", "d-66e33b4d42a94c3db8e39b10911bca44", "メールテンプレートID")
	flag.StringVar(&mailFromName, "mail-from-name", "", "メール送信元名")
	flag.StringVar(&mailFromAddress, "mail-from-address", "", "メール送信元アドレス")
	flag.StringVar(&sourceCSVFile, "source-csv-file", "", "参照元CSVファイル")
	flag.BoolVar(&debug, "debug", true, "デバッグモード")
	flag.Parse()

	if mailTemplateID == "" {
		return nil, fmt.Errorf("mail-template-id is required")
	}
	if mailFromName == "" {
		return nil, fmt.Errorf("mail-from-name is required")
	}
	if mailFromAddress == "" {
		return nil, fmt.Errorf("mail-from-address is required")
	}
	if sourceCSVFile == "" {
		return nil, fmt.Errorf("source-csv-file is required")
	}

	logger, err := zap.NewDevelopment()
	if err != nil {
		return nil, fmt.Errorf("failed to create logger: %w", err)
	}

	mailerParams := &mailer.Params{
		APIKey:      sendgridAPIKey,
		FromName:    mailFromName,
		FromAddress: mailFromAddress,
		TemplateMap: map[string]string{"default": mailTemplateID},
	}
	mailer := mailer.NewClient(mailerParams, mailer.WithLogger(logger))

	app := &app{
		logger: logger,
		client: mailer,
		source: sourceCSVFile,
		debug:  debug,
	}
	return app, nil
}

func (a *app) run(ctx context.Context) error {
	file, err := os.Open(a.source)
	if err != nil {
		return fmt.Errorf("failed to open source csv file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// CSVファイルのヘッダーを読み込む
	header, err := reader.Read()
	if err != nil {
		return fmt.Errorf("failed to read header: %w", err)
	}

	// CSVファイルの行を読み込む
	rows, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read rows: %w", err)
	}

	// メールの送信先と動的な変数を設定する
	req := make([]*builder, len(rows))
	for i, row := range rows {
		builder := newBuilder(header, row)

		if _, ok := builder.getName(); !ok {
			return fmt.Errorf("name is required: index=%d", i)
		}
		if _, ok := builder.getEmail(); !ok {
			return fmt.Errorf("email is required: index=%d", i)
		}

		req[i] = builder
	}

	output, err := os.Create(fmt.Sprintf("output_%s.txt", jst.Now().Format("20060102150405")))
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer output.Close()

	writer := io.Writer(output)

	// メールを送信する
	a.logger.Info("Start to send emails")
	for i, builder := range req {
		name := builder.substitutions["氏名"]
		sei := builder.substitutions["姓"]
		email := builder.substitutions["メールアドレス"]
		res := "OK"

		retryFn := func() error {
			return a.send(ctx, builder)
		}
		retry := backoff.NewExponentialBackoff(3)
		err := backoff.Retry(ctx, retry, retryFn, backoff.WithRetryablel(a.isRetryable))
		if err != nil {
			res = "NG"
			a.logger.Error("failed to retry", zap.Int("index", i), zap.String("name", name), zap.Error(err))
		}

		out := fmt.Sprintf("index=%d: 氏名=%s, 姓=%s, メールアドレス=%s, 結果=%s\n", i, name, sei, email, res)
		if _, err := fmt.Fprint(writer, out); err != nil {
			a.logger.Error("failed to write output", zap.Int("index", i), zap.String("name", name), zap.Error(err))
		}
	}
	a.logger.Info("Finish to send emails")

	return nil
}

func (a *app) isRetryable(err error) bool {
	if err == nil {
		return false
	}
	switch {
	case errors.Is(err, context.Canceled), errors.Is(err, context.DeadlineExceeded):
		return true
	case errors.Is(err, mailer.ErrCanceled), errors.Is(err, mailer.ErrTimeout),
		errors.Is(err, mailer.ErrInternal), errors.Is(err, mailer.ErrUnavailable):
		return true
	default:
		return false
	}
}

func (a *app) send(ctx context.Context, builder *builder) error {
	name, email, substitutions, ok := builder.build()
	if !ok {
		return errors.New("failed to build")
	}

	if a.debug {
		a.logger.Info("Send email",
			zap.String("氏名", name),
			zap.String("アドレス", email),
			zap.Any("substitutions", substitutions),
		)
		return nil
	}

	return a.client.SendFromInfo(ctx, "default", name, email, substitutions)
}

type builder struct {
	substitutions map[string]string
}

func newBuilder(header, body []string) *builder {
	substitutions := make(map[string]string, len(header))
	for i, key := range header {
		substitutions[key] = body[i]
	}
	builder := &builder{
		substitutions: substitutions,
	}
	builder.set()
	return builder
}

func (b *builder) build() (string, string, map[string]interface{}, bool) {
	if b == nil {
		return "", "", nil, false
	}
	name, ok := b.getName()
	if !ok {
		return "", "", nil, false
	}
	email, ok := b.getEmail()
	if !ok {
		return "", "", nil, false
	}
	substitutions := make(map[string]interface{}, len(b.substitutions))
	for k, v := range b.substitutions {
		substitutions[k] = v
	}
	return name, email, substitutions, true
}

func (b *builder) set() {
	b.setName()
	b.setEmail()
}

func (b *builder) setName() {
	if _, ok := b.substitutions["氏名"]; ok {
		return
	}
	var sei, mei string
	if v, ok := b.substitutions["姓"]; ok {
		sei = v
	}
	if v, ok := b.substitutions["名"]; ok {
		mei = v
	}
	name := fmt.Sprintf("%s %s", sei, mei)
	b.substitutions["氏名"] = strings.TrimSpace(name)
}

func (b *builder) setEmail() {
	if _, ok := b.substitutions["メールアドレス"]; ok {
		return
	}
	if v, ok := b.substitutions["e-mail"]; ok {
		b.substitutions["メールアドレス"] = v
		return
	}
}

func (b *builder) getName() (string, bool) {
	name, ok := b.substitutions["氏名"]
	return name, ok
}

func (b *builder) getEmail() (string, bool) {
	email, ok := b.substitutions["メールアドレス"]
	return email, ok
}
