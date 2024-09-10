package entity

import (
	"encoding/json"
	"time"

	"github.com/and-period/furumaru/api/pkg/jst"
	"gorm.io/datatypes"
)

// OrderExperience - 注文体験情報
type OrderExperience struct {
	OrderID               string                 `gorm:"primaryKey;<-:create"`        // 注文履歴ID
	ExperienceRevisionID  int64                  `gorm:"primaryKey;<-:create"`        // 体験ID
	AdultCount            int64                  `gorm:""`                            // 大人人数
	JuniorHighSchoolCount int64                  `gorm:""`                            // 中学生人数
	ElementarySchoolCount int64                  `gorm:""`                            // 小学生人数
	PreschoolCount        int64                  `gorm:""`                            // 幼児人数
	SeniorCount           int64                  `gorm:""`                            // シニア人数
	Remarks               OrderExperienceRemarks `gorm:"-"`                           // 備考
	RemarksJSON           datatypes.JSON         `gorm:"default:null;column:remarks"` // 備考(JSON)
	CreatedAt             time.Time              `gorm:"<-:create"`                   // 登録日時
	UpdatedAt             time.Time              `gorm:""`                            // 更新日時
}

// OrderExperienceRemarks - 注文体験情報詳細
type OrderExperienceRemarks struct {
	Transportation string    `json:"transportation"` // 交通手段
	RequestedDate  time.Time `json:"requestedDate"`  // 体験希望日
	RequestedTime  time.Time `json:"requestedTime"`  // 体験希望時間
}

type OrderExperiences []*OrderExperience

type NewOrderExperienceParams struct {
	OrderID               string
	Experience            *Experience
	AdultCount            int64
	JuniorHighSchoolCount int64
	ElementarySchoolCount int64
	PreschoolCount        int64
	SeniorCount           int64
	Transportation        string
	RequestedDate         string
	RequestedTime         string
}

type NewOrderExperienceRemarksParams struct {
	Transportation string
	RequestedDate  string
	RequestedTime  string
}

func NewOrderExperience(params *NewOrderExperienceParams) (*OrderExperience, error) {
	rparams := &NewOrderExperienceRemarksParams{
		Transportation: params.Transportation,
		RequestedDate:  params.RequestedDate,
		RequestedTime:  params.RequestedTime,
	}
	remarks, err := NewOrderExperienceRemarks(rparams)
	if err != nil {
		return nil, err
	}
	return &OrderExperience{
		OrderID:               params.OrderID,
		ExperienceRevisionID:  params.Experience.ExperienceRevision.ID,
		AdultCount:            params.AdultCount,
		JuniorHighSchoolCount: params.JuniorHighSchoolCount,
		ElementarySchoolCount: params.ElementarySchoolCount,
		PreschoolCount:        params.PreschoolCount,
		SeniorCount:           params.SeniorCount,
		Remarks:               *remarks,
	}, nil
}

func (o *OrderExperience) Fill() error {
	remarks, err := o.unmarshalRemarks()
	if err != nil {
		return err
	}
	o.Remarks = *remarks
	return nil
}

func (o *OrderExperience) unmarshalRemarks() (*OrderExperienceRemarks, error) {
	if o.RemarksJSON == nil {
		return &OrderExperienceRemarks{}, nil
	}
	var remarks *OrderExperienceRemarks
	return remarks, json.Unmarshal(o.RemarksJSON, &remarks)
}

func (o *OrderExperience) FillJSON() error {
	remarks, err := o.Remarks.Marshal()
	if err != nil {
		return err
	}
	o.RemarksJSON = remarks
	return nil
}

func (os OrderExperiences) MapByOrderID() map[string]*OrderExperience {
	m := make(map[string]*OrderExperience)
	for _, o := range os {
		m[o.OrderID] = o
	}
	return m
}

func (os OrderExperiences) Fill() error {
	for _, o := range os {
		if err := o.Fill(); err != nil {
			return err
		}
	}
	return nil
}

func NewOrderExperienceRemarks(params *NewOrderExperienceRemarksParams) (*OrderExperienceRemarks, error) {
	var (
		requestedDate, requestedTime time.Time
		err                          error
	)
	if params.RequestedDate != "" {
		requestedDate, err = jst.ParseFromYYYYMMDD(params.RequestedDate)
		if err != nil {
			return nil, err
		}
	}
	if params.RequestedTime != "" {
		requestedTime, err = jst.ParseFromHHMM(params.RequestedTime)
		if err != nil {
			return nil, err
		}
	}
	return &OrderExperienceRemarks{
		Transportation: params.Transportation,
		RequestedDate:  requestedDate,
		RequestedTime:  requestedTime,
	}, nil
}

func (r *OrderExperienceRemarks) Marshal() ([]byte, error) {
	if len(r.Transportation) == 0 {
		return nil, nil
	}
	return json.Marshal(r)
}
