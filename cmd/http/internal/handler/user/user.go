package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nuea/backend-golang-test/internal/client"
	"github.com/nuea/backend-golang-test/internal/util"
	userv1 "github.com/nuea/backend-golang-test/proto/gen/backend_golang_test/user/v1"
)

type Handler struct {
	begotc userv1.UserServiceClient
}

func ProvideUserHandler(c *client.GRPCClients) *Handler {
	return &Handler{
		begotc: c.BackendGolangTestGRPCService.UserServiceClient,
	}
}

// @id CreateUser
// @accept  json
// @produce  json
// @tags User
// @param req body CreateRequest true "req"
// @success 200 {object} CreateResponse
// @router /api/v1/users [POST]
func (h *Handler) CreateUser(ctx *gin.Context) {
	var req *CreateRequest
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

	if _, err := h.begotc.CreateUser(ctx, &userv1.CreateUserRequest{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, &CreateResponse{
		Message: "Register success",
	})
}

// @id GetUsers
// @accept  json
// @produce  json
// @security BearerAuth
// @tags User
// @param req formData GetUsersRequest true "req"
// @success 200 {object} GetUsersResponse
// @router /admin/v1/users [GET]
func (h *Handler) GetUsers(ctx *gin.Context) {
	var req GetUsersRequest

	if err := ctx.ShouldBind(&req); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	users, err := h.begotc.GetUsers(ctx, &userv1.GetUsersRequest{
		Name:  req.Name,
		Email: req.Email,
	})
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	datas, err := util.MapToSlice(mapToUser, users.Data)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, &GetUsersResponse{
		Data: datas,
	})
}

// @id GetUser
// @accept  json
// @produce  json
// @security BearerAuth
// @tags User
// @param id path string true "id"
// @success 200 {object} GetUserResponse
// @router /admin/v1/users/{id} [GET]
func (h *Handler) GetUser(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "path parameter is missing.",
		})
		return
	}

	gRes, err := h.begotc.GetUser(ctx, &userv1.GetUserRequest{
		Id: id,
	})
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	user, err := mapToUser(gRes.User)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, &GetUserResponse{
		User: *user,
	})
}

// @id UpdateUser
// @accept  json
// @produce  json
// @security BearerAuth
// @tags User
// @param req body UpdateUserRequest true "req"
// @param id path string true "id"
// @success 200 {object} UpdateUserResponse
// @router /admin/v1/users/{id} [PATCH]
func (h *Handler) UpdateUser(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "path parameter is missing.",
		})
		return
	}

	var req *UpdateUserRequest
	if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if _, err := h.begotc.UpdateUser(ctx, &userv1.UpdateUserRequest{
		Id:    id,
		Name:  req.Name,
		Email: req.Email,
	}); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, &UpdateUserResponse{
		Message: "Update success",
	})
}

// @id DeleteUser
// @accept  json
// @produce  json
// @security BearerAuth
// @tags User
// @param id path string true "id"
// @success 200 {object} DeleteUserResponse
// @router /admin/v1/users/{id} [DELETE]
func (h *Handler) DeleteUser(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "path parameter is missing.",
		})
		return
	}

	if _, err := h.begotc.DeleteUser(ctx, &userv1.DeleteUserRequest{
		Id: id,
	}); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, &DeleteUserResponse{
		Message: "Delete success",
	})
}
