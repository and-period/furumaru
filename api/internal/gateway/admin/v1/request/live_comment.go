package request

type UpdateLiveCommentRequest struct {
	Disabled bool `json:"disabled,omitempty"` // コメント無効フラグ
}
