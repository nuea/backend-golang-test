package user

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/gotidy/ptr"
	"github.com/nuea/backend-golang-test/internal/client"
	"github.com/nuea/backend-golang-test/internal/client/backendgolangtest"
	userv1 "github.com/nuea/backend-golang-test/proto/gen/backend_golang_test/user/v1"
	pbmock "github.com/nuea/backend-golang-test/proto/mock"
	"github.com/stretchr/testify/assert"
)

func setupTestRequest(t *testing.T, method, path string, payload interface{}) (*httptest.ResponseRecorder, *gin.Context) {
	rec := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rec)

	var body *bytes.Buffer
	if str, ok := payload.(string); ok {
		body = bytes.NewBufferString(str)
	} else if payload != nil {
		jb, err := json.Marshal(payload)
		assert.NoError(t, err)
		body = bytes.NewBuffer(jb)
	} else {
		body = bytes.NewBuffer(nil)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, path, body)
	req.Header.Set("Content-Type", "application/json")

	assert.NoError(t, err)
	ctx.Request = req

	return rec, ctx
}

func TestProvideUserHandler(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		musc := pbmock.NewMockUserServiceClient(ctrl)

		h := ProvideUserHandler(&client.GRPCClients{
			BackendGolangTestGRPCService: &backendgolangtest.BackendGolangTestGRPCService{
				UserServiceClient: musc,
			},
		})
		assert.NotNil(t, h)
	})
}

func TestCreateUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	musc := pbmock.NewMockUserServiceClient(ctrl)
	h := &Handler{
		begotc: musc,
	}
	path := "/admin/v1/users"

	t.Run("success", func(t *testing.T) {
		req := &CreateRequest{
			Name:     "test",
			Email:    "test@example.com",
			Password: "password",
		}
		greq := &userv1.CreateUserRequest{
			Name:     req.Name,
			Email:    req.Email,
			Password: req.Password,
		}

		rec, ctx := setupTestRequest(t, http.MethodPost, path, req)
		musc.EXPECT().CreateUser(ctx, greq).Return(nil, nil).Times(1)
		h.CreateUser(ctx)

		assert.Equal(t, http.StatusOK, rec.Code)
		var res CreateResponse
		err := json.Unmarshal(rec.Body.Bytes(), &res)
		assert.NoError(t, err)
		assert.Equal(t, "Completed successfully", res.Message)
	})

	t.Run("bad request - invalid json", func(t *testing.T) {
		req := "invalid json"
		rec, ctx := setupTestRequest(t, http.MethodPost, path, req)

		h.CreateUser(ctx)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "invalid character")
	})

	t.Run("bad request - validation failed", func(t *testing.T) {
		req := &CreateRequest{
			Email:    "test@example.com",
			Password: "password",
		}

		rec, ctx := setupTestRequest(t, http.MethodPost, path, req)
		h.CreateUser(ctx)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "name is required")
	})

	t.Run("internal server error", func(t *testing.T) {
		req := &CreateRequest{
			Name:     "test",
			Email:    "test@example.com",
			Password: "password",
		}
		greq := &userv1.CreateUserRequest{
			Name:     req.Name,
			Email:    req.Email,
			Password: req.Password,
		}

		rec, ctx := setupTestRequest(t, http.MethodPost, path, req)
		msgerr := errors.New("internal server error")
		musc.EXPECT().CreateUser(ctx, greq).Return(nil, msgerr).Times(1)
		h.CreateUser(ctx)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), msgerr.Error())
	})
}

