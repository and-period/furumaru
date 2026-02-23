# テスト生成テンプレート

AI エージェント（Claude Code / Codex）がテストコードを生成する際の参照テンプレート集。
本リポジトリの既存コードから抽出した実践的パターンに基づく。

> **参照元**: 各テンプレートの元となるファイルパスをコメントで記載している。
> パターンが不明な場合は、記載の参照元ファイルを直接確認すること。

---

## データベース層テスト (TiDB)

### TestMain セットアップ

データベーステストでは `testcontainers-go` によるコンテナ起動、または外部 DB 接続を使い分ける。
各パッケージの `*_test.go` に共通の `TestMain` を 1 つだけ定義する。

> 参照: `api/internal/store/database/tidb/tidb_test.go`

```go
package tidb

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/and-period/furumaru/api/pkg/mysql"
)

var (
	dbClient *mysql.Client
	current  = time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
)

func TestMain(m *testing.M) {
	setEnv()
	ctx := context.Background()

	client, cleanup, err := newTestDBClient(ctx)
	if err != nil {
		panic(err)
	}
	defer cleanup()
	dbClient = client

	os.Exit(m.Run())
}

func newTestDBClient(ctx context.Context) (*mysql.Client, func(), error) {
	if !mysql.ShouldUseContainerDB() {
		return newExternalDBClient()
	}
	schemaDir, err := schemaDir()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to resolve schema dir: %w", err)
	}
	client, cleanup, err := mysql.NewContainerDB(ctx,
		mysql.WithContainerDatabase("<database-name>"),
		mysql.WithSchemaDir(schemaDir),
	)
	if err != nil {
		return nil, nil, err
	}
	return client, cleanup, nil
}

// schemaDir: runtime.Caller で自パッケージ位置を基点にスキーマディレクトリを解決
func schemaDir() (string, error) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("failed to get caller information")
	}
	apiRoot, err := mysql.FindProjectRoot(filepath.Dir(filename))
	if err != nil {
		return "", err
	}
	repoRoot := filepath.Dir(apiRoot)
	dir := filepath.Join(repoRoot, "infra", "tidb", "schema", "<database-name>")
	if _, err := os.Stat(dir); err != nil {
		return "", fmt.Errorf("schema dir not found: %w", err)
	}
	return dir, nil
}
```

### deleteAll ヘルパー

テスト前にテーブルをクリーンアップする。**削除順序は外部キー制約に従う**（子テーブルから先に削除）。

```go
func deleteAll(ctx context.Context) error {
	tables := []string{
		// 子テーブルから順に列挙
		"order_items",
		"orders",
		"product_revisions",
		"products",
		"product_types",
		"categories",
	}
	return delete(ctx, tables...)
}

func delete(ctx context.Context, tables ...string) error {
	for _, table := range tables {
		sql := fmt.Sprintf("DELETE FROM %s", table)
		if err := dbClient.DB.WithContext(ctx).Exec(sql).Error; err != nil {
			return err
		}
	}
	return nil
}
```

### テストヘルパー関数

テストデータ生成用のヘルパーを定義する。ファイル末尾にまとめて配置する。

> 参照: `api/internal/store/database/tidb/category_test.go`（末尾）

```go
func testCategory(id, name string, now time.Time) *entity.Category {
	return &entity.Category{
		ID:        id,
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
```

### List テスト

> 参照: `api/internal/store/database/tidb/category_test.go` — `TestCategory_List`

