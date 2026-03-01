package types

// AuditLog - 監査ログ情報
type AuditLog struct {
	ID            string   `json:"id"`            // 監査ログID
	CreatedAt     int64    `json:"createdAt"`     // 操作日時
	AdminID       string   `json:"adminId"`       // 操作者ID
	AdminType     int32    `json:"adminType"`     // 操作者種別
	Action        int32    `json:"action"`        // 操作種別
	ResourceType  string   `json:"resourceType"`  // リソース種別
	ResourceID    string   `json:"resourceId"`    // リソースID
	Result        int32    `json:"result"`        // 結果
	ResultDetail  string   `json:"resultDetail"`  // エラー詳細
	HttpMethod    string   `json:"httpMethod"`    // HTTPメソッド
	HttpPath      string   `json:"httpPath"`      // HTTPパス
	HttpRoute     string   `json:"httpRoute"`     // HTTPルート
	HttpStatus    int      `json:"httpStatus"`    // HTTPステータス
	ClientIP      string   `json:"clientIp"`      // 接続元IP
	UserAgent     string   `json:"userAgent"`     // User-Agent
	RequestID     string   `json:"requestId"`     // リクエストID
	DurationMs    int      `json:"durationMs"`    // 処理時間(ms)
	ChangedFields []string `json:"changedFields"` // 変更フィールド名一覧
}

// AuditLogsResponse - 監査ログ一覧レスポンス
type AuditLogsResponse struct {
	AuditLogs []*AuditLog `json:"auditLogs"` // 監査ログ一覧
	Total     int64       `json:"total"`      // 合計件数
}
