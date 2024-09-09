package request

type UpdateVideoCommentRequest struct {
	Disabled bool `json:"disabled,omitempty"` // コメント無効フラグ
}
