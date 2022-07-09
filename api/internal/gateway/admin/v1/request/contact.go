package request

type UpdateContactRequest struct {
	Status   int32  `json:"status,omitempty"`   // 対応状況
	Priority int32  `json:"priority,omitempty"` // 優先度
	Note     string `json:"note,omitempty"`     // 対応者メモ
}
