package request

type CreateGuestLiveCommentRequest struct {
	Comment string `json:"comment" validate:"required,max=200"` // コメント
}
