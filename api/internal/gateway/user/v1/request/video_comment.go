package request

type CreateVideoCommentRequest struct {
	Comment string `json:"comment,omitempty"` // コメント
}
