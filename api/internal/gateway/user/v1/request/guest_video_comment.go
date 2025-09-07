package request

type CreateGuestVideoCommentRequest struct {
	Comment string `json:"comment" validate:"required,max=200"` // コメント
}
