package request

type CreateLiveCommentRequest struct {
	Comment string `json:"comment" binding:"required,max=200"` // コメント
}
