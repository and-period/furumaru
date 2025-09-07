package request

type CreateVideoCommentRequest struct {
	Comment string `json:"comment" validate:"required,max=200"` // コメント
}
