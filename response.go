package main

const (
	ResponseStateSuccess = "success"
	ResponseStateError   = "error"
)

type Response struct {
	State  string      `json:"state"`
	Result interface{} `json:"result"`
}

func NewSuccessResponse(result interface{}) *Response {
	return &Response{
		State:  ResponseStateSuccess,
		Result: result,
	}
}

func NewErrorResponse(err *Error) *Response {
	return &Response{
		State:  ResponseStateError,
		Result: err,
	}
}
