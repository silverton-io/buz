package response

type Response struct {
	Message string `json:"message"`
}

var Ok = Response{
	Message: "ok",
}
