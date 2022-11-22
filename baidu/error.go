package baidu

type ErrorResponse struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

func (e ErrorResponse) IsSuccess() bool {
	return e.Error == ""
}
