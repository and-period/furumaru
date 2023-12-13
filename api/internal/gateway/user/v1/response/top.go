package response

type TopCommonResponse struct {
	Lives        []*LiveSummary    `json:"lives"`        // 配信中・配信予定のマルシェ一覧
	Archives     []*ArchiveSummary `json:"archives"`     // 過去のマルシェ一覧
	Coordinators []*Coordinator    `json:"coordinators"` // コーディネータ一覧
}
