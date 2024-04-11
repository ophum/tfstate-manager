package middlewares

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func SessionAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session := sessions.Default(ctx)

		userID, ok := session.Get("user_id").(uint64)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "unauthorized",
			})
			return
		}

		ctx.Request = ctx.Request.WithContext(
			intoContextUserID(ctx.Request.Context(), userID),
		)
		ctx.Next()
	}
}

type userIDKey struct{}

func intoContextUserID(parent context.Context, userID uint64) context.Context {
	return context.WithValue(parent, userIDKey{}, userID)
}

func GetUserID(ctx context.Context) (uint64, error) {
	id, ok := ctx.Value(userIDKey{}).(uint64)
	if !ok {
		return 0, errors.New("not set user_id")
	}
	return id, nil

}
