package middleware

import (
	"log"
	"net/http"
	"strings"
	"time"

	"go_learning/internal/web"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type LoginJWTMiddlewareBuilder struct {
}

func (m *LoginJWTMiddlewareBuilder) CheckLogin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		path := ctx.Request.URL.Path
		if path == "/users/signup" || path == "/users/login" || path == "/hello" {
			return
		}

		authCode := ctx.GetHeader("Authorization")
		if authCode == "" {
			ctx.AbortWithStatus(http.StatusUnauthorized)
		}
		segs := strings.Split(authCode, " ")
		if len(segs) != 2 {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		tokenStr := segs[1]
		var uc web.UserClaims
		token, err := jwt.ParseWithClaims(tokenStr, &uc, func(token *jwt.Token) (interface{}, error) {
			return web.JWTKey, nil
		})
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if token == nil || !token.Valid {
			ctx.AbortWithStatus(http.StatusUnauthorized)
		}
		if uc.UserAgent != ctx.GetHeader("User-Agent") {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		expireTime := uc.ExpiresAt
		if expireTime.Sub(time.Now()) < time.Minute*15 {
			uc.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Hour))
			tokenStr, err = token.SignedString(web.JWTKey)
			ctx.Header("x-jwt-token", tokenStr)
			if err != nil {
				log.Println(err)
			}
		}
		ctx.Set("user", uc)
	}
}