func TestGetUsers(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	musc := pbmock.NewMockUserServiceClient(ctrl)
	h := &Handler{
		begotc: musc,
	}
	path := "/admin/v1/users"

	t.Run("success - without query params", func(t *testing.T) {
		req := &GetUsersRequest{
			Name:  nil,
			Email: nil,
		}
		gReq := &userv1.GetUsersRequest{
			Name:  req.Name,
			Email: req.Email,
		}
		data := []*userv1.User{
			{Id: "686b6ce8dbf72bfc4d0fef95", Name: "test", Email: "test@example.com"},
			{Id: "686b6ce8dbf72bfc4d0fef96", Name: "test", Email: "testtest@example.com"},
		}
		gRes := &userv1.GetUsersResponse{
			Data: data,
		}

		rec, ctx := setupTestRequest(t, http.MethodGet, path, req)
		musc.EXPECT().GetUsers(ctx, gReq).Return(gRes, nil).Times(1)
		h.GetUsers(ctx)

		assert.Equal(t, http.StatusOK, rec.Code)
		var res GetUsersResponse
		err := json.Unmarshal(rec.Body.Bytes(), &res)
		assert.NoError(t, err)
		assert.Equal(t, data[0].Id, res.Data[0].ID)
		assert.Equal(t, data[1].Id, res.Data[1].ID)
	})

	t.Run("success - with query params", func(t *testing.T) {
		req := &GetUsersRequest{
			Name:  nil,
			Email: ptr.String("test@example.com"),
		}
		gReq := &userv1.GetUsersRequest{
			Name:  req.Name,
			Email: req.Email,
		}
		data := []*userv1.User{
			{Id: "686b6ce8dbf72bfc4d0fef95", Name: "test", Email: "test@example.com"},
		}
		gRes := &userv1.GetUsersResponse{
			Data: data,
		}

		rec, ctx := setupTestRequest(t, http.MethodGet, path, req)
		musc.EXPECT().GetUsers(ctx, gReq).Return(gRes, nil).Times(1)
		h.GetUsers(ctx)

		assert.Equal(t, http.StatusOK, rec.Code)
		var res GetUsersResponse
		err := json.Unmarshal(rec.Body.Bytes(), &res)
		assert.NoError(t, err)
		assert.Equal(t, data[0].Id, res.Data[0].ID)
	})

	t.Run("bad request - invalid json", func(t *testing.T) {
		req := "invalid json"
		rec, ctx := setupTestRequest(t, http.MethodGet, path, req)

		h.GetUsers(ctx)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "invalid character")
	})

	t.Run("internal server error - gRPC", func(t *testing.T) {
		req := &GetUsersRequest{
			Name:  nil,
			Email: ptr.String("test@example.com"),
		}
		gReq := &userv1.GetUsersRequest{
			Name:  req.Name,
			Email: req.Email,
		}

		rec, ctx := setupTestRequest(t, http.MethodGet, path, req)
		msgerr := errors.New("internal server error")
		musc.EXPECT().GetUsers(ctx, gReq).Return(nil, msgerr).Times(1)
		h.GetUsers(ctx)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), msgerr.Error())
	})
}

func TestGetUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	musc := pbmock.NewMockUserServiceClient(ctrl)
	h := &Handler{
		begotc: musc,
	}

	path := "/admin/v1/users"
	uid := "686b6ce8dbf72bfc4d0fef95"

	t.Run("success", func(t *testing.T) {
		gReq := &userv1.GetUserRequest{
			Id: uid,
		}
		gRes := &userv1.GetUserResponse{
			User: &userv1.User{
				Id:    uid,
				Name:  "test",
				Email: "test@example.com",
			},
		}

		rec, ctx := setupTestRequest(t, http.MethodGet, path, nil)
		if uid != "" {
			ctx.Params = append(ctx.Params, gin.Param{Key: "id", Value: uid})
		}
		musc.EXPECT().GetUser(ctx, gReq).Return(gRes, nil).Times(1)
		h.GetUser(ctx)

		assert.Equal(t, http.StatusOK, rec.Code)
		var res GetUserResponse
		err := json.Unmarshal(rec.Body.Bytes(), &res)
		assert.NoError(t, err)
		assert.Equal(t, uid, res.ID)
	})

	t.Run("bad request - path parameter is missing", func(t *testing.T) {
		rec, ctx := setupTestRequest(t, http.MethodGet, path, nil)
		h.GetUser(ctx)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "path parameter is missing.")
	})

	t.Run("internal server error - gRPC", func(t *testing.T) {
		rec, ctx := setupTestRequest(t, http.MethodGet, path, nil)
		if uid != "" {
			ctx.Params = append(ctx.Params, gin.Param{Key: "id", Value: uid})
		}
		msgerr := errors.New("internal server error")
		musc.EXPECT().GetUser(ctx, gomock.Any()).Return(nil, msgerr).Times(1)
		h.GetUser(ctx)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), msgerr.Error())
	})
}

func TestUpdateUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	musc := pbmock.NewMockUserServiceClient(ctrl)
	h := &Handler{
		begotc: musc,
	}

	path := "/admin/v1/users"
	uid := "686b6ce8dbf72bfc4d0fef95"

	t.Run("success", func(t *testing.T) {
		req := &UpdateUserRequest{
			Name:  ptr.String("test"),
			Email: ptr.String("test@example.com"),
		}
		gReq := &userv1.UpdateUserRequest{
			Id:    uid,
			Name:  req.Name,
			Email: req.Email,
		}

		rec, ctx := setupTestRequest(t, http.MethodPatch, path, req)
		if uid != "" {
			ctx.Params = append(ctx.Params, gin.Param{Key: "id", Value: uid})
		}
		musc.EXPECT().UpdateUser(ctx, gReq).Return(nil, nil).Times(1)
		h.UpdateUser(ctx)

		assert.Equal(t, http.StatusOK, rec.Code)
		var res UpdateUserResponse
		err := json.Unmarshal(rec.Body.Bytes(), &res)
		assert.NoError(t, err)
		assert.Equal(t, "Updated successfully", res.Message)
	})

	t.Run("bad request - path parameter is missing", func(t *testing.T) {
		rec, ctx := setupTestRequest(t, http.MethodPatch, path, nil)
		h.UpdateUser(ctx)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "path parameter is missing.")
	})

	t.Run("bad request - invalid json", func(t *testing.T) {
		req := "invalid json"
		rec, ctx := setupTestRequest(t, http.MethodPatch, path, req)
		if uid != "" {
			ctx.Params = append(ctx.Params, gin.Param{Key: "id", Value: uid})
		}

		h.UpdateUser(ctx)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "invalid character")
	})

	t.Run("internal server error", func(t *testing.T) {
		req := &UpdateUserRequest{
			Name:  ptr.String("test"),
			Email: ptr.String("test@example.com"),
		}
		rec, ctx := setupTestRequest(t, http.MethodGet, path, req)
		if uid != "" {
			ctx.Params = append(ctx.Params, gin.Param{Key: "id", Value: uid})
		}
		msgerr := errors.New("internal server error")
		musc.EXPECT().UpdateUser(ctx, gomock.Any()).Return(nil, msgerr).Times(1)
		h.UpdateUser(ctx)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), msgerr.Error())
	})
}

func TestDeleteUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	musc := pbmock.NewMockUserServiceClient(ctrl)
	h := &Handler{
		begotc: musc,
	}

	path := "/admin/v1/users"
	uid := "686b6ce8dbf72bfc4d0fef95"

	t.Run("success", func(t *testing.T) {
		gReq := &userv1.DeleteUserRequest{
			Id: uid,
		}

		rec, ctx := setupTestRequest(t, http.MethodDelete, path, nil)
		if uid != "" {
			ctx.Params = append(ctx.Params, gin.Param{Key: "id", Value: uid})
		}
		musc.EXPECT().DeleteUser(ctx, gReq).Return(nil, nil).Times(1)
		h.DeleteUser(ctx)

		assert.Equal(t, http.StatusOK, rec.Code)
		var res UpdateUserResponse
		err := json.Unmarshal(rec.Body.Bytes(), &res)
		assert.NoError(t, err)
		assert.Equal(t, "Deleted successfully", res.Message)
	})

	t.Run("bad request - path parameter is missing", func(t *testing.T) {
		rec, ctx := setupTestRequest(t, http.MethodDelete, path, nil)
		h.DeleteUser(ctx)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "path parameter is missing.")
	})

	t.Run("internal server error", func(t *testing.T) {
		rec, ctx := setupTestRequest(t, http.MethodDelete, path, nil)
		if uid != "" {
			ctx.Params = append(ctx.Params, gin.Param{Key: "id", Value: uid})
		}
		msgerr := errors.New("internal server error")
		musc.EXPECT().DeleteUser(ctx, gomock.Any()).Return(nil, msgerr).Times(1)
		h.DeleteUser(ctx)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), msgerr.Error())
	})
}
