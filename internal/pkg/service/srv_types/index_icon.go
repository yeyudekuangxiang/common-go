package srv_types

import "mio/internal/pkg/model/entity"

type UpdateIndexIconDTO struct {
	Id      int64
	Title   string
	RowNum  string
	Sort    int8
	Status  entity.IndexIconStatus
	IsOpen  int8
	Pic     string
	Display string
}
