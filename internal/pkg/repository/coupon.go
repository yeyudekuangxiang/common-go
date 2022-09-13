package repository

import (
	"gorm.io/gorm"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model/entity/coupon"
	"mio/internal/pkg/repository/repotypes"
)

var DefaultCouponRepository ICouponRepository = NewCouponRepository(app.DB)

type ICouponRepository interface {
	CouponListOfOpenid(openid string, couponTypeIds []string) ([]coupon.CouponRes, error)
	FindCoupon(by FindCouponBy) coupon.Coupon
	Save(coupon2 *coupon.Coupon) error
	CreateBatch(list *[]coupon.Coupon) error
	GetPageCouponRecord(do repotypes.GetPageUserCouponTypeDO) ([]coupon.CouponRecord, int64, error)
}

func NewCouponRepository(DB *gorm.DB) CouponRepository {
	return CouponRepository{DB: DB}
}

type CouponRepository struct {
	DB *gorm.DB
}

func (p CouponRepository) CouponListOfOpenid(openid string, couponTypeIds []string) ([]coupon.CouponRes, error) {
	var Coupons []coupon.CouponRes
	if err := p.DB.Table("coupon").Where("Coupon_type_id in (?)", couponTypeIds).Where("openid = ?", openid).Where("redeemed = ?", "true").Find(&Coupons).Error; err != nil {
		return nil, err
	}
	return Coupons, nil
}
func (p CouponRepository) FindCoupon(by FindCouponBy) coupon.Coupon {
	cp := coupon.Coupon{}
	db := p.DB.Model(cp)

	if by.CouponTypeId != "" {
		db.Where("coupon_type_id = ?", by.CouponTypeId)
	}
	if by.CouponId != "" {
		db.Where("coupon_id = ?", by.CouponId)
	}
	if by.OpenId != "" {
		db.Where("openid = ?", by.OpenId)
	}

	err := db.First(&cp).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		panic(err)
	}
	return cp
}
func (p CouponRepository) Save(coupon2 *coupon.Coupon) error {
	return p.DB.Save(coupon2).Error
}
func (p CouponRepository) CreateBatch(list *[]coupon.Coupon) error {
	return p.DB.Create(list).Error
}
func (p CouponRepository) GetPageCouponRecord(do repotypes.GetPageUserCouponTypeDO) ([]coupon.CouponRecord, int64, error) {
	db := p.DB.Table("coupon c").
		Joins("inner join coupon_type ct on c.coupon_type_id = ct.coupon_type_id")
	if do.OpenId != "" {
		db.Where("c.openid = ?", do.OpenId)
	}
	list := make([]coupon.CouponRecord, 0)
	var count int64
	return list, count, db.Count(&count).
		Select("c.id,c.openid,c.update_time,ct.name,c.coupon_id,c.coupon_type_id,ct.product_item_id,ct.product_item_count,ct.point").
		Order("c.update_time desc").
		Offset(do.Offset).
		Limit(do.Limit).
		Find(&list).
		Error
}
