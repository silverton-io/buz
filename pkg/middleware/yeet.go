package middleware

import "github.com/gin-gonic/gin"

type yeet struct {
	Msg string
}

func Yeet() gin.HandlerFunc {
	return func(c *gin.Context) {
		// If msg or source doesn't conform... :wave:
		// y := yeet{
		// 	Msg: "yeeted!",
		// }
	}
}
