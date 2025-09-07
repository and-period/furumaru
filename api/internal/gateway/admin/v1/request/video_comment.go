package request

type UpdateVideoCommentRequest struct {
	Disabled bool `json:"disabled" validate:""` // コメント無効フラグ
}
