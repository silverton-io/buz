// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the AGPLv3 license, a copy of
// which may be found at https://github.com/silverton-io/honeypot/blob/main/LICENSE

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

var ManifoldDistributionError = Response{
	Message: "distribution error",
}
