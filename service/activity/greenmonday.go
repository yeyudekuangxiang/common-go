package activity

type GreenMondayService struct {
}

// EnableOrder 检测用户是否有权限下单
func (srv GreenMondayService) EnableOrder(userId int64) error {
	return nil
}
