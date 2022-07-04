package activity

import (
	"mio/internal/pkg/repository/activity"
)

type ServiceContext struct {
	DonationBookModel activity.DonationBookModel
}

func NewServiceContext() *ServiceContext {
	return &ServiceContext{
		DonationBookModel: activity.NewDonationBookModel(),
	}
}
