// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

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

var SchemaNotAvailable = Response{
	Message: "schema not available",
}

var CachePurged = Response{
	Message: "cache purged",
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

var MissingAuthHeader = Response{
	Message: "missing authorization header",
}

var MissingAuthSchemeOrToken = Response{
	Message: "missing auth scheme or token",
}

var InvalidAuthScheme = Response{
	Message: "invalid scheme",
}

var InvalidAuthToken = Response{
	Message: "invalid token",
}
