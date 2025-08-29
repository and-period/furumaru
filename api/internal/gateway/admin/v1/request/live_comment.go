package request

type UpdateLiveCommentRequest struct {
	Disabled bool `json:"disabled"` // コメント無効フラグ
}
