// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the AGPLv3 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package middleware

import "github.com/gin-gonic/gin"

// nolint: unused
type yeet struct {
	Msg string
}

func Yeet() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO
		// If msg or source doesn't conform to something... :wave:
		// y := yeet{
		// 	Msg: "yeeted!",
		// }
	}
}
