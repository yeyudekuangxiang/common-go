package baidu

type CommonRespCode struct {
	ErrorCode int    `json:"error_code"`
	ErrorMsg  string `json:"error_msg"`
}

func (result CommonRespCode) IsSuccess() bool {
	return result.ErrorCode == 0
}

type StatusRespCode struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

func (result StatusRespCode) IsSuccess() bool {
	return result.Status == 0
}
