package middlewares

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func StripPrefix(prefix string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Request.URL.Path = strings.TrimPrefix(ctx.Request.URL.Path, prefix)
		ctx.Next()
	}
}
