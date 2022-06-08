package business

import rbusiness "mio/internal/pkg/repository/business"

var DefaultPointService = PointService{repo: rbusiness.DefaultPointRepository}

type PointService struct {
	repo rbusiness.PointRepository
}
