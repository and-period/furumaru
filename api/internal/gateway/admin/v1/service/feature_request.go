package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/types"
	"github.com/and-period/furumaru/api/internal/messenger/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
)

type FeatureRequest struct {
	types.FeatureRequest
}

type FeatureRequests []*FeatureRequest

func NewFeatureRequest(featureRequest *entity.FeatureRequest) *FeatureRequest {
	return &FeatureRequest{
		FeatureRequest: types.FeatureRequest{
			ID:            featureRequest.ID,
			Title:         featureRequest.Title,
			Description:   featureRequest.Description,
			Category:      types.FeatureRequestCategory(featureRequest.Category),
			Priority:      types.FeatureRequestPriority(featureRequest.Priority),
			Status:        types.FeatureRequestStatus(featureRequest.Status),
			Note:          featureRequest.Note,
			SubmittedBy:   featureRequest.SubmittedBy,
			SubmitterName: "",
			CreatedAt:     jst.Unix(featureRequest.CreatedAt),
			UpdatedAt:     jst.Unix(featureRequest.UpdatedAt),
		},
	}
}

func (f *FeatureRequest) Fill(submitterName string) {
	f.SubmitterName = submitterName
}

func (f *FeatureRequest) Response() *types.FeatureRequest {
	return &f.FeatureRequest
}

func NewFeatureRequests(featureRequests entity.FeatureRequests) FeatureRequests {
	res := make(FeatureRequests, len(featureRequests))
	for i := range featureRequests {
		res[i] = NewFeatureRequest(featureRequests[i])
	}
	return res
}

func (fs FeatureRequests) Fill(admins map[string]*Admin) {
	for _, f := range fs {
		if admin, ok := admins[f.SubmittedBy]; ok {
			f.SubmitterName = admin.Name()
		}
	}
}

func (fs FeatureRequests) SubmittedByIDs() []string {
	ids := make([]string, 0, len(fs))
	seen := make(map[string]struct{}, len(fs))
	for _, f := range fs {
		if _, ok := seen[f.SubmittedBy]; !ok {
			seen[f.SubmittedBy] = struct{}{}
			ids = append(ids, f.SubmittedBy)
		}
	}
	return ids
}

func (fs FeatureRequests) Response() []*types.FeatureRequest {
	res := make([]*types.FeatureRequest, len(fs))
	for i := range fs {
		res[i] = fs[i].Response()
	}
	return res
}
