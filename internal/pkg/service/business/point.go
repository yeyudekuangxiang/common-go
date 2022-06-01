package business

import brepo "mio/internal/pkg/repository/business"

var DefaultPointService = PointService{repo: brepo.DefaultPointRepository}

type PointService struct {
	repo brepo.PointRepository
}
