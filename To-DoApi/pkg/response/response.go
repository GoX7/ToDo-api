package response

type Response struct {
	Status  string `json:"status"`
	Error   string `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
}

func NewOk() Response {
	return Response{
		Status: "OK",
	}
}

func NewError(msg string) Response {
	return Response{
		Status: "Error",
		Error:  msg,
	}
}

func NewMessage(msg string) Response {
	return Response{
		Status:  "OK",
		Message: msg,
	}
}
