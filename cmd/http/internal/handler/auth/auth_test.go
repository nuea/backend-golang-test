package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nuea/backend-golang-test/internal/service"
	"github.com/nuea/backend-golang-test/internal/service/auth"
	userv1 "github.com/nuea/backend-golang-test/proto/gen/backend_golang_test/user/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockAuthService struct {
	mock.Mock
	auth.AuthService
}

func (m *mockAuthService) Login(ctx context.Context, req *userv1.LoginRequest) (accessToken string, err error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return "", args.Error(1)
	}
	return args.Get(0).(string), args.Error(1)
}

func TestProvideHandler(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		sv := new(mockAuthService)
		h := ProvideAuthHandler(&service.Service{AuthService: sv})

		assert.NotNil(t, h)
	})
}

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

	req, err := http.NewRequestWithContext(ctx, method, path, body)
	req.Header.Set("Content-Type", "application/json")

	assert.NoError(t, err)
	ctx.Request = req

	return rec, ctx
}

func TestLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	path := "/admin/v1/login"

	t.Run("success", func(t *testing.T) {
		sv := new(mockAuthService)
		h := &Handler{authsv: sv}
		req := &LoginRequest{
			Email:    "test@example.com",
			Password: "password",
		}
		greq := &userv1.LoginRequest{
			Email:    req.Email,
			Password: req.Password,
		}

		rec, ctx := setupTestRequest(t, http.MethodPost, path, req)
		ac := "fake-access-token"
		sv.On("Login", ctx, greq).Return(ac, nil).Once()

		h.Login(ctx)

		assert.Equal(t, http.StatusOK, rec.Code)
		var response LoginResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, ac, response.AccessToken)

		sv.AssertExpectations(t)
	})

	t.Run("bad request - invalid json", func(t *testing.T) {
		sv := new(mockAuthService)
		h := &Handler{authsv: sv}
		req := "invalid json"

		rec, ctx := setupTestRequest(t, http.MethodPost, path, req)
		h.Login(ctx)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "invalid character")
	})

	t.Run("bad request - validation failed", func(t *testing.T) {
		sv := new(mockAuthService)
		h := &Handler{authsv: sv}
		req := &LoginRequest{
			Email: "test@example.com",
		}

		rec, ctx := setupTestRequest(t, http.MethodPost, path, req)
		h.Login(ctx)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "password is required")
	})

	t.Run("internal server error", func(t *testing.T) {
		sv := new(mockAuthService)
		h := &Handler{authsv: sv}
		req := &LoginRequest{
			Email:    "test@example.com",
			Password: "password",
		}
		greq := &userv1.LoginRequest{
			Email:    req.Email,
			Password: req.Password,
		}

		rec, ctx := setupTestRequest(t, http.MethodPost, path, req)

		msgerr := errors.New("internal server error")
		sv.On("Login", ctx, greq).Return(nil, msgerr).Once()

		h.Login(ctx)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), msgerr.Error())
		sv.AssertExpectations(t)
	})
}
