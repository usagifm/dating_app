package auth

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/mitchellh/mapstructure"
	"github.com/usagifm/dating-app/lib/logger"
	"github.com/usagifm/dating-app/lib/response"
	"github.com/usagifm/dating-app/src/app"
	"github.com/usagifm/dating-app/src/entity"
)

type (
	ctxKeyAuthUser struct{}
)

const (
	AUTHORIZATION = "x-user-id"
)

type AuthenticatedRequestValidator func(r *http.Request) *entity.User

func DefaultAuthenticatedRequestValidator(r *http.Request) *entity.User {
	authToken := r.Header.Get(AUTHORIZATION)
	if authToken == "" {
		return nil
	}

	if strings.HasPrefix(authToken, "Bearer ") {
		authToken = strings.TrimPrefix(authToken, "Bearer ")
	}

	tokenByte, err := jwt.Parse(authToken, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %s", jwtToken.Header["alg"])
		}

		return []byte(app.Config().JWTSecret), nil
	})
	if err != nil {
		logger.GetLogger(r.Context()).Errorf("invalid auth token :%s", r.Header.Get(AUTHORIZATION), err)
		return nil
	}

	claims, ok := tokenByte.Claims.(jwt.MapClaims)
	if !ok || !tokenByte.Valid {
		logger.GetLogger(r.Context()).Errorf("invalid auth token :%s", r.Header.Get(AUTHORIZATION), err)
		return nil
	}

	user := entity.User{}
	if err := mapstructure.Decode(claims, &user); err != nil {
		logger.GetLogger(r.Context()).Errorf("error decoding claims: %s", err)
		return nil
	}

	return &user
}

var CtxKeyAuthUser ctxKeyAuthUser = ctxKeyAuthUser{}

func New(validator AuthenticatedRequestValidator) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			user := validator(r)
			if user == nil {
				response.JSONUnauthorizedResponse(ctx, w, "Unauthorized")
				return
			}

			ctx = context.WithValue(ctx, CtxKeyAuthUser, *user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUser(ctx context.Context) *entity.User {
	if user, ok := ctx.Value(CtxKeyAuthUser).(entity.User); ok {
		return &user
	}
	return nil
}
