package entity

import (
	"time"

	"github.com/and-period/furumaru/api/pkg/uuid"
)

// AuditLogAction - 監査ログ操作種別
type AuditLogAction int32

const (
	AuditLogActionUnknown AuditLogAction = 0
	AuditLogActionCreate  AuditLogAction = 1 // 作成
	AuditLogActionUpdate  AuditLogAction = 2 // 更新
	AuditLogActionDelete  AuditLogAction = 3 // 削除
	AuditLogActionSignIn  AuditLogAction = 4 // サインイン
	AuditLogActionSignOut AuditLogAction = 5 // サインアウト
	AuditLogActionExport  AuditLogAction = 6 // エクスポート
	AuditLogActionUpload  AuditLogAction = 7 // アップロード
)

// AuditLogResult - 監査ログ結果
type AuditLogResult int32

const (
	AuditLogResultUnknown AuditLogResult = 0
	AuditLogResultSuccess AuditLogResult = 1 // 成功
	AuditLogResultFailure AuditLogResult = 2 // 失敗
	AuditLogResultDenied  AuditLogResult = 3 // 拒否
	AuditLogResultError   AuditLogResult = 4 // エラー
)

// AuditLog - 監査ログ
type AuditLog struct {
	ID            string         `gorm:"primaryKey;<-:create"` // 監査ログID
	CreatedAt     time.Time      `gorm:"<-:create"`            // 操作日時
	AdminID       string         `gorm:"<-:create"`            // 操作者ID
	AdminType     AdminType      `gorm:"<-:create"`            // 操作者種別
	Action        AuditLogAction `gorm:"<-:create"`            // 操作種別
	ResourceType  string         `gorm:"<-:create"`            // リソース種別
	ResourceID    string         `gorm:"<-:create"`            // リソースID
	Result        AuditLogResult `gorm:"<-:create"`            // 結果
	ResultDetail  string         `gorm:"<-:create"`            // エラー詳細
	HttpMethod    string         `gorm:"<-:create"`            // HTTPメソッド
	HttpPath      string         `gorm:"<-:create"`            // HTTPパス
	HttpRoute     string         `gorm:"<-:create"`            // HTTPルート
	HttpStatus    int            `gorm:"<-:create"`            // HTTPステータス
	ClientIP      string         `gorm:"<-:create"`            // 接続元IP
	UserAgent     string         `gorm:"<-:create"`            // User-Agent
	RequestID     string         `gorm:"<-:create"`            // リクエストID
	DurationMs    int            `gorm:"<-:create"`            // 処理時間(ms)
	ChangedFields []string       `gorm:"serializer:json"`      // 変更フィールド名一覧
	Metadata      []byte         `gorm:"column:metadata"`      // 追加メタデータ
	UpdatedAt     time.Time      `gorm:""`                     // GORM互換用
}

type AuditLogs []*AuditLog

type NewAuditLogParams struct {
	AdminID       string
	AdminType     AdminType
	Action        AuditLogAction
	ResourceType  string
	ResourceID    string
	Result        AuditLogResult
	ResultDetail  string
	HttpMethod    string
	HttpPath      string
	HttpRoute     string
	HttpStatus    int
	ClientIP      string
	UserAgent     string
	RequestID     string
	DurationMs    int
	ChangedFields []string
}

func NewAuditLog(params *NewAuditLogParams) *AuditLog {
	return &AuditLog{
		ID:            uuid.Base58Encode(uuid.New()),
		AdminID:       params.AdminID,
		AdminType:     params.AdminType,
		Action:        params.Action,
		ResourceType:  params.ResourceType,
		ResourceID:    params.ResourceID,
		Result:        params.Result,
		ResultDetail:  params.ResultDetail,
		HttpMethod:    params.HttpMethod,
		HttpPath:      params.HttpPath,
		HttpRoute:     params.HttpRoute,
		HttpStatus:    params.HttpStatus,
		ClientIP:      params.ClientIP,
		UserAgent:     params.UserAgent,
		RequestID:     params.RequestID,
		DurationMs:    params.DurationMs,
		ChangedFields: params.ChangedFields,
	}
}

func (l *AuditLog) TableName() string {
	return "audit_logs"
}
