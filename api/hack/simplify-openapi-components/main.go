package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"regexp"
	"time"

	"github.com/and-period/furumaru/api/pkg/jst"
	"gopkg.in/yaml.v3"
)

var (
	sourceFile    string
	addCookieAuth bool
	debug         bool
)

type app struct {
	source        string
	pattern       *regexp.Regexp
	addCookieAuth bool
	debug         bool
}

func main() {
	startedAt := jst.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
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

func setup(_ context.Context) (*app, error) {
	flag.StringVar(&sourceFile, "source-file", "", "対象のOpenAPI仕様書ファイル (swagger.yaml)")
	flag.BoolVar(&addCookieAuth, "add-cookie-auth", false, "cookie認証スキームを追加")
	flag.BoolVar(&debug, "debug", false, "デバッグモード")
	flag.Parse()

	// コマンドライン引数から取得
	if sourceFile == "" && flag.NArg() > 0 {
		sourceFile = flag.Arg(0)
	}

	if sourceFile == "" {
		return nil, fmt.Errorf("source-file is required")
	}

	// 正規表現パターンを定義
	// github_com_and-period_furumaru_api_internal_gateway_[^.]+\.(response|request)\.(\w+) または
	// github_com_and-period_furumaru_api_internal_gateway_[^.]+\.(\w+) にマッチ
	pattern := regexp.MustCompile(`github_com_and-period_furumaru_api_internal_gateway_[^.]+\.((?:response|request)\.)?(\w+)`)

	app := &app{
		source:        sourceFile,
		pattern:       pattern,
		addCookieAuth: addCookieAuth,
		debug:         debug,
	}
	return app, nil
}

func (a *app) run(ctx context.Context) error {
	slog.Debug("Start to run", slog.Bool("debug", a.debug))

	// ファイルを読み込む
	content, err := os.ReadFile(a.source)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	slog.Info("Processing OpenAPI specification",
		slog.String("file", a.source),
		slog.Int("size", len(content)),
	)

	// コンポーネント名を簡潔化
	count := 0
	result := a.pattern.ReplaceAllFunc(content, func(match []byte) []byte {
		count++

		// マッチした部分を解析
		submatches := a.pattern.FindSubmatch(match)
		if len(submatches) >= 3 {
			prefix := submatches[1]   // response. or request. or empty
			typeName := submatches[2] // 型名

			// デバッグモードの場合、変換内容を出力
			if a.debug {
				slog.Debug("Replacing component name",
					slog.String("from", string(match)),
					slog.String("to", string(append(prefix, typeName...))),
				)
			}

			// prefixが空でない場合はそのまま使用、空の場合は型名のみ
			if len(prefix) > 0 {
				return append(prefix, typeName...)
			}
			return typeName
		}
		return match
	})

	// cookie認証スキームを追加
	if a.addCookieAuth {
		result, err = a.addCookieAuthScheme(result)
		if err != nil {
			return fmt.Errorf("failed to add cookie auth scheme: %w", err)
		}
	}

	// ファイルに書き戻す
	err = os.WriteFile(a.source, result, 0644)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	slog.Info("Successfully processed OpenAPI specification",
		slog.String("file", a.source),
		slog.Int("component_count", count),
		slog.Bool("cookie_auth_added", a.addCookieAuth),
	)

	return nil
}

// addCookieAuthScheme adds cookie authentication scheme to the OpenAPI specification
func (a *app) addCookieAuthScheme(content []byte) ([]byte, error) {
	// YAMLをパース
	var doc yaml.Node
	if err := yaml.Unmarshal(content, &doc); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	// componentsノードを探す
	var componentsNode *yaml.Node
	if doc.Kind == yaml.DocumentNode && len(doc.Content) > 0 {
		rootNode := doc.Content[0]
		if rootNode.Kind == yaml.MappingNode {
			for i := 0; i < len(rootNode.Content); i += 2 {
				if rootNode.Content[i].Value == "components" {
					componentsNode = rootNode.Content[i+1]
					break
				}
			}
		}
	}

	if componentsNode == nil {
		slog.Warn("components section not found, skipping cookie auth addition")
		return content, nil
	}

	// securitySchemesノードを探す
	var securitySchemesNode *yaml.Node
	if componentsNode.Kind == yaml.MappingNode {
		for i := 0; i < len(componentsNode.Content); i += 2 {
			if componentsNode.Content[i].Value == "securitySchemes" {
				securitySchemesNode = componentsNode.Content[i+1]
				break
			}
		}
	}

	if securitySchemesNode == nil {
		slog.Warn("securitySchemes section not found, skipping cookie auth addition")
		return content, nil
	}

	// cookieauthが既に存在するか確認
	if securitySchemesNode.Kind == yaml.MappingNode {
		for i := 0; i < len(securitySchemesNode.Content); i += 2 {
			if securitySchemesNode.Content[i].Value == "cookieauth" {
				slog.Info("cookieauth already exists, skipping addition")
				return content, nil
			}
		}

		// cookieauthノードを作成
		cookieAuthKey := &yaml.Node{
			Kind:  yaml.ScalarNode,
			Value: "cookieauth",
		}

		cookieAuthValue := &yaml.Node{
			Kind: yaml.MappingNode,
			Content: []*yaml.Node{
				{Kind: yaml.ScalarNode, Value: "in"},
				{Kind: yaml.ScalarNode, Value: "cookie"},
				{Kind: yaml.ScalarNode, Value: "name"},
				{Kind: yaml.ScalarNode, Value: "session_id"},
				{Kind: yaml.ScalarNode, Value: "type"},
				{Kind: yaml.ScalarNode, Value: "apiKey"},
			},
		}

		// securitySchemesに追加
		securitySchemesNode.Content = append(securitySchemesNode.Content, cookieAuthKey, cookieAuthValue)

		slog.Info("Successfully added cookieauth to securitySchemes")
	}

	// YAMLに戻す
	result, err := yaml.Marshal(&doc)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal YAML: %w", err)
	}

	return result, nil
}