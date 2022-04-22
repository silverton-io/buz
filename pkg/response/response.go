package response

type Response struct {
	Message string `json:"message"`
}

var Ok = Response{
	Message: "ok",
}

var InvalidContentType = Response{
	Message: "invalid content type",
}

var BadRequest = Response{
	Message: "bad request",
}

var SchemaNotCached = Response{
	Message: "schema not cached",
}

var Timeout = Response{
	Message: "request timed out",
}

var RateLimitExceeded = Response{
	Message: "rate limit exceeded",
}

var DistributionError = Response{
	Message: "distribution error",
}
