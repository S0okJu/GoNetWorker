package api

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func NewResponse(code Code, message string) Response {
	return Response{
		Code:    code.Number(),
		Message: message,
	}
}