```go
func TestXxx_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}
	err := deleteAll(t.Context())
	require.NoError(t, err)

	// --- テストデータ投入 ---
	categories := make(entity.Categories, 3)
	categories[0] = testCategory("category-id01", "野菜", now())
	categories[1] = testCategory("category-id02", "水産加工物", now())
	categories[2] = testCategory("category-id03", "水産物", now())
	err = db.DB.Create(&categories).Error
	require.NoError(t, err)

	// --- テーブル駆動テスト ---
	type args struct {
		params *database.ListCategoriesParams
	}
	type want struct {
		categories entity.Categories
		hasErr     bool
	}
	tests := []struct {
		name  string
		setup func(ctx context.Context, t *testing.T, db *mysql.Client)
		args  args
		want  want
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				params: &database.ListCategoriesParams{
					Name:   "水産",
					Limit:  1,
					Offset: 1,
				},
			},
			want: want{
				categories: categories[2:],
				hasErr:     false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()
			tt.setup(ctx, t, db)

			db := &category{db: db, now: now}
			actual, err := db.List(ctx, tt.args.params)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
			assert.ElementsMatch(t, tt.want.categories, actual)
		})
	}
}
```

### Get テスト

> 参照: `api/internal/store/database/tidb/category_test.go` — `TestCategory_Get`

```go
func TestXxx_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}
	err := deleteAll(t.Context())
	require.NoError(t, err)

	c := testCategory("category-id", "野菜", now())
	err = db.DB.Create(&c).Error
	require.NoError(t, err)

	type args struct {
		categoryID string
	}
	type want struct {
		category *entity.Category
		hasErr   bool
	}
	tests := []struct {
		name  string
		setup func(ctx context.Context, t *testing.T, db *mysql.Client)
		args  args
		want  want
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				categoryID: "category-id",
			},
			want: want{
				category: c,
				hasErr:   false,
			},
		},
		{
			name:  "not found",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				categoryID: "other-id",
			},
			want: want{
				category: nil,
				hasErr:   true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()
			tt.setup(ctx, t, db)

			db := &category{db: db, now: now}
			actual, err := db.Get(ctx, tt.args.categoryID)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
			assert.Equal(t, tt.want.category, actual)
		})
	}
}
```

### Create テスト

Create テストは `t.Parallel()` を使わず、各テストケース開始時にテーブルを個別にクリーンする点に注意。

> 参照: `api/internal/store/database/tidb/category_test.go` — `TestCategory_Create`

```go
func TestXxx_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}
	err := deleteAll(t.Context())
	require.NoError(t, err)

	type args struct {
		category *entity.Category
	}
	type want struct {
		hasErr bool
	}
	tests := []struct {
		name  string
		setup func(ctx context.Context, t *testing.T, db *mysql.Client)
		args  args
		want  want
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				category: testCategory("category-id", "野菜", now()),
			},
			want: want{
				hasErr: false,
			},
		},
		{
			name: "failed to duplicate entry",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {
				category := testCategory("category-id", "野菜", now())
				err = db.DB.Create(&category).Error
				require.NoError(t, err)
			},
			args: args{
				category: testCategory("category-id", "野菜", now()),
			},
			want: want{
				hasErr: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := t.Context()
			// Create テストでは個別にクリーンアップ（並列実行しない）
			err := delete(ctx, categoryTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &category{db: db, now: now}
			err = db.Create(ctx, tt.args.category)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
		})
	}
}
```

---

## サービス層テスト

### テストインフラ構造

サービス層テストでは `gomock` で依存をモック化し、`testService` ヘルパーで定型処理をラップする。

> 参照: `api/internal/store/service/service_test.go`

```go
// mocks 構造体: サービスが依存する全インターフェースのモックを保持
type mocks struct {
	db         *dbMocks
	cache      *mock_dynamodb.MockClient
	user       *mock_user.MockService
	messenger  *mock_messenger.MockService
	media      *mock_media.MockService
	postalCode *mock_postalcode.MockClient
	// ... 他の依存
}

type dbMocks struct {
	Category *mock_database.MockCategory
	Product  *mock_database.MockProduct
	// ... 他のテーブル
}

// testCaller: テスト本体のシグネチャ
type testCaller func(ctx context.Context, t *testing.T, service *service)

// testService: setup → サービス生成 → テスト実行 をラップ
func testService(
	setup func(ctx context.Context, mocks *mocks),
	testFunc testCaller,
	opts ...testOption,
) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()
		ctx, cancel := context.WithCancel(t.Context())
		defer cancel()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mocks := newMocks(ctrl)
		srv := newService(mocks, opts...)
		setup(ctx, mocks)
		testFunc(ctx, t, srv)
		srv.waitGroup.Wait()
	}
}
```

