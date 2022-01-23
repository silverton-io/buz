package main

type Response struct {
	Message string `json:"message"`
}

var AllOk = Response{
	Message: "ok",
}
