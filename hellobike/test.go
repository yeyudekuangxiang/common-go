package hellobike

func (receiver ResultCode) name() {
	a := "action=hellobike.activity.bikecard&app_id=20200907153742407&biz_content={\"activityId\":\"1296103963453030400\",\"mobilePhone\":\"13661502232\",\"transactionId\":\"202009091447R0001243\"}&utc_timestamp=1599634041750&version=1.075e3747b359246379b2447dfd5090b8a"

	b := "action=hellobike.activity.bikecard&app_id=20200907153742407&biz_content={\"activityId\":\"1296103963453030400\",\"mobilePhone\":\"13661502232\",\"transactionId\":\"202009091447R0001243\"}&utc_timestamp=1599634041750&version=1.075e3747b359246379b2447dfd5090b8a"

	println(a)
	println(b)
}