### List テスト

> 参照: `api/internal/store/service/product_test.go` — `TestListProducts`

```go
func TestListXxx(t *testing.T) {
	t.Parallel()

	now := jst.Date(2022, 6, 28, 18, 30, 0, 0)

	// 期待されるDBパラメータ
	params := &database.ListXxxParams{
		Name:   "検索ワード",
		Limit:  30,
		Offset: 0,
	}

	// テストデータ
	items := entity.Xxxs{
		{
			ID:   "item-id",
			Name: "テストアイテム",
			// ... フィールド
		},
	}

	tests := []struct {
		name        string
		setup       func(ctx context.Context, mocks *mocks)
		input       *store.ListXxxInput
		expect      entity.Xxxs
		expectTotal int64
		expectErr   error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Xxx.EXPECT().List(gomock.Any(), params).Return(items, nil)
				mocks.db.Xxx.EXPECT().Count(gomock.Any(), params).Return(int64(1), nil)
			},
			input: &store.ListXxxInput{
				Name:   "検索ワード",
				Limit:  30,
				Offset: 0,
			},
			expect:      items,
			expectTotal: 1,
			expectErr:   nil,
		},
		{
			name:        "invalid argument",
			setup:       func(ctx context.Context, mocks *mocks) {},
			input:       &store.ListXxxInput{},
			expect:      nil,
			expectTotal: 0,
			expectErr:   exception.ErrInvalidArgument,
		},
		{
			name: "failed to list",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Xxx.EXPECT().List(gomock.Any(), params).Return(nil, assert.AnError)
				mocks.db.Xxx.EXPECT().Count(gomock.Any(), params).Return(int64(1), nil)
			},
			input: &store.ListXxxInput{
				Name:   "検索ワード",
				Limit:  30,
				Offset: 0,
			},
			expect:      nil,
			expectTotal: 0,
			expectErr:   exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, total, err := service.ListXxx(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.ElementsMatch(t, tt.expect, actual)
			assert.Equal(t, tt.expectTotal, total)
		}, withNow(now)))
	}
}
```

### Get / MultiGet テスト

```go
func TestMultiGetXxx(t *testing.T) {
	t.Parallel()

	now := jst.Date(2022, 5, 2, 18, 30, 0, 0)
	items := entity.Xxxs{
		{
			ID:        "item-id",
			Name:      "テストアイテム",
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.MultiGetXxxInput
		expect    entity.Xxxs
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Xxx.EXPECT().MultiGet(ctx, []string{"item-id"}).Return(items, nil)
			},
			input: &store.MultiGetXxxInput{
				XxxIDs: []string{"item-id"},
			},
			expect:    items,
			expectErr: nil,
		},
		{
			name: "failed to multi get",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Xxx.EXPECT().MultiGet(ctx, []string{"item-id"}).Return(nil, assert.AnError)
			},
			input: &store.MultiGetXxxInput{
				XxxIDs: []string{"item-id"},
			},
			expect:    nil,
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.MultiGetXxx(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.ElementsMatch(t, tt.expect, actual)
		}, withNow(now)))
	}
}
```

---

## エンティティテスト

### コンストラクタテスト

エンティティの `NewXxx()` コンストラクタをテストする。生成される ID はランダムなので `""` で無視する。

> 参照: `api/internal/store/entity/category_test.go` — `TestCategory`

```go
func TestXxx(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		params *NewXxxParams
		expect *Xxx
		hasErr bool
	}{
		{
			name: "success",
			params: &NewXxxParams{
				Name: "テスト",
				// ... 他のフィールド
			},
			expect: &Xxx{
				Name: "テスト",
				// ... 他のフィールド
			},
			hasErr: false,
		},
		{
			name: "invalid argument",
			params: &NewXxxParams{
				// 不正な入力
			},
			expect: nil,
			hasErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual, err := NewXxx(tt.params)
			assert.Equal(t, tt.hasErr, err != nil, err)
			if actual != nil {
				actual.ID = "" // ランダム生成IDを無視
			}
			assert.Equal(t, tt.expect, actual)
		})
	}
}
```

