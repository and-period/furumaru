package handler

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	sdk "github.com/anthropics/anthropic-sdk-go"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/types"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/pkg/anthropic"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var aiChatTools = []anthropic.Tool{
	{
		Name: "updateProductForm",
		Description: `商品登録フォームのフィールドを更新する。
ユーザーが自然言語で伝えた商品情報をフォームフィールドに変換する。
変更するフィールドのみを含め、変更不要なフィールドは省略すること。`,
		InputSchema: map[string]interface{}{
			"name":                 map[string]any{"type": "string", "maxLength": 128, "description": "商品名"},
			"description":          map[string]any{"type": "string", "maxLength": 20000, "description": "商品説明文"},
			"price":                map[string]any{"type": "number", "minimum": 0, "description": "販売価格（税込・円）"},
			"cost":                 map[string]any{"type": "number", "minimum": 0, "description": "原価（税込・円）"},
			"inventory":            map[string]any{"type": "number", "minimum": 0, "description": "在庫数"},
			"weight":               map[string]any{"type": "number", "minimum": 0, "description": "重量（kg）"},
			"itemUnit":             map[string]any{"type": "string", "maxLength": 16, "description": "単位（例: 個, 瓶, 箱, kg）"},
			"itemDescription":      map[string]any{"type": "string", "maxLength": 64, "description": "内容量の説明（例: 1箱3kg入り）"},
			"deliveryType":         map[string]any{"type": "number", "minimum": 1, "maximum": 3, "description": "配送方法: 1=常温便, 2=冷蔵便, 3=冷凍便"},
			"storageMethodType":    map[string]any{"type": "number", "minimum": 1, "maximum": 4, "description": "保存方法: 1=常温保存, 2=冷暗所保存, 3=冷蔵保存, 4=冷凍保存"},
			"expirationDate":       map[string]any{"type": "number", "minimum": 0, "description": "賞味期限（製造日からの日数）"},
			"recommendedPoint1":    map[string]any{"type": "string", "maxLength": 128, "description": "おすすめポイント1"},
			"recommendedPoint2":    map[string]any{"type": "string", "maxLength": 128, "description": "おすすめポイント2"},
			"recommendedPoint3":    map[string]any{"type": "string", "maxLength": 128, "description": "おすすめポイント3"},
			"originPrefectureCode": map[string]any{"type": "number", "minimum": 1, "maximum": 47, "description": "原産地の都道府県コード（1=北海道, 2=青森, ..., 47=沖縄）"},
			"originCity":           map[string]any{"type": "string", "maxLength": 32, "description": "原産地の市区町村"},
			"scope":                map[string]any{"type": "number", "minimum": 1, "maximum": 3, "description": "公開範囲: 1=全体公開, 2=LINE限定, 3=下書き"},
			"box60Rate":            map[string]any{"type": "number", "minimum": 0, "maximum": 600, "description": "60サイズ箱の占有率（%）"},
			"box80Rate":            map[string]any{"type": "number", "minimum": 0, "maximum": 250, "description": "80サイズ箱の占有率（%）"},
			"box100Rate":           map[string]any{"type": "number", "minimum": 0, "maximum": 100, "description": "100サイズ箱の占有率（%）"},
		},
	},
	{
		Name: "suggestDescription",
		Description: `商品の特徴から魅力的な商品説明文を生成する。
現在のフォーム情報や会話内容をもとに、商品の魅力が伝わる説明文を作成する。`,
		InputSchema: map[string]interface{}{
			"description": map[string]any{"type": "string", "maxLength": 20000, "description": "生成された商品説明文"},
		},
	},
	{
		Name: "suggestPoints",
		Description: `商品のおすすめポイントを3つ生成する。
各ポイントは128文字以内で、商品の魅力を端的に伝える内容にする。`,
		InputSchema: map[string]interface{}{
			"point1": map[string]any{"type": "string", "maxLength": 128, "description": "おすすめポイント1"},
			"point2": map[string]any{"type": "string", "maxLength": 128, "description": "おすすめポイント2"},
			"point3": map[string]any{"type": "string", "maxLength": 128, "description": "おすすめポイント3"},
		},
	},
}

