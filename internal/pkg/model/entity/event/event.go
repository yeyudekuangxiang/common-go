package event

import (
	"errors"
	"fmt"
	"mio/internal/pkg/model"
	"mio/internal/pkg/model/entity"
	"strconv"
	"strings"
	"time"
)

var OrderByEventSortDesc entity.OrderBy = "order_by_event_sort_desc"

type Event struct {
	ID                    int                  `json:"-" gorm:"primaryKey;type:serial4;not null;comment:公益活动表"`
	EventCategoryId       string               `json:"eventCategoryId" gorm:"type:varchar(255);not null;comment:公益活动所属分类标识"`
	EventId               string               `json:"eventId" gorm:"type:varchar(255);not null;comment:公益活动标识"`
	EventTemplateType     EventTemplateType    `json:"eventTemplateType" gorm:"type:varchar(50);not null;default:'';comment:证书模版类型"`
	Title                 string               `json:"title" gorm:"type:varchar(255);not null;comment:公益活动标题"`
	Subtitle              string               `json:"subtitle" gorm:"type:varchar(255);not null;default:'';comment:公益活动副标题"`
	Active                bool                 `json:"active" gorm:"type:bool;not null;default:false;comment:是否上线"`
	CoverImageUrl         string               `json:"coverImageUrl" gorm:"type:varchar(500);not null;comment:项目主图 343 × 200"`
	StartTime             model.Time           `json:"startTime" gorm:"type:timestamptz;not null;comment:公益活动开始时间"`
	EndTime               model.Time           `json:"endTime" gorm:"type:timestamptz;comment:公益活动结束时间"`
	ProductItemId         string               `json:"productItemId" gorm:"type:varchar(255);not null;comment:关联的商品编号"`
	ParticipationCount    int                  `json:"participationCount" gorm:"type:int4;not null;default:0;comment:已参与次数"`
	ParticipationTitle    string               `json:"ParticipationTitle" gorm:"type:varchar(255);"`
	ParticipationSubtitle string               `json:"participationSubtitle" gorm:"type:varchar(255);default:'';comment:用于展示支持次数或者co2"`
	Sort                  int                  `json:"sort" gorm:"type:int4;default:0;comment:排序 从小到大排序"`
	Tag                   model.ArrayString    `json:"tag" gorm:"type:varchar(255);default:'';comment:标签,多个标签用英文逗号隔开"`
	TemplateSetting       EventTemplateSetting `json:"templateSetting" gorm:"type:varchar(2000);not null;default:'';comment:公益活动模版配置"`
	Limit                 EventLimit           `json:"limit" gorm:"comment:空表示不限制1-1D 按天限制次数1-1W 按周限制次数1-1M 按月限制次数1-1Y 按年限制次数"`
}

func (Event) TableName() string {
	return "event"
}

type EventLimit string

func (e EventLimit) Parse() (time.Duration, int64, error) {
	var err error

	if e == "" {
		return 0, 0, nil
	}
	sl := strings.Split(string(e), "-")
	fmt.Println(sl)
	if len(sl) != 2 {
		return 0, 0, errors.New("event limit format err")
	}
	count, err := strconv.ParseInt(sl[0], 10, 64)
	if err != nil {
		return 0, 0, errors.New("event limit format err")
	}
	fmt.Println(count)

	var tNum int64 = 1
	if len(sl[1]) > 1 {
		tNum, err = strconv.ParseInt(sl[1][:len(sl[1])-1], 10, 64)
		if err != nil {
			return 0, 0, errors.New("event limit format err")
		}
		if tNum <= 0 {
			return 0, 0, errors.New("event limit format err")
		}
	}
	fmt.Println(tNum)

	tFlag := sl[1][len(sl[1])-1:]
	fmt.Println(tNum, tFlag)
	now := time.Now()
	switch strings.ToUpper(tFlag) {
	case "D":
		return now.AddDate(0, 0, int(tNum)).Sub(now), count, nil
	case "W":
		return now.AddDate(0, 0, int(tNum)*7).Sub(now), count, nil
	case "M":
		return now.AddDate(0, int(tNum), 0).Sub(now), count, nil
	case "Y":
		return now.AddDate(int(tNum), 0, 0).Sub(now), count, nil
	}
	return 0, 0, errors.New("event limit format err")
}
