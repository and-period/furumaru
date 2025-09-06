package request

type CreateGuestVideoCommentRequest struct {
	Comment string `json:"comment" binding:"required,max=200"` // コメント
}
