package auth

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/mitchellh/mapstructure"
	"github.com/usagifm/dating-app/lib/logger"
	"github.com/usagifm/dating-app/src/app"
	"github.com/usagifm/dating-app/src/entity"
	"github.com/usagifm/dating-app/src/middleware/response"
)

type (
	ctxKeyAuthUser struct{}
)

const (
	Authorization = "Authorization"
)

type AuthenticatedRequestValidator func(r *http.Request) *entity.User

func DefaultAuthenticatedRequestValidator(r *http.Request) *entity.User {
	authToken := r.Header.Get(Authorization)

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
		logger.GetLogger(r.Context()).Errorf("invalid auth token :%s", r.Header.Get(Authorization), err)
		return nil
	}
	claims, ok := tokenByte.Claims.(jwt.MapClaims)
	if !ok || !tokenByte.Valid {
		logger.GetLogger(r.Context()).Errorf("invalid auth token :%s", r.Header.Get(Authorization), err)
		return nil
	}

	user := entity.User{}

	if err := mapstructure.Decode(claims["user"], &user); err != nil {
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
				response.JSONUnauthorizedResponse(ctx, w)
				return
			}

			ctxValue := context.WithValue(ctx, CtxKeyAuthUser, *user)
			ctx = logger.WithLogger(ctxValue, logger.GetLogger(ctxValue).WithFields(logger.Fields{
				"user_id": user.Id,
			}))

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
