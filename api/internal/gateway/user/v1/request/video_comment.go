package request

type CreateVideoCommentRequest struct {
	Comment string `json:"comment" binding:"required,max=200"` // コメント
}