> **注意**: コンストラクタがエラーを返さない場合は `hasErr` フィールドを省略し、
> `assert.Equal(t, tt.expect, actual)` のみで比較する。
> 参照元に合わせてシグネチャを確認すること。

### メソッドテスト

エンティティのメソッド（Fill、Validate、集計メソッドなど）のテスト。

```go
func TestXxx_MethodName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		xxx    *Xxx
		expect ResultType
	}{
		{
			name: "success",
			xxx: &Xxx{
				ID:   "xxx-id",
				Name: "テスト",
			},
			expect: expectedValue,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := tt.xxx.MethodName()
			assert.Equal(t, tt.expect, actual)
		})
	}
}
```

---

## Gateway サービスラッパーテスト

### NewXxx コンストラクタテスト

Entity → Response DTO 変換をテストする。

> 参照: `api/internal/gateway/admin/v1/service/administrator_test.go` — `TestAdministrator`

```go
func TestXxx(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		entity *entity.Xxx
		expect *Xxx
	}{
		{
			name: "success",
			entity: &entity.Xxx{
				ID:        "xxx-id",
				Name:      "テスト",
				CreatedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
				UpdatedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
			},
			expect: &Xxx{
				Xxx: types.Xxx{
					ID:        "xxx-id",
					Name:      "テスト",
					CreatedAt: 1640962800, // Unix タイムスタンプ
					UpdatedAt: 1640962800,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, NewXxx(tt.entity))
		})
	}
}
```

### Response メソッドテスト

```go
func TestXxx_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		xxx    *Xxx
		expect *types.Xxx
	}{
		{
			name: "success",
			xxx: &Xxx{
				Xxx: types.Xxx{
					ID:        "xxx-id",
					Name:      "テスト",
					CreatedAt: 1640962800,
					UpdatedAt: 1640962800,
				},
			},
			expect: &types.Xxx{
				ID:        "xxx-id",
				Name:      "テスト",
				CreatedAt: 1640962800,
				UpdatedAt: 1640962800,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.xxx.Response())
		})
	}
}
```

### バッチ変換テスト (NewXxxs / Response)

```go
func TestXxxs(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		entities entity.Xxxs
		expect   Xxxs
	}{
		{
			name: "success",
			entities: entity.Xxxs{
				{ID: "xxx-id01", Name: "テスト01"},
				{ID: "xxx-id02", Name: "テスト02"},
			},
			expect: Xxxs{
				{Xxx: types.Xxx{ID: "xxx-id01", Name: "テスト01"}},
				{Xxx: types.Xxx{ID: "xxx-id02", Name: "テスト02"}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, NewXxxs(tt.entities))
		})
	}
}
```

---

## Gateway ハンドラーテスト

### テストインフラ構造

ハンドラーテストでは `gin.TestMode` + `httptest.NewRecorder` を使い、HTTP リクエスト/レスポンスを検証する。

> 参照: `api/internal/gateway/user/v1/handler/handler_test.go`

