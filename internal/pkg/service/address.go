package service

import (
	"gorm.io/gorm"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model/entity"
)

var DefaultAddressService = AddressService{}

type AddressService struct {
}

func (srv AddressService) FindDefaultAddress(openid string) (*entity.Address, error) {
	address := entity.Address{}
	err := app.DB.Where("openid = ? and is_default = true", openid).Take(&address).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		panic(err)
	}
	return &address, nil
}
