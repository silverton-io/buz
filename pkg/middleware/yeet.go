package middleware

import "github.com/gin-gonic/gin"

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
