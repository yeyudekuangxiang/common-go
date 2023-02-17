package wxworkmsg

import "gitlab.miotech.com/miotech-application/backend/common-go/wxwork"

type RobotMessage struct {
	Key     string
	Type    wxwork.MsgType
	Message wxwork.IMessage
}
