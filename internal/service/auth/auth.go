package auth

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/nuea/backend-golang-test/internal/client"
	"github.com/nuea/backend-golang-test/internal/config"
	userv1 "github.com/nuea/backend-golang-test/proto/gen/backend_golang_test/user/v1"
)

type AuthService interface {
	Login(ctx context.Context, req *userv1.LoginRequest) (accessToken string, err error)
	GenerateAccessToken(userID string) (string, error)
	VerifyAccessToken(accessToken string) (*JwtToken, error)
}

type authService struct {
	cfg        *config.AuthConfig
	authclient userv1.AuthServiceClient
}

func ProvideAuthenticationService(cfg *config.AppConfig, c *client.Clients) AuthService {
	return &authService{
		cfg:        &cfg.Auth,
		authclient: c.BackendGolangTestGRPCService.AuthServiceClient,
	}
}

type JwtToken struct {
	jwt.StandardClaims
	UserID string `json:"uid,omitempty"`
}

func (s *authService) Login(ctx context.Context, req *userv1.LoginRequest) (accessToken string, err error) {
	res, err := s.authclient.Login(ctx, req)
	if err != nil {
		return "", err
	}

	accessToken, err = s.GenerateAccessToken(res.UserId)
	if err != nil {
		return "", err
	}

	if accessToken != "" {
		s.setCookies(ctx, accessToken)
	}

	return accessToken, nil
}

func (s *authService) GenerateAccessToken(userID string) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, &JwtToken{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(s.cfg.AccessTokenExpireTTL).UnixMilli(),
		},
	}).SignedString([]byte(s.cfg.SecretKey))
}

func (s *authService) VerifyAccessToken(accessToken string) (*JwtToken, error) {
	if accessToken == "" {
		return nil, errors.New("Access token is empty.")
	}

	token, err := jwt.ParseWithClaims(accessToken, &JwtToken{},
		func(t *jwt.Token) (interface{}, error) {
			return []byte(s.cfg.SecretKey), nil
		},
	)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JwtToken)
	if !ok {
		return nil, errors.New("Access token is invalid.")
	}
	return claims, nil
}

func (s *authService) setCookies(ctx context.Context, accessToken string) {
	gCtx := ctx.(*gin.Context)
	gCtx.SetSameSite(http.SameSiteDefaultMode)
	gCtx.SetCookie("_uac", accessToken, int(s.cfg.AccessTokenExpireTTL/time.Second), "/", "localhost", false, true)
	gCtx.SetCookie("_uac_s", accessToken, int(s.cfg.AccessTokenExpireTTL/time.Second), "/", "localhost", true, true)

}
