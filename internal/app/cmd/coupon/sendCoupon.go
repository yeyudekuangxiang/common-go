package coupon

import (
	"github.com/spf13/cobra"
	"mio/internal/pkg/core/app"
	mioctx "mio/internal/pkg/core/context"
	"mio/internal/pkg/core/initialize"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
	"mio/internal/pkg/service"
	"mio/internal/pkg/service/platform/jhx"
)

var JhxSendCouponCmd = &cobra.Command{
	Use:     "jhxSendCoupon",
	Aliases: []string{"jhxSendCoupon"},
	Short:   "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		initialize.Initialize("/Users/yunfeng/Documents/workspace/mp2c-go/config.yaml")

		userIds := []int64{
			1535006, 1535024, 1535031, 1535032, 1535041, 1535046, 1535050, 1535063, 1535065, 1535077, 1535087, 1535095, 1535105, 1535107, 1535114, 1535129, 1535130, 1535174, 1535187, 1535190, 1535198, 1535201, 1535216, 1535221, 1535239, 1535245, 1535252, 1535256, 1535257,
			1535259, 1535262, 1535264, 1535287, 1535298, 1535304, 1535334, 1535341, 1535377, 1535390, 1535396, 1535419, 1535420, 1535456, 1535473, 1535483, 1535490, 1535495, 1535498, 1535500, 1535506, 1535507, 1535513, 1535538, 1535548, 1535558, 1535623, 1535681, 1535723,
			1535731, 1535787, 1535807, 1535816, 1535843, 1535881, 1535897, 1535916, 1535925, 1535945, 1536004, 1536014, 1536016, 1536031, 1536035, 1536068, 1536075, 1536134, 1536151, 1536185, 1536187, 1536196, 1536232, 1536248, 1536269, 1536280, 1536342, 1536432, 1536433,
			1536435, 1536439, 1536441, 1536442, 1536443, 1536448, 1536452, 1536457, 1536459, 1536491, 1536504, 1536530, 1536559, 1536579, 1536596, 1536613, 1536614, 1536653, 1536661, 1536708, 1536719, 1536778, 1536802, 1536967, 1537001, 1537040, 1537119, 1537198, 1537402,
			1537437, 1537501, 1537508, 1537562, 1537594, 1537614, 1537675, 1537706, 1537825, 1538066, 1538084, 1538235, 1538701, 1538725, 1539250, 1539263, 1539297, 1539447, 1539461, 1539517, 1539535, 1539791, 1539804, 1539889, 1539980, 1540083, 1540133, 1540524, 1540611,
			1540794, 1541662, 1541923, 1542480,
		}

		uList, err := service.DefaultUserService.GetUserListBy(repository.GetUserListBy{UserIds: userIds})
		if err != nil {
			return
		}

		jhxService := jhx.NewJhxService(mioctx.NewMioContext())
		for _, u := range uList {
			if u.ChannelId == 1059 {
				go func(u entity.User) {
					_, err := jhxService.SendCoupon(1000, u)
					if err != nil {
						app.Logger.Errorf("金华行发券失败:%s", err.Error())
						return
					}
					return
				}(u)
			}
		}
	},
}
