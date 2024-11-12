package handler

import (
	"net/http"

	"github.com/usagifm/dating-app/lib/logger"
	"github.com/usagifm/dating-app/lib/response"
	"github.com/usagifm/dating-app/src/app"
	"github.com/usagifm/dating-app/src/v1/contract"
)

// GetUserPreference godoc
//
//	@Summary		Get user preference
//	@Description	Retrieve the user's preference settings
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Router			/api/v1/dating/preference [get]
func GetUserPreference(h *DatingAppHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := app.Tracer().Start(r.Context(), "GetUserPreferenceHandler")
		defer span.End()

		userPreference, err := h.sDating.GetUserPreference(ctx)
		if err != nil {
			logger.GetLogger(r.Context()).Error("GetUserPreference err:", err)
			response.JSONInternalErrorResponse(ctx, w, err.Error())
			return
		}

		response.JSONSuccessResponse(ctx, w, userPreference, "Success get user preference")

	}
}

// UpdateUserPreference godoc
//
//	@Summary		Update user preference
//	@Description	Update the user's preference settings
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			body	body	contract.UpdateUserPreferenceRequest	true	"Update user preference request body"
//	@Router			/api/v1/dating/preference [put]
func UpdateUserPreference(h *DatingAppHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := app.Tracer().Start(r.Context(), "UpdateUserPreferenceHandler")
		defer span.End()

		body, err := contract.ValidateAndBuildUpdateUserPreferenceRequest(r.WithContext(ctx))
		if err != nil {
			logger.GetLogger(r.Context()).Error("ValidateAndBuildUpdateUserPreferenceRequest err:", err)
			response.JSONBadRequestResponse(ctx, w, "Validation error")
			return
		}

		err = h.sDating.UpdateUserPreference(ctx, body)
		if err != nil {
			logger.GetLogger(r.Context()).Error("UpdateUserPreference err:", err)
			response.JSONInternalErrorResponse(ctx, w, err.Error())
			return
		}

		response.JSONSuccessResponse(ctx, w, nil, "Update preference data success")

	}
}

// GetProfilesByPreference godoc
//
//	@Summary		Get profiles by user preference
//	@Description	Retrieve a list of profiles based on user preference
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Router			/api/v1/dating [get]
func GetProfilesByPreference(h *DatingAppHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := app.Tracer().Start(r.Context(), "GetProfilesByPreferenceHandler")
		defer span.End()

		userProfiles, err := h.sDating.GetProfilesByPreference(ctx)
		if err != nil {
			logger.GetLogger(r.Context()).Error("GetUserPreference err:", err)
			response.JSONInternalErrorResponse(ctx, w, err.Error())
			return
		}

		response.JSONSuccessResponse(ctx, w, userProfiles, "Success get user profiles")

	}
}

// GetUserMatches godoc
//
//	@Summary		Get matched users
//	@Description	Retrieve a list of matched users
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Router			/api/v1/dating/matches [get]
func GetUserMatches(h *DatingAppHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := app.Tracer().Start(r.Context(), "GetUserMatchesHandler")
		defer span.End()

		userProfiles, err := h.sDating.GetUserMatches(ctx)
		if err != nil {
			logger.GetLogger(r.Context()).Error("GetUserMatches err:", err)
			response.JSONInternalErrorResponse(ctx, w, err.Error())
			return
		}

		response.JSONSuccessResponse(ctx, w, userProfiles, "Success get matched user")
	}
}

// Swipe godoc
//
//	@Summary		Swipe user
//	@Description	Swipe a user and check if it's a match
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			body	body	contract.SwipeRequest	true	"Swipe request body"
//	@Router			/api/v1/dating/swipe [post]
func Swipe(h *DatingAppHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := app.Tracer().Start(r.Context(), "SwipeHandler")
		defer span.End()

		body, err := contract.ValidateAndBuildUpdateSwipeRequest(r.WithContext(ctx))
		if err != nil {
			logger.GetLogger(r.Context()).Error("ValidateAndBuildUpdateSwipeRequest err:", err)
			response.JSONBadRequestResponse(ctx, w, "Validation error")
			return
		}

		isMatched, err := h.sDating.Swipe(ctx, body)
		if err != nil {
			logger.GetLogger(r.Context()).Error("Swipe err:", err)
			response.JSONInternalErrorResponse(ctx, w, err.Error())
			return
		}

		if isMatched {
			response.JSONSuccessResponse(ctx, w, nil, "congratulation, its a match")
			return
		}

		response.JSONSuccessResponse(ctx, w, nil, "swiped successfully")

	}
}

// GetPackages godoc
//
//	@Summary		Get available packages
//	@Description	Retrieve the list of packages available for purchase
//	@Tags			Package
//	@Accept			json
//	@Produce		json
//	@Router			/api/v1/dating/package/list [get]
func GetPackages(h *DatingAppHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := app.Tracer().Start(r.Context(), "GetPackagesHandler")
		defer span.End()

		packages, err := h.sDating.GetPackages(ctx)
		if err != nil {
			logger.GetLogger(r.Context()).Error("GetPackages err:", err)
			response.JSONInternalErrorResponse(ctx, w, err.Error())
			return
		}

		response.JSONSuccessResponse(ctx, w, packages, "Success get package list")
	}
}

// BuyPackage godoc
//
//	@Summary		Buy a package
//	@Description	Purchase a package
//	@Tags			Package
//	@Accept			json
//	@Produce		json
//	@Param			body	body	contract.BuyPackageRequest	true	"Buy package request body"
//	@Router			/api/v1/dating/package/buy [post]
func BuyPackage(h *DatingAppHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := app.Tracer().Start(r.Context(), "BuyPackageHandler")
		defer span.End()

		body, err := contract.ValidateAndBuildBuyPackageRequest(r.WithContext(ctx))
		if err != nil {
			logger.GetLogger(r.Context()).Error("ValidateAndBuildBuyPackageRequest err:", err)
			response.JSONBadRequestResponse(ctx, w, "Validation error")
			return
		}

		err = h.sDating.BuyPackage(ctx, body)
		if err != nil {
			logger.GetLogger(r.Context()).Error("BuyPackage err:", err)
			response.JSONInternalErrorResponse(ctx, w, err.Error())
			return
		}

		response.JSONSuccessResponse(ctx, w, nil, "Success buy package")

	}
}
