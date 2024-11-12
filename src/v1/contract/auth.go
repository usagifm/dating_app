package contract

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/usagifm/dating-app/lib/logger"
	"github.com/usagifm/dating-app/src/app"
	"github.com/usagifm/dating-app/src/entity"
)

type CustomClaims struct {
	User       entity.User
	Authorized bool  `json:"authorized"`
	Exp        int64 `json:"exp"`
}

type SignUpRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Name     string `json:"name" validate:"required,min=6"`
	Gender   string `json:"gender" validate:"required"`
	Age      int    `json:"age" validate:"required"`
	Bio      string `json:"bio" validate:"required"`
	PhotoUrl string `json:"photo_url" validate:"required"`

	PreferredGender string `json:"preferred_gender" validate:"required"`
	MinAge          int    `json:"min_age"`
	MaxAge          int    `json:"max_age"`
}

func ValidateAndBuildSignUpRequest(r *http.Request) (SignUpRequest, error) {
	_, span := app.Tracer().Start(r.Context(), "ValidateAndBuildSignUpRequest")
	defer span.End()

	var payload SignUpRequest
	bodyByte, err := io.ReadAll(r.Body)

	if err != nil {
		logger.GetLogger(r.Context()).Error("read request body err: ", err)
		return payload, err
	}

	if err := json.Unmarshal(bodyByte, &payload); err != nil {
		logger.GetLogger(r.Context()).Error("unmarshal request body err: ", err)
		return payload, err
	}

	if err := app.RequestValidator().Struct(payload); err != nil {
		logger.GetLogger(r.Context()).Error("validate request body err: ", err)
		return payload, err
	}

	if payload.Gender != "male" && payload.Gender != "female" {
		logger.GetLogger(r.Context()).Error("validate request body err: ", "gender can only be male or female")
		return payload, err
	}

	if payload.PreferredGender != "male" && payload.PreferredGender != "female" && payload.PreferredGender != "both" {
		logger.GetLogger(r.Context()).Error("validate request body err: ", "gender can only be male or female or both")
		return payload, err
	}

	return payload, nil
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func ValidateAndBuildLoginRequest(r *http.Request) (LoginRequest, error) {
	_, span := app.Tracer().Start(r.Context(), "ValidateAndBuildLoginRequest")
	defer span.End()

	var payload LoginRequest
	bodyByte, err := io.ReadAll(r.Body)

	if err != nil {
		logger.GetLogger(r.Context()).Error("read request body err: ", err)
		return payload, err
	}

	if err := json.Unmarshal(bodyByte, &payload); err != nil {
		logger.GetLogger(r.Context()).Error("unmarshal request body err: ", err)
		return payload, err
	}

	if err := app.RequestValidator().Struct(payload); err != nil {
		logger.GetLogger(r.Context()).Error("validate request body err: ", err)
		return payload, err
	}

	return payload, nil
}
