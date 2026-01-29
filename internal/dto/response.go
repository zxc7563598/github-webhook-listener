package dto

type BaseResponse struct {
	Success bool   `json:"success"`
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Error   any    `json:"error,omitempty"`
}

func Success(code int, data any, message ...string) *BaseResponse {
	msg := "success"
	if len(message) > 0 {
		msg = message[0]
	}
	return &BaseResponse{
		Success: true,
		Code:    code,
		Data:    data,
		Message: msg,
	}
}

func Error(code int, message string, err any) *BaseResponse {
	return &BaseResponse{
		Success: true,
		Code:    code,
		Error:   err,
		Message: message,
	}
}
