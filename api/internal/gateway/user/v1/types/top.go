package types

type TopCommonResponse struct {
	Lives            []*LiveSummary    `json:"lives"`            // 配信中・配信予定のマルシェ一覧
	Archives         []*ArchiveSummary `json:"archives"`         // 過去のマルシェ一覧
	ProductVideos    []*VideoSummary   `json:"productVideos"`    // 商品動画一覧
	ExperienceVideos []*VideoSummary   `json:"experienceVideos"` // 体験動画一覧
	Coordinators     []*Coordinator    `json:"coordinators"`     // コーディネータ一覧
}
