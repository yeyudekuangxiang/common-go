package wxworkqueue

import "mio/pkg/wxwork"

type RobotMessage struct {
	Key     string
	Type    wxwork.MsgType
	Message wxwork.IMessage
}
