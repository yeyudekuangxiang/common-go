package repotypes

import "mio/internal/pkg/model/entity"

type GetIndexIconOneDO struct {
	ID int64
}

type GetIndexIconExistDO struct {
	Name     string
	ImageUrl string
	NotId    int64
}

type DeleteIndexIconDO struct {
	Id     int64
	Status entity.IndexIconStatus
}

type GetIndexIconPageDO struct {
	Offset int
	Limit  int
	Title  string
	Status entity.IndexIconStatus
	IsOpen int8
}
