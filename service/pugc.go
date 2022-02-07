package service

import (
	Pugc "mio/model/pugc"
	"mio/repository"
)

var DefaultPugcService = NewPugcService(repository.DefaultPugcRepository)

func NewPugcService(r repository.IPugcRepository) PugcService {
	return PugcService{
		r: r,
	}
}

type PugcService struct {
	r repository.IPugcRepository
}

func (u PugcService) InsertPugc(pugc *Pugc.PugcAddModel) error {
	return u.r.Insert(pugc)
}
