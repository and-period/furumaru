package types

// AiChatRequest - AIチャットリクエスト
type AiChatRequest struct {
	SessionID string          `json:"sessionId"`  // セッションID（空なら新規作成）
	Messages  []AiChatMessage `json:"messages"`   // UIMessage 形式
	FormData  map[string]any  `json:"formData"`   // 現在のフォーム状態
}

// AiChatMessage - UIMessage 形式のメッセージ
type AiChatMessage struct {
	ID    string       `json:"id"`
	Role  string       `json:"role"`
	Parts []AiChatPart `json:"parts"`
}

// AiChatPart - メッセージの各パート
type AiChatPart struct {
	Type             string `json:"type"`                       // "text", "tool-invocation"
	Text             string `json:"text,omitempty"`             // type=text の場合
	ToolInvocationID string `json:"toolInvocationId,omitempty"` // type=tool-invocation
	ToolName         string `json:"toolName,omitempty"`
	State            string `json:"state,omitempty"` // "call", "result"
	Args             any    `json:"args,omitempty"`
	Result           any    `json:"result,omitempty"`
}

// AiChatSessionsResponse - セッション一覧レスポンス
type AiChatSessionsResponse struct {
	Sessions []AiChatSessionResponse `json:"sessions"`
}

// AiChatSessionResponse - セッション情報レスポンス
type AiChatSessionResponse struct {
	SessionID string `json:"sessionId"`
	Title     string `json:"title"`
	CreatedAt int64  `json:"createdAt"`
}

// AiChatMessagesResponse - メッセージ一覧レスポンス
type AiChatMessagesResponse struct {
	Messages []AiChatMessageResponse `json:"messages"`
}

// AiChatMessageResponse - メッセージ情報レスポンス
type AiChatMessageResponse struct {
	ID        string `json:"id"`
	SessionID string `json:"sessionId"`
	Role      string `json:"role"`
	Content   string `json:"content"`
	CreatedAt int64  `json:"createdAt"`
}
