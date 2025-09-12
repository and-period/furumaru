package main

import (
	"bytes"
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
	convertEnums  bool
	addRequired   bool
)

type app struct {
	source        string
	pattern       *regexp.Regexp
	addCookieAuth bool
	debug         bool
	convertEnums  bool
	addRequired   bool
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
	flag.BoolVar(&convertEnums, "convert-enums", true, "enum型の変換を実行")
	flag.BoolVar(&addRequired, "add-required", true, "全フィールドにrequiredを自動付与")
	flag.Parse()

	// コマンドライン引数から取得
	if sourceFile == "" && flag.NArg() > 0 {
		sourceFile = flag.Arg(0)
	}

	if sourceFile == "" {
		return nil, fmt.Errorf("source-file is required")
	}

	logLevel := slog.LevelInfo
	if debug {
		logLevel = slog.LevelDebug
	}

	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	})))

	// 正規表現パターンを定義
	pattern := regexp.MustCompile(`(util|types|github_com_and-period_furumaru_api_internal_.+)\.(\w+)`)

	app := &app{
		source:        sourceFile,
		pattern:       pattern,
		addCookieAuth: addCookieAuth,
		debug:         debug,
		convertEnums:  convertEnums,
		addRequired:   addRequired,
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
		if len(submatches) >= 2 {
			typeName := submatches[len(submatches)-1] // 型名

			// デバッグモードの場合、変換内容を出力
			slog.Debug("Replacing component name", slog.String("from", string(match)), slog.String("to", string(typeName)))

			// types プレフィックスは削除して型名のみを使用
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

	// schemas処理を実行（enum変換とrequired付与を統合）
	if a.convertEnums || a.addRequired {
		result, err = a.processSchemas(result)
		if err != nil {
			return fmt.Errorf("failed to process schemas: %w", err)
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
		slog.Bool("enums_converted", a.convertEnums),
		slog.Bool("required_fields_added", a.addRequired),
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

	// YAMLに戻す（インデント2スペース）
	var buf bytes.Buffer
	encoder := yaml.NewEncoder(&buf)
	encoder.SetIndent(2)
	if err := encoder.Encode(&doc); err != nil {
		return nil, fmt.Errorf("failed to marshal YAML: %w", err)
	}
	encoder.Close()

	return buf.Bytes(), nil
}

// convertEnumSchema converts a single enum schema to standard format
func (a *app) convertEnumSchema(schemaName string, schemaNode *yaml.Node) bool {
	if schemaNode.Kind != yaml.MappingNode {
		return false
	}

	// x-enum-varnamesとx-enum-commentsの存在を確認
	var xEnumVarnamesNode, xEnumCommentsNode, typeNode *yaml.Node
	hasEnum := false

	for i := 0; i < len(schemaNode.Content); i += 2 {
		key := schemaNode.Content[i].Value
		value := schemaNode.Content[i+1]

		switch key {
		case "x-enum-varnames":
			xEnumVarnamesNode = value
		case "x-enum-comments":
			xEnumCommentsNode = value
		case "type":
			typeNode = value
		case "enum":
			hasEnum = true
		}
	}

	// 既にenum配列がある場合はスキップ
	if hasEnum {
		return false
	}

	// x-enum-varnamesがなければenum型ではない
	if xEnumVarnamesNode == nil || xEnumVarnamesNode.Kind != yaml.SequenceNode {
		return false
	}

	// typeがintegerでない場合はスキップ
	if typeNode == nil || typeNode.Value != "integer" {
		return false
	}

	slog.Debug("Converting enum schema", slog.String("schema", schemaName))

	// x-enum-varnamesから重複を除去
	uniqueVarnames := make([]*yaml.Node, 0)
	seen := make(map[string]bool)
	for _, varnameNode := range xEnumVarnamesNode.Content {
		varname := varnameNode.Value
		if !seen[varname] {
			seen[varname] = true
			uniqueVarnames = append(uniqueVarnames, varnameNode)
		}
	}

	// 重複を除去したx-enum-varnamesで上書き
	xEnumVarnamesNode.Content = uniqueVarnames

	// enum配列を生成（0から順番に）
	enumCount := len(uniqueVarnames)
	enumValues := make([]*yaml.Node, enumCount)
	for i := 0; i < enumCount; i++ {
		enumValues[i] = &yaml.Node{
			Kind:  yaml.ScalarNode,
			Value: fmt.Sprintf("%d", i),
		}
	}

	// x-enum-descriptionsを生成
	enumDescriptions := make([]*yaml.Node, enumCount)
	for i, varnameNode := range uniqueVarnames {
		varname := varnameNode.Value
		description := "不明" // デフォルト値

		// x-enum-commentsから対応する説明を取得
		if xEnumCommentsNode != nil && xEnumCommentsNode.Kind == yaml.MappingNode {
			for j := 0; j < len(xEnumCommentsNode.Content); j += 2 {
				if xEnumCommentsNode.Content[j].Value == varname {
					description = xEnumCommentsNode.Content[j+1].Value
					break
				}
			}
		}

		enumDescriptions[i] = &yaml.Node{
			Kind:  yaml.ScalarNode,
			Value: description,
		}
	}

	// 新しいフィールドを追加
	addEnumField := func(key string, value *yaml.Node) {
		keyNode := &yaml.Node{
			Kind:  yaml.ScalarNode,
			Value: key,
		}
		schemaNode.Content = append(schemaNode.Content, keyNode, value)
	}

	// enum配列を追加
	enumArrayNode := &yaml.Node{
		Kind:    yaml.SequenceNode,
		Content: enumValues,
	}
	addEnumField("enum", enumArrayNode)

	// format: int32を追加
	formatNode := &yaml.Node{
		Kind:  yaml.ScalarNode,
		Value: "int32",
	}
	addEnumField("format", formatNode)

	// nullable: falseを追加
	nullableNode := &yaml.Node{
		Kind:  yaml.ScalarNode,
		Tag:   "!!bool",
		Value: "false",
	}
	addEnumField("nullable", nullableNode)

	// x-enum-descriptionsを追加
	enumDescriptionsArrayNode := &yaml.Node{
		Kind:    yaml.SequenceNode,
		Content: enumDescriptions,
	}
	addEnumField("x-enum-descriptions", enumDescriptionsArrayNode)

	return true
}

// addRequiredToSchema adds required fields to a single schema for all properties
func (a *app) addRequiredToSchema(schemaName string, schemaNode *yaml.Node) int {
	if schemaNode.Kind != yaml.MappingNode {
		return 0
	}

	var propertiesNode, requiredNode *yaml.Node
	var typeValue string

	// スキーマの構造を解析
	for i := 0; i < len(schemaNode.Content); i += 2 {
		key := schemaNode.Content[i].Value
		value := schemaNode.Content[i+1]

		switch key {
		case "properties":
			propertiesNode = value
		case "required":
			requiredNode = value
		case "type":
			typeValue = value.Value
		}
	}

	// objectタイプのスキーマのみ処理
	if typeValue != "object" || propertiesNode == nil || propertiesNode.Kind != yaml.MappingNode {
		return 0
	}

	// 既存のrequiredフィールドを取得
	existingRequired := make(map[string]bool)
	if requiredNode != nil && requiredNode.Kind == yaml.SequenceNode {
		for _, item := range requiredNode.Content {
			existingRequired[item.Value] = true
		}
	}

	// すべてのフィールドを収集（既存のrequiredに含まれていないもの）
	var requiredFields []string
	for j := 0; j < len(propertiesNode.Content); j += 2 {
		fieldName := propertiesNode.Content[j].Value

		if !existingRequired[fieldName] {
			requiredFields = append(requiredFields, fieldName)
		}
	}

	if len(requiredFields) == 0 {
		return 0
	}

	slog.Debug("Adding required fields to schema",
		slog.String("schema", schemaName),
		slog.Any("fields", requiredFields))

	// requiredノードが存在しない場合は作成
	if requiredNode == nil {
		// requiredキーと値ノードを作成
		requiredKeyNode := &yaml.Node{
			Kind:  yaml.ScalarNode,
			Value: "required",
		}
		requiredNode = &yaml.Node{
			Kind:    yaml.SequenceNode,
			Content: make([]*yaml.Node, 0),
		}

		// スキーマに追加
		schemaNode.Content = append(schemaNode.Content, requiredKeyNode, requiredNode)
	}

	// 新しいrequiredフィールドを追加
	for _, fieldName := range requiredFields {
		fieldNode := &yaml.Node{
			Kind:  yaml.ScalarNode,
			Value: fieldName,
		}
		requiredNode.Content = append(requiredNode.Content, fieldNode)
	}

	return len(requiredFields)
}

// processSchemas processes schemas for both enum conversion and required field addition in a single pass
func (a *app) processSchemas(content []byte) ([]byte, error) {
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
		slog.Warn("components section not found, skipping schemas processing")
		return content, nil
	}

	// schemasノードを探す
	var schemasNode *yaml.Node
	if componentsNode.Kind == yaml.MappingNode {
		for i := 0; i < len(componentsNode.Content); i += 2 {
			if componentsNode.Content[i].Value == "schemas" {
				schemasNode = componentsNode.Content[i+1]
				break
			}
		}
	}

	if schemasNode == nil {
		slog.Warn("schemas section not found, skipping schemas processing")
		return content, nil
	}

	// カウンター
	enumConvertedCount := 0
	requiredAddedCount := 0

	// 各スキーマをチェック
	if schemasNode.Kind == yaml.MappingNode {
		for i := 0; i < len(schemasNode.Content); i += 2 {
			schemaName := schemasNode.Content[i].Value
			schemaNode := schemasNode.Content[i+1]

			// enum処理
			if a.convertEnums && a.convertEnumSchema(schemaName, schemaNode) {
				enumConvertedCount++
			}

			// required処理
			if a.addRequired {
				count := a.addRequiredToSchema(schemaName, schemaNode)
				requiredAddedCount += count
			}
		}
	}

	if a.convertEnums {
		slog.Info("Enum conversion completed", slog.Int("converted_count", enumConvertedCount))
	}
	if a.addRequired {
		slog.Info("Required fields addition completed", slog.Int("fields_added", requiredAddedCount))
	}

	// YAMLに戻す（インデント2スペース）
	var buf bytes.Buffer
	encoder := yaml.NewEncoder(&buf)
	encoder.SetIndent(2)
	if err := encoder.Encode(&doc); err != nil {
		return nil, fmt.Errorf("failed to marshal YAML: %w", err)
	}
	encoder.Close()

	return buf.Bytes(), nil
}
