package auth

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nuea/backend-golang-test/internal/service"
	"github.com/nuea/backend-golang-test/internal/service/auth"
)

type AuthMiddleware interface {
	Middleware() gin.HandlerFunc
}

type authMiddleware struct {
	authsv auth.AuthService
}

func ProvideAuthMiddleware(sv *service.Service) AuthMiddleware {
	return &authMiddleware{
		authsv: sv.AuthService,
	}
}

func (m *authMiddleware) Middleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if err := m.authentication(ctx); err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			return
		}
		ctx.Next()
	}
}

func (m *authMiddleware) authentication(ctx *gin.Context) error {
	var ac string
	if acs := strings.Fields(ctx.GetHeader("Authorization")); len(acs) > 1 && acs[0] == "Bearer" {
		ac = acs[1]
	} else {
		return errors.New("Unauthorized.")
	}

	claims, err := m.authsv.VerifyAccessToken(ac)
	if err != nil {
		return err
	}

	if claims.ExpiresAt < time.Now().Local().UnixMilli() {
		return errors.New("Unauthorized.")
	}
	return nil
}
