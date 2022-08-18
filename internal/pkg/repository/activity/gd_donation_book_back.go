package activity

import (
	"errors"
	"gorm.io/gorm"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model/entity/activity"
)

var _ DonationBookModel = (*customDonationBookModel)(nil)

type (
	// DonationBookModel is an interface to be customized, add more methods here,
	// and implement the added methods in customDonationBookModel.
	DonationBookModel interface {
		donationBookModel
	}

	customDonationBookModel struct {
		*defaultDonationBookModel
	}

	donationBookModel interface {
		Save(data *activity.GDDonationBookRecord) error
		Create(data *activity.GDDonationBookRecord) error
		FindById(id int64) (*activity.GDDonationBookRecord, error)
		FindBy(by FindRecordBy) ([]*activity.GDDonationBookRecord, error)
	}

	defaultDonationBookModel struct {
		*gorm.DB
		table string
	}
)

// NewDonationBookModel returns a model for the database table.
func NewDonationBookModel() DonationBookModel {
	return &customDonationBookModel{
		defaultDonationBookModel: newDonationBookModel(),
	}
}

func newDonationBookModel() *defaultDonationBookModel {
	return &defaultDonationBookModel{
		DB:    &gorm.DB{},
		table: "`gd_donation_book`",
	}
}

func (m *defaultDonationBookModel) Save(data *activity.GDDonationBookRecord) error {
	return m.DB.Save(data).Error
}

func (m *defaultDonationBookModel) Create(data *activity.GDDonationBookRecord) error {
	return m.DB.Create(data).Error
}

func (m *defaultDonationBookModel) FindById(id int64) (*activity.GDDonationBookRecord, error) {
	var resp *activity.GDDonationBookRecord
	result := m.DB.First(resp, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, errors.New("数据不存在")
		}
		panic(result.Error)
	}
	return resp, nil
}

func (m *defaultDonationBookModel) FindBy(by FindRecordBy) ([]*activity.GDDonationBookRecord, error) {
	var resp []*activity.GDDonationBookRecord
	//var GDDonationBookRecord activity.GDDonationBookRecord
	db := app.DB
	if by.UserId > 0 {
		db.Where("user_id = ?", by.UserId)
	}
	result := db.Find(&resp)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, errors.New("数据不存在")
		}
		panic(result.Error)
	}
	return resp, nil
}
