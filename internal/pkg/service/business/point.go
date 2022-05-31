package business

import brepo "mio/internal/pkg/repository/business"

var DefultPointService = PointService{repo: brepo.DefaultPointRepository}

type PointService struct {
	repo brepo.PointRepository
}