func (h *handler) aiChatRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/ai", h.authentication)
	r.POST("/chat", h.AiChat)
	r.GET("/chat/sessions", h.ListAiChatSessions)
	r.GET("/chat/sessions/:sessionId/messages", h.ListAiChatMessages)
}

func (h *handler) AiChat(ctx *gin.Context) {
	if h.anthropic == nil {
		h.httpError(ctx, status.Error(codes.Unavailable, "AI chat is not available: ANTHROPIC_API_KEY is not configured"))
		return
	}

	req := &types.AiChatRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}

	adminID := getAdminID(ctx)

	// セッション取得 or 作成
	sessionID := req.SessionID
	if sessionID == "" {
		title := ""
		if len(req.Messages) > 0 {
			for _, part := range req.Messages[0].Parts {
				if part.Type == "text" && part.Text != "" {
					title = part.Text
					runes := []rune(title)
					if len(runes) > 50 {
						title = string(runes[:50])
					}
					break
				}
			}
		}
		in := &store.CreateAiChatSessionInput{
			AdminID: adminID,
			Title:   title,
		}
		session, err := h.store.CreateAiChatSession(ctx, in)
		if err != nil {
			h.httpError(ctx, err)
			return
		}
		sessionID = session.ID
	}

	// ユーザーメッセージを DB に保存
	if len(req.Messages) > 0 {
		lastMsg := req.Messages[len(req.Messages)-1]
		if lastMsg.Role == "user" {
			contentJSON, err := json.Marshal(lastMsg.Parts)
			if err != nil {
				h.badRequest(ctx, err)
				return
			}
			msgIn := &store.CreateAiChatMessageInput{
				SessionID: sessionID,
				Role:      "user",
				Content:   string(contentJSON),
			}
			if _, err := h.store.CreateAiChatMessage(ctx, msgIn); err != nil {
				slog.WarnContext(ctx, "Failed to save user message", slog.String("error", err.Error()))
			}
		}
	}

	// UIMessage → Anthropic メッセージ変換
	anthropicMessages := convertToAnthropicMessages(req.Messages)

	// システムプロンプト構築
	systemPrompt := buildSystemPrompt(req.FormData)

	// Anthropic ストリーミング呼び出し
	streamParams := &anthropic.StreamParams{
		System:   systemPrompt,
		Messages: anthropicMessages,
		Tools:    aiChatTools,
	}
	stream := h.anthropic.StreamMessage(ctx.Request.Context(), streamParams)
	defer stream.Close()

	// レスポンスヘッダ設定 (Data Stream Protocol v2)
	ctx.Header("Content-Type", "text/plain; charset=utf-8")
	ctx.Header("X-Vercel-AI-Data-Stream", "v2")
	ctx.Header("Cache-Control", "no-cache")
	ctx.Header("Connection", "keep-alive")
	ctx.Header("X-Session-Id", sessionID)
	ctx.Writer.WriteHeader(http.StatusOK)
	ctx.Writer.Flush()

	// ストリーミングイベント → Data Stream Protocol パート変換
	var assistantContent []map[string]interface{}

	for stream.Next() {
		event := stream.Current()
		parts := convertToDataStreamParts(event, &assistantContent)
		for _, part := range parts {
			if _, err := fmt.Fprint(ctx.Writer, part); err != nil {
				slog.WarnContext(ctx, "Failed to write stream part", slog.String("error", err.Error()))
				return
			}
			ctx.Writer.Flush()
		}
	}

	if err := stream.Err(); err != nil {
		slog.ErrorContext(ctx, "Anthropic stream error", slog.String("error", err.Error()))
		errPart := fmt.Sprintf("3:%s\n", mustMarshalJSON(err.Error()))
		fmt.Fprint(ctx.Writer, errPart)
		ctx.Writer.Flush()
	}

	// アシスタントメッセージを DB に保存
	if len(assistantContent) > 0 {
		contentJSON, err := json.Marshal(assistantContent)
		if err == nil {
			msgIn := &store.CreateAiChatMessageInput{
				SessionID: sessionID,
				Role:      "assistant",
				Content:   string(contentJSON),
			}
			if _, err := h.store.CreateAiChatMessage(ctx, msgIn); err != nil {
				slog.WarnContext(ctx, "Failed to save assistant message", slog.String("error", err.Error()))
			}
		}
	}
}