```go
// テスト用ヘルパー関数群
func testGet(
	t *testing.T,
	setup func(*testing.T, *mocks, *gomock.Controller),
	expect *testResponse,
	path string,
	opts ...testOption,
) {
	testHTTP(t, setup, expect, newHTTPRequest(t, http.MethodGet, path, nil), opts...)
}

func testPost(
	t *testing.T,
	setup func(*testing.T, *mocks, *gomock.Controller),
	expect *testResponse,
	path string,
	body interface{},
	opts ...testOption,
) {
	testHTTP(t, setup, expect, newHTTPRequest(t, http.MethodPost, path, body), opts...)
}

// testHTTP: 共通のHTTPテスト実行関数
func testHTTP(
	t *testing.T,
	setup func(*testing.T, *mocks, *gomock.Controller),
	expect *testResponse,
	req *http.Request,
	opts ...testOption,
) {
	t.Parallel()
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mocks := newMocks(ctrl)

	dopts := &testOptions{now: jst.Now}
	for i := range opts {
		opts[i](dopts)
	}

	h := newHandler(mocks, dopts)
	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)
	newRoutes(h, r)
	setup(t, mocks, ctrl)

	r.ServeHTTP(w, req)
	require.Equal(t, expect.code, w.Code)
	if isError(w) || expect.body == nil {
		return
	}
	body, err := json.Marshal(expect.body)
	require.NoError(t, err)
	require.JSONEq(t, string(body), w.Body.String())
}
```

### エンドポイントテスト（テーブル駆動）

```go
func TestListXxx(t *testing.T) {
	t.Parallel()

	items := entity.Xxxs{
		{ID: "xxx-id", Name: "テスト"},
	}

	tests := []struct {
		name   string
		setup  func(*testing.T, *mocks, *gomock.Controller)
		expect *testResponse
		path   string
	}{
		{
			name: "success",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				in := &store.ListXxxInput{Limit: 20, Offset: 0}
				mocks.store.EXPECT().ListXxx(gomock.Any(), in).Return(items, int64(1), nil)
			},
			expect: &testResponse{
				code: http.StatusOK,
				body: &types.XxxsResponse{
					Xxxs:  service.NewXxxs(items).Response(),
					Total: 1,
				},
			},
			path: "/v1/xxxs",
		},
		{
			name: "failed to list",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				in := &store.ListXxxInput{Limit: 20, Offset: 0}
				mocks.store.EXPECT().ListXxx(gomock.Any(), in).Return(nil, int64(0), assert.AnError)
			},
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
			path: "/v1/xxxs",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testGet(t, tt.setup, tt.expect, tt.path)
		})
	}
}
```

---

## モック生成

### go:generate ディレクティブ

モックは `go tool mockgen` で生成する。インターフェース定義ファイルの先頭に以下を記述。

> 参照: `api/internal/store/database/database.go`, `api/internal/store/service.go`

```go
//go:generate go tool mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=./../../../mock/store/$GOPACKAGE/$GOFILE
package database
```

生成コマンド:

```bash
cd api && go generate ./internal/store/...
```

### モックの使用パターン

```go
import (
	mock_database "github.com/and-period/furumaru/api/mock/store/database"
	"go.uber.org/mock/gomock"
)

ctrl := gomock.NewController(t)
defer ctrl.Finish()

mockProduct := mock_database.NewMockProduct(ctrl)
mockProduct.EXPECT().Get(gomock.Any(), "product-id").Return(product, nil)
mockProduct.EXPECT().List(gomock.Any(), gomock.Any()).Return(products, nil)
```

---

## テスト実行コマンド

```bash
# 全テスト実行
cd api && make test

# 特定パッケージ
go test ./internal/store/database/tidb/... -v

# 特定テスト関数
go test ./internal/store/service/... -run TestListProducts -v

# カバレッジ付き
go test ./internal/store/... -cover
```

---

## 注意事項

1. **テストライブラリ**: `stretchr/testify`（`assert` / `require`）と `go.uber.org/mock/gomock` を使用
2. **時刻固定**: `jst.Date()` または `time.Date()` で固定時刻を使う。`time.Now()` は使わない
3. **並列実行**: `t.Parallel()` を基本的に使用する。ただし Create/Update/Delete テストでは DB 状態が競合するため使用しない場合がある
4. **エラー検証**: `assert.Equal(t, tt.want.hasErr, err != nil, err)` または `assert.ErrorIs(t, err, tt.expectErr)` を使い分ける
5. **ID 無視**: コンストラクタテストでは `actual.ID = ""` でランダム生成 ID を無視する
6. **goleak**: サービス層の `TestMain` では `go.uber.org/goleak` で goroutine リークを検出する
