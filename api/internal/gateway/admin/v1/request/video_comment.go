package request

type UpdateVideoCommentRequest struct {
	Disabled bool `json:"disabled" binding:""` // コメント無効フラグ
}
