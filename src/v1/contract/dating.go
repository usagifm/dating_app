package contract

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/usagifm/dating-app/lib/logger"
	"github.com/usagifm/dating-app/src/app"
)

type UpdateUserPreferenceRequest struct {
	MinAge          int    `json:"min_age"`
	MaxAge          int    `json:"max_age"`
	PreferredGender string `json:"preferred_gender"`
}

func ValidateAndBuildUpdateUserPreferenceRequest(r *http.Request) (UpdateUserPreferenceRequest, error) {
	_, span := app.Tracer().Start(r.Context(), "ValidateAndBuildUpdateUserPreferenceRequest")
	defer span.End()

	var payload UpdateUserPreferenceRequest
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
	if payload.PreferredGender != "male" && payload.PreferredGender != "female" && payload.PreferredGender != "both" {
		logger.GetLogger(r.Context()).Error("validate request body err: ", "gender can only be male or female or both")
		return payload, err
	}

	return payload, nil
}

type SwipeRequest struct {
	SwipedId  int    `json:"swiped_id"`
	SwipeType string `json:"swipe_type"`
}

func ValidateAndBuildUpdateSwipeRequest(r *http.Request) (SwipeRequest, error) {
	_, span := app.Tracer().Start(r.Context(), "ValidateAndBuildUpdateUserPreferenceRequest")
	defer span.End()

	var payload SwipeRequest
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
	if payload.SwipeType != "pass" && payload.SwipeType != "like" {
		logger.GetLogger(r.Context()).Error("validate request body err: ", "swipe type can only be pass or like")
		return payload, err
	}

	return payload, nil
}

type BuyPackageRequest struct {
	PackageId int `json:"package_id" validate:"required"`
}

func ValidateAndBuildBuyPackageRequest(r *http.Request) (BuyPackageRequest, error) {
	_, span := app.Tracer().Start(r.Context(), "ValidateAndBuildBuyPackageRequest")
	defer span.End()

	var payload BuyPackageRequest
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
