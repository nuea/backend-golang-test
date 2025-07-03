package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nuea/backend-golang-test/internal/service"
	"github.com/nuea/backend-golang-test/internal/service/auth"
	"github.com/nuea/backend-golang-test/internal/util"
	userv1 "github.com/nuea/backend-golang-test/proto/gen/backend_golang_test/user/v1"
)

type Handler struct {
	authsv auth.AuthService
}

func ProvideUserHandler(sv *service.Service) *Handler {
	return &Handler{
		authsv: sv.AuthService,
	}
}

func (h *Handler) Login(ctx *gin.Context) {
	var req *LoginRequest
	if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := util.ValidateStruct(req); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	gReq := &userv1.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	ac, err := h.authsv.Login(ctx, gReq)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, &LoginResponse{
		AccessToken: ac,
	})
}
