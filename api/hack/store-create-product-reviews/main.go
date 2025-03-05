package main

import (
	"context"
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"io"
	"math/rand/v2"
	"os"
	"strconv"
	"time"

	sdb "github.com/and-period/furumaru/api/internal/store/database"
	stidb "github.com/and-period/furumaru/api/internal/store/database/tidb"
	sentity "github.com/and-period/furumaru/api/internal/store/entity"
	udb "github.com/and-period/furumaru/api/internal/user/database"
	utidb "github.com/and-period/furumaru/api/internal/user/database/tidb"
	uentity "github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/mysql"
	"github.com/and-period/furumaru/api/pkg/secret"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"go.uber.org/zap"
)

var (
	tidbHost     string
	tidbPort     string
	tidbUsername string
	tidbPassword string
	awsRegion    string
	dbSecretName string
	filepath     string
	debug        bool

	userParams = []*uentity.NewUserParams{
		{
			Registered:    false,
			Lastname:      "佐藤",
			Firstname:     "優子",
			LastnameKana:  "サトウ",
			FirstnameKana: "ユウコ",
			AccountID:     "yuko_sato",
			Username:      "ゆうたろう",
			Email:         "sato.yuko@example.com",
		},
		{
			Registered:    false,
			Lastname:      "鈴木",
			Firstname:     "大翔",
			LastnameKana:  "スズキ",
			FirstnameKana: "ヒロト",
			AccountID:     "hiroto_suzuki",
			Username:      "Dragon_Tamer",
			Email:         "suzuki.hiroto@example.com",
		},
		{
			Registered:    false,
			Lastname:      "田中",
			Firstname:     "美咲",
			LastnameKana:  "タナカ",
			FirstnameKana: "ミサキ",
			AccountID:     "misaki_tanaka",
			Username:      "みさにゃん",
			Email:         "tanaka.misaki@example.com",
		},
		{
			Registered:    false,
			Lastname:      "渡辺",
			Firstname:     "翔太",
			LastnameKana:  "ワタナベ",
			FirstnameKana: "ショウタ",
			AccountID:     "shota_watanabe",
			Username:      "空色ショータ",
			Email:         "watanabe.shota@example.com",
		},
		{
			Registered:    false,
			Lastname:      "伊藤",
			Firstname:     "愛",
			LastnameKana:  "イトウ",
			FirstnameKana: "アイ",
			AccountID:     "ai_ito",
			Username:      "あいぽん",
			Email:         "ito.ai@example.com",
		},
		{
			Registered:    false,
			Lastname:      "中村",
			Firstname:     "陽一",
			LastnameKana:  "ナカムラ",
			FirstnameKana: "ヨウイチ",
			AccountID:     "yoichi_nakamura",
			Username:      "よっちん",
			Email:         "nakamura.yoichi@example.com",
		},
		{
			Registered:    false,
			Lastname:      "小林",
			Firstname:     "凛",
			LastnameKana:  "コバヤシ",
			FirstnameKana: "リン",
			AccountID:     "rin_kobayashi",
			Username:      "りんりん星人",
			Email:         "kobayashi.rin@example.com",
		},
		{
			Registered:    false,
			Lastname:      "加藤",
			Firstname:     "大和",
			LastnameKana:  "カトウ",
			FirstnameKana: "ヤマト",
			AccountID:     "yamato_kato",
			Username:      "やまとん",
			Email:         "kato.yamato@example.com",
		},
		{
			Registered:    false,
			Lastname:      "木村",
			Firstname:     "結衣",
			LastnameKana:  "キムラ",
			FirstnameKana: "ユイ",
			AccountID:     "yui_kimura",
			Username:      "ゆいぴょん",
			Email:         "kimura.yui@example.com",
		},
		{
			Registered:    false,
			Lastname:      "山本",
			Firstname:     "蓮",
			LastnameKana:  "ヤマモト",
			FirstnameKana: "レン",
			AccountID:     "ren_yamamoto",
			Username:      "れんれん侍",
			Email:         "yamamoto.ren@example.com",
		},
	}
)

type app struct {
	logger *zap.Logger
	store  *sdb.Database
	user   *udb.Database
}

func main() {
	startedAt := jst.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Minute)
	defer cancel()

	app, err := setup(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to setup: %v\n", err)
		os.Exit(1)
	}

	if err := app.run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "failed to run: %v\n", err)
		os.Exit(1)
	}

	endAt := jst.Now()

	const format = "2006-01-02 15:04:05"
	fmt.Printf("Done: %s -> %s (%s)\n",
		jst.Format(startedAt, format), jst.Format(endAt, format),
		endAt.Sub(startedAt).Truncate(time.Second).String(),
	)
}

