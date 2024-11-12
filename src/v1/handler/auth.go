package handler

import (
	"net/http"

	"github.com/usagifm/dating-app/lib/logger"
	"github.com/usagifm/dating-app/lib/response"
	"github.com/usagifm/dating-app/src/app"
	"github.com/usagifm/dating-app/src/v1/contract"
)

// SignUp godoc
//
//	@Summary		Register a new user
//	@Description	Sign up a new user with the provided details.
//	@Tags			Authentication
//	@Accept			json
//	@Produce		json
//	@Param			body	body	contract.SignUpRequest	true	"Sign Up Request Body"
//	@Router			/api/v1/auth/signup [post]
func SignUp(h *DatingAppHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := app.Tracer().Start(r.Context(), "SignUpHandler")
		defer span.End()

		body, err := contract.ValidateAndBuildSignUpRequest(r.WithContext(ctx))
		if err != nil {
			logger.GetLogger(r.Context()).Error("ValidateAndBuildSignUpRequest err:", err)
			response.JSONBadRequestResponse(ctx, w, "Validation error")
			return
		}

		errSignUp := h.sAuth.SignUp(ctx, body)
		if errSignUp != nil {
			logger.GetLogger(r.Context()).Error("SignUp err:", errSignUp)
			response.JSONInternalErrorResponse(ctx, w, errSignUp.Error())
			return
		}

		response.JSONSuccessResponse(ctx, w, nil, "Sign Up success, please proceed to login")

	}
}

// Login godoc
//
//	@Summary		User login
//	@Description	Log in an existing user and get a token.
//	@Tags			Authentication
//	@Accept			json
//	@Produce		json
//	@Param			body	body	contract.LoginRequest true "Login Request Body"
//	@Router			/api/v1/auth/login [post]
func Login(h *DatingAppHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx, span := app.Tracer().Start(r.Context(), "LoginHandler")
		defer span.End()

		body, err := contract.ValidateAndBuildLoginRequest(r.WithContext(ctx))
		if err != nil {
			logger.GetLogger(r.Context()).Error("ValidateAndBuildLoginRequest err:", err)
			response.JSONBadRequestResponse(ctx, w, "Validation error")
			return
		}

		token, err := h.sAuth.Login(ctx, body)
		if err != nil {
			logger.GetLogger(r.Context()).Error("Login err:", err)
			response.JSONInternalErrorResponse(ctx, w, err.Error())
			return
		}

		response.JSONSuccessResponse(ctx, w, token, "Login success")

	}
}

// GetProfile godoc
//
//	@Summary		Get user profile
//	@Description	Retrieve the profile of the currently authenticated user.
//	@Tags			User
//	@Produce		json
//	@Router			/api/v1/auth/profile [get]
func GetProfile(h *DatingAppHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := app.Tracer().Start(r.Context(), "GetProfileHandler")
		defer span.End()

		user, err := h.sAuth.GetProfile(ctx)
		if err != nil {
			logger.GetLogger(r.Context()).Error("Login err:", err)
			response.JSONInternalErrorResponse(ctx, w, err.Error())
			return
		}

		response.JSONSuccessResponse(ctx, w, user, "Get profile success")

	}
}
