package request

type UpdateVideoCommentRequest struct {
	Disabled bool `json:"disabled"` // コメント無効フラグ
}
