package request

type CreateGuestLiveCommentRequest struct {
	Comment string `json:"comment" binding:"required,max=200"` // コメント
}