func setup(ctx context.Context) (*app, error) {
	flag.StringVar(&tidbHost, "tidb-host", "localhost", "TiDB host")
	flag.StringVar(&tidbPort, "tidb-port", "4000", "TiDB port")
	flag.StringVar(&tidbUsername, "tidb-username", "root", "TiDB username")
	flag.StringVar(&tidbPassword, "tidb-password", "", "TiDB password")
	flag.StringVar(&awsRegion, "aws-region", "ap-northeast-1", "AWS region")
	flag.StringVar(&dbSecretName, "db-secret-name", "", "AWS Secret Manager secret name for TiDB")
	flag.StringVar(&filepath, "filepath", "", "filepath")
	flag.BoolVar(&debug, "debug", true, "debug mode")
	flag.Parse()

	// Loggerの設定
	logger, err := zap.NewDevelopment()
	if err != nil {
		return nil, fmt.Errorf("failed to create logger: %w", err)
	}

	if dbSecretName != "" {
		// AWS SDKの設定
		awscfg, err := awsconfig.LoadDefaultConfig(ctx, awsconfig.WithRegion(awsRegion))
		if err != nil {
			return nil, fmt.Errorf("failed to load aws config: %w", err)
		}

		// AWS Secrets Managerの設定
		secret := secret.NewClient(awscfg)

		secrets, err := secret.Get(ctx, dbSecretName)
		if err != nil {
			return nil, fmt.Errorf("failed to get secret: %w", err)
		}
		tidbHost = secrets["host"]
		tidbPort = secrets["port"]
		tidbUsername = secrets["username"]
		tidbPassword = secrets["password"]
	}

	// TiDBの設定
	storedb, err := mysql.NewTiDBClient(&mysql.Params{
		Host:     tidbHost,
		Port:     tidbPort,
		Database: "stores",
		Username: tidbUsername,
		Password: tidbPassword,
	}, mysql.WithNow(jst.Now), mysql.WithLocation(jst.Location()))
	if err != nil {
		return nil, fmt.Errorf("failed to create tidb client for store: %w", err)
	}

	userdb, err := mysql.NewTiDBClient(&mysql.Params{
		Host:     tidbHost,
		Port:     tidbPort,
		Database: "users",
		Username: tidbUsername,
		Password: tidbPassword,
	}, mysql.WithNow(jst.Now), mysql.WithLocation(jst.Location()))
	if err != nil {
		return nil, fmt.Errorf("failed to create tidb client for user: %w", err)
	}

	app := &app{
		logger: logger,
		store:  stidb.NewDatabase(storedb),
		user:   utidb.NewDatabase(userdb),
	}
	return app, nil
}

func (a *app) run(ctx context.Context) error {
	guests, err := a.createGuests(ctx)
	if err != nil {
		return fmt.Errorf("failed to create guests: %w", err)
	}
	a.logger.Info("created guests", zap.Int("count", len(guests)))

	reviews, err := a.readCSVFile(ctx, guests)
	if err != nil {
		return fmt.Errorf("failed to read csv file: %w", err)
	}
	a.logger.Info("read csv file", zap.Int("count", len(reviews)))

	if err := a.createProductReviews(ctx, reviews); err != nil {
		return fmt.Errorf("failed to create product reviews: %w", err)
	}
	a.logger.Info("created product reviews", zap.Int("count", len(reviews)))

	return nil
}

func (a *app) createGuests(ctx context.Context) (uentity.Guests, error) {
	res := make(uentity.Guests, 0, len(userParams))
	for _, params := range userParams {
		guest, err := a.user.Guest.GetByEmail(ctx, params.Email)
		if err != nil && !errors.Is(err, udb.ErrNotFound) {
			return nil, fmt.Errorf("failed to get guest: %w", err)
		}
		if guest == nil { // 登録処理
			user := uentity.NewUser(params)
			if err := a.user.Guest.Create(ctx, user); err != nil {
				return nil, fmt.Errorf("failed to create guest: %w", err)
			}
			guest = &user.Guest
		} else { // 更新処理
			params := &udb.UpdateGuestParams{
				Lastname:      params.Lastname,
				Firstname:     params.Firstname,
				LastnameKana:  params.LastnameKana,
				FirstnameKana: params.FirstnameKana,
			}
			if err := a.user.Guest.Update(ctx, guest.UserID, params); err != nil {
				return nil, fmt.Errorf("failed to update guest: %w", err)
			}
		}
		res = append(res, guest)
	}
	return res, nil
}

func (a *app) readCSVFile(_ context.Context, guests uentity.Guests) (sentity.ProductReviews, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open csv file: %w", err)
	}
	defer f.Close()

	reader := csv.NewReader(f)

	// ヘッダーを読み飛ばす
	if _, err := reader.Read(); err != nil {
		return nil, fmt.Errorf("failed to read csv file: %w", err)
	}

	res := make(sentity.ProductReviews, 0)
	for {
		// "商品ID","評価","タイトル","コメント"
		record, err := reader.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read csv file: %w", err)
		}

		guest := guests[rand.N(len(guests))]

		rate, err := strconv.ParseInt(record[1], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse rate: %w", err)
		}

		review := &sentity.NewProductReviewParams{
			ProductID: record[0],
			UserID:    guest.UserID,
			Rate:      rate,
			Title:     record[2],
			Comment:   record[3],
		}
		res = append(res, sentity.NewProductReview(review))
	}
	return res, nil
}

func (a *app) createProductReviews(ctx context.Context, reviews sentity.ProductReviews) error {
	for _, review := range reviews {
		current, _, err := a.store.ProductReview.List(ctx, &sdb.ListProductReviewsParams{
			ProductID: review.ProductID,
			UserID:    review.UserID,
		})
		if err != nil {
			return fmt.Errorf("failed to list product reviews: %w", err)
		}
		if len(current) > 0 {
			a.logger.Info("product review already exists", zap.Any("review", review))
			continue
		}
		if debug {
			a.logger.Debug("create product review", zap.Any("review", review))
			continue
		}
		if err := a.store.ProductReview.Create(ctx, review); err != nil {
			return fmt.Errorf("failed to create product review: %w", err)
		}
	}
	return nil
}
