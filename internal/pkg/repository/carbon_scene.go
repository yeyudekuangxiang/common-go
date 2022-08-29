package repository

import (
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository/repotypes"
)

func NewCarbonSceneRepository(ctx *context.MioContext) CarbonSceneRepository {
	return CarbonSceneRepository{ctx: ctx}
}

type CarbonSceneRepository struct {
	ctx *context.MioContext
}

func (repo CarbonSceneRepository) Save(transaction *entity.CarbonScene) error {
	return repo.ctx.DB.Save(transaction).Error
}

func (repo CarbonSceneRepository) Create(transaction *entity.CarbonScene) error {
	return repo.ctx.DB.Create(transaction).Error
}

func (repo CarbonSceneRepository) FindBy(by repotypes.CarbonSceneBy) entity.CarbonScene {
	pt := entity.CarbonScene{}
	db := repo.ctx.DB.Model(pt)
	if by.Type != "" {
		db.Where("type", by.Type)
	}
	err := db.First(&pt).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		panic(err)
	}
	return pt
}

//GetValue 根据类型和值获取碳量
func (repo CarbonSceneRepository) GetValue(ret entity.CarbonScene, initVal float64) float64 {
	//ret := repo.FindBy(repotypes.CarbonSceneBy{Type: sceneType})
	if ret.ID != 0 {
		unitNumerator := decimal.NewFromFloat(ret.UnitNumerator)
		unitDenominator := decimal.NewFromFloat(ret.UnitDenominator)
		initValDec := decimal.NewFromFloat(initVal)
		val, _ := initValDec.Mul(unitNumerator).Div(unitDenominator).Round(2).Float64()
		return val
	}
	return 0
}
