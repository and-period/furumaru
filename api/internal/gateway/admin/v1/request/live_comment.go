package request

type UpdateLiveCommentRequest struct {
	Disabled bool `json:"disabled" validate:""` // コメント無効フラグ
}