func (h *handler) ListAiChatSessions(ctx *gin.Context) {
	adminID := getAdminID(ctx)

	in := &store.ListAiChatSessionsInput{
		AdminID: adminID,
		Limit:   50,
		Offset:  0,
	}
	sessions, err := h.store.ListAiChatSessions(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	res := &types.AiChatSessionsResponse{
		Sessions: make([]types.AiChatSessionResponse, len(sessions)),
	}
	for i, s := range sessions {
		res.Sessions[i] = types.AiChatSessionResponse{
			SessionID: s.ID,
			Title:     s.Title,
			CreatedAt: s.CreatedAt.Unix(),
		}
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) ListAiChatMessages(ctx *gin.Context) {
	sessionID := ctx.Param("sessionId")
	if sessionID == "" {
		h.badRequest(ctx, fmt.Errorf("sessionId is required"))
		return
	}

	in := &store.ListAiChatMessagesInput{
		SessionID: sessionID,
	}
	messages, err := h.store.ListAiChatMessages(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	res := &types.AiChatMessagesResponse{
		Messages: make([]types.AiChatMessageResponse, len(messages)),
	}
	for i, m := range messages {
		res.Messages[i] = types.AiChatMessageResponse{
			ID:        m.ID,
			SessionID: m.SessionID,
			Role:      m.Role,
			Content:   m.Content,
			CreatedAt: m.CreatedAt.Unix(),
		}
	}
	ctx.JSON(http.StatusOK, res)
}

// convertToAnthropicMessages は UIMessage 形式を Anthropic API のメッセージ形式に変換する
func convertToAnthropicMessages(messages []types.AiChatMessage) []anthropic.Message {
	var result []anthropic.Message

	for _, msg := range messages {
		switch msg.Role {
		case "user":
			var blocks []anthropic.ContentBlock
			for _, part := range msg.Parts {
				if part.Type == "text" && part.Text != "" {
					blocks = append(blocks, anthropic.ContentBlock{
						Type: "text",
						Text: part.Text,
					})
				}
			}
			if len(blocks) > 0 {
				result = append(result, anthropic.Message{
					Role:    "user",
					Content: blocks,
				})
			}

		case "assistant":
			var assistantBlocks []anthropic.ContentBlock
			var toolResults []anthropic.ContentBlock

			for _, part := range msg.Parts {
				switch part.Type {
				case "text":
					if part.Text != "" {
						assistantBlocks = append(assistantBlocks, anthropic.ContentBlock{
							Type: "text",
							Text: part.Text,
						})
					}
				case "tool-invocation":
					if part.State == "result" || part.State == "call" {
						// assistant 側に tool_use を追加
						argsJSON, _ := json.Marshal(part.Args)
						assistantBlocks = append(assistantBlocks, anthropic.ContentBlock{
							Type:  "tool_use",
							ID:    part.ToolInvocationID,
							Name:  part.ToolName,
							Input: string(argsJSON),
						})

						// result がある場合は tool_result メッセージを追加
						if part.State == "result" {
							resultJSON, _ := json.Marshal(part.Result)
							toolResults = append(toolResults, anthropic.ContentBlock{
								Type:      "tool_result",
								ToolUseID: part.ToolInvocationID,
								Content:   string(resultJSON),
							})
						}
					}
				}
			}

			if len(assistantBlocks) > 0 {
				result = append(result, anthropic.Message{
					Role:    "assistant",
					Content: assistantBlocks,
				})
			}

			if len(toolResults) > 0 {
				result = append(result, anthropic.Message{
					Role:    "user",
					Content: toolResults,
				})
			}
		}
	}

	return result
}

// convertToDataStreamParts は Anthropic ストリーミングイベントを Data Stream Protocol パートに変換する
func convertToDataStreamParts(event sdk.MessageStreamEventUnion, assistantContent *[]map[string]interface{}) []string {
	var parts []string

	switch event.Type {
	case "content_block_start":
		cb := event.ContentBlock
		switch cb.Type {
		case "text":
			// テキストブロック開始: 特別なパートは不要
		case "tool_use":
			// Tool Call 開始
			toolCallStart := map[string]string{"toolCallId": cb.ID, "toolName": cb.Name}
			part := fmt.Sprintf("9:%s\n", mustMarshalJSON(toolCallStart))
			parts = append(parts, part)

			*assistantContent = append(*assistantContent, map[string]interface{}{
				"type":      "tool_use",
				"id":        cb.ID,
				"name":      cb.Name,
				"input_raw": "",
			})
		}

	case "content_block_delta":
		delta := event.Delta
		switch delta.Type {
		case "text_delta":
			if delta.Text != "" {
				part := fmt.Sprintf("0:%s\n", mustMarshalJSON(delta.Text))
				parts = append(parts, part)

				// アシスタントコンテンツにテキストを蓄積
				found := false
				for i := len(*assistantContent) - 1; i >= 0; i-- {
					if (*assistantContent)[i]["type"] == "text" {
						(*assistantContent)[i]["text"] = (*assistantContent)[i]["text"].(string) + delta.Text
						found = true
						break
					}
				}
				if !found {
					*assistantContent = append(*assistantContent, map[string]interface{}{
						"type": "text",
						"text": delta.Text,
					})
				}
			}
		case "input_json_delta":
			if delta.PartialJSON != "" {
				// Tool の引数デルタを tool_use エントリに蓄積
				for i := len(*assistantContent) - 1; i >= 0; i-- {
					if (*assistantContent)[i]["type"] == "tool_use" {
						if raw, ok := (*assistantContent)[i]["input_raw"].(string); ok {
							(*assistantContent)[i]["input_raw"] = raw + delta.PartialJSON
						}
						break
					}
				}

				// ツール呼び出しの引数差分: toolCallId は直近の content_block_start から取得
				toolCallID := ""
				for i := len(*assistantContent) - 1; i >= 0; i-- {
					if (*assistantContent)[i]["type"] == "tool_use" {
						if id, ok := (*assistantContent)[i]["id"].(string); ok {
							toolCallID = id
						}
						break
					}
				}
				argsDelta := map[string]string{"toolCallId": toolCallID, "argsTextDelta": delta.PartialJSON}
				part := fmt.Sprintf("b:%s\n", mustMarshalJSON(argsDelta))
				parts = append(parts, part)
			}
		}

	case "message_delta":
		md := event.AsMessageDelta()
		finishReason := mapStopReason(md.Delta.StopReason)
		stepFinish := map[string]interface{}{
			"finishReason": finishReason,
			"usage": map[string]int64{
				"promptTokens":     md.Usage.InputTokens,
				"completionTokens": md.Usage.OutputTokens,
			},
		}
		part := fmt.Sprintf("e:%s\n", mustMarshalJSON(stepFinish))
		parts = append(parts, part)

	case "message_stop":
		part := "d:{\"finishReason\":\"stop\",\"usage\":{}}\n"
		parts = append(parts, part)
	}

	return parts
}

// buildSystemPrompt はフォームデータを含むシステムプロンプトを構築する
func buildSystemPrompt(formData map[string]any) string {
	formDataJSON, _ := json.MarshalIndent(formData, "", "  ")

	return fmt.Sprintf(`あなたは「ふるマル」の商品登録アシスタントです。
管理者が地域特産品を登録・編集する際に、自然言語での指示をフォームデータに変換します。

## あなたの役割
- ユーザーの自然言語による商品説明からフォームフィールドを抽出・提案する
- 商品説明文やおすすめポイントの生成・改善を支援する
- 不足している情報があれば質問して確認する
- 日本の地域特産品・農産物・水産物・加工食品に関する知識を活用する

## 現在のフォーム状態
以下が現在のフォームに入力されている内容です:
%s

## フィールドの説明
- name: 商品名（最大128文字）
- description: 商品説明文（最大20,000文字）
- price: 販売価格（税込・円）
- cost: 原価（税込・円）
- inventory: 在庫数
- weight: 重量（kg）
- itemUnit: 単位（個, 瓶, 箱, kg など）
- itemDescription: 内容量の説明（例: 1箱3kg入り）
- deliveryType: 配送方法（1=常温便, 2=冷蔵便, 3=冷凍便）
- storageMethodType: 保存方法（1=常温, 2=冷暗所, 3=冷蔵, 4=冷凍）
- expirationDate: 賞味期限（製造日からの日数）
- recommendedPoint1〜3: おすすめポイント（各128文字以内）
- originPrefectureCode: 都道府県コード（1=北海道〜47=沖縄）
- originCity: 市区町村（最大32文字）
- scope: 公開範囲（1=全体公開, 2=LINE限定, 3=下書き）
- box60Rate/box80Rate/box100Rate: 箱サイズ占有率

## 都道府県コード対応表
1=北海道, 2=青森県, 3=岩手県, 4=宮城県, 5=秋田県, 6=山形県, 7=福島県,
8=茨城県, 9=栃木県, 10=群馬県, 11=埼玉県, 12=千葉県, 13=東京都, 14=神奈川県,
15=新潟県, 16=富山県, 17=石川県, 18=福井県, 19=山梨県, 20=長野県,
21=岐阜県, 22=静岡県, 23=愛知県, 24=三重県,
25=滋賀県, 26=京都府, 27=大阪府, 28=兵庫県, 29=奈良県, 30=和歌山県,
31=鳥取県, 32=島根県, 33=岡山県, 34=広島県, 35=山口県,
36=徳島県, 37=香川県, 38=愛媛県, 39=高知県,
40=福岡県, 41=佐賀県, 42=長崎県, 43=熊本県, 44=大分県, 45=宮崎県, 46=鹿児島県, 47=沖縄県

## ルール
- 更新するフィールドのみをToolに含めること（変更不要なフィールドは省略）
- 価格は税込の円で扱う
- 配送方法と保存方法は商品の特性から適切に推測する（生鮮食品→冷蔵/冷凍など）
- 説明文を生成する際は商品の魅力が伝わるよう工夫する
- ユーザーが日本語で話しかけたら日本語で返答する
- 情報が不足している場合は推測せず、ユーザーに確認する
- Toolを呼び出す前に、何を更新するか簡潔に説明する`, string(formDataJSON))
}

// mapStopReason は Anthropic の stop_reason を Data Stream Protocol の finishReason にマッピングする
func mapStopReason(reason sdk.StopReason) string {
	switch reason {
	case sdk.StopReasonEndTurn:
		return "stop"
	case sdk.StopReasonToolUse:
		return "tool-calls"
	case sdk.StopReasonMaxTokens:
		return "length"
	default:
		return "stop"
	}
}

// mustMarshalJSON は値を JSON 文字列にマーシャルする
func mustMarshalJSON(v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		return `""`
	}
	return string(b)
}

