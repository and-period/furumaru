package request

type UpdateLiveCommentRequest struct {
	Disabled bool `json:"disabled" binding:""` // コメント無効フラグ
}
