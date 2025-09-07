package request

type CreateLiveCommentRequest struct {
	Comment string `json:"comment" validate:"required,max=200"` // コメント
}
