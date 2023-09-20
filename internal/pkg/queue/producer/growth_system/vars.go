package growth_system

import (
	"encoding/json"
	"fmt"
	"gitlab.miotech.com/miotech-application/backend/common-go/tool/timetool"
	"mio/internal/pkg/queue/types/message/growthsystemmsg"
	"mio/internal/pkg/util"
	"strconv"
	"strings"
	"time"
)

var eventMap = map[string]string{
	"invite":                    "invite",
	"check_in":                  "check_in",
	"quiz":                      "quiz",
	"step":                      "step",
	"coffee_cup":                "coffee_cup",
	"bike_ride":                 "ride",
	"cycling":                   "ride",
	"jhx":                       "bus",
	"ytx":                       "subway",
	"ecar":                      "recharge",
	"fast_electricity":          "recharge",
	"power_replace":             "battery_swapping",
	"recharge_mio":              "recharge_mio",
	"recycling":                 "recycling",
	"recycling_aihuishou":       "recycling",
	"recycling_clothing":        "recycling",
	"recycling_computer":        "recycling",
	"recycling_appliance":       "recycling",
	"recycling_book":            "recycling",
	"recycling_fmy_clothing":    "recycling",
	"recycling_shishanghuishou": "recycling",
	"recycling_dangdangyixia":   "recycling",
	"reduce_plastic":            "reduce_plastic",
	"secondhand_order":          "second_hand_market",
	"hello_bike_ride":           "hello_bike",
	"hello_bike":                "hello_bike",
	"bottle_recycling":          "bottle_recycling",
	"bottle":                    "bottle_recycling",
	"like":                      "post_like",
	"article":                   "post_push",
	"comment":                   "post_comment",
	"charitable":                "charitable",
	"mall_redemption":           "mall_redemption",
}

func getEvent(paramsType string) string {
	if taskType, ok := eventMap[paramsType]; ok {
		return taskType
	}
	return ""
}

func checkParams(params growthsystemmsg.GrowthSystemParam) ([]byte, error) {
	if params.TaskType == "" && params.TaskSubType == "" {
		return nil, fmt.Errorf("TaskType and TaskSubType cannot be empty")
	}

	if params.TaskSubType == "" {
		params.TaskSubType = params.TaskType
	}

	if params.TaskType == "" {
		taskType := getEvent(strings.ToLower(params.TaskSubType))
		if taskType == "" {
			params.TaskType = params.TaskSubType
		} else {
			params.TaskType = taskType
		}
	}

	if params.TaskValue == 0 {
		params.TaskValue = 1
	}

	id, err := util.SnowflakeID()
	if err != nil {
		return nil, err
	}

	msg := &growthsystemmsg.GrowthSystemReq{
		MessageId:   id.String(),
		TaskType:    params.TaskType,
		TaskSubType: params.TaskSubType,
		UserId:      params.UserId,
		TaskValue:   strconv.FormatInt(params.TaskValue, 10),
		Time:        time.Now().Format(timetool.DateFormat),
	}

	marshal, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}
	return marshal, nil
}
