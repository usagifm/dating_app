package response

import (
	"context"
	"net/http"

	i18n_err "github.com/usagifm/dating-app/lib/i18n/errors"

	"github.com/usagifm/dating-app/lib/middleware/request"
)

func JSONSuccessResponse(ctx context.Context, w http.ResponseWriter, data interface{}, message string) {
	JSONResponse(ctx, w, createSuccessResponse(data, request.GetRequestID(ctx), http.StatusOK, message), http.StatusOK)
}

func JSONCreatedResponse(ctx context.Context, w http.ResponseWriter, data interface{}, message string) {
	JSONResponse(ctx, w, createSuccessResponse(data, request.GetRequestID(ctx), http.StatusCreated, message), http.StatusCreated)
}

func JSONUnauthorizedResponse(ctx context.Context, w http.ResponseWriter, message string) {
	JSONResponse(ctx, w, createErrorResponse(http.StatusUnauthorized, message, i18n_err.ErrUnauthorized, request.GetRequestID(ctx), request.GetLanguage(ctx), nil), http.StatusUnauthorized)
}

func JSONInternalErrorResponse(ctx context.Context, w http.ResponseWriter, message string) {
	JSONResponse(ctx, w, createErrorResponse(http.StatusInternalServerError, message, i18n_err.ErrInternalServer, request.GetRequestID(ctx), request.GetLanguage(ctx), nil), http.StatusInternalServerError)
}

func JSONBadRequestResponse(ctx context.Context, w http.ResponseWriter, message string) {
	JSONResponse(ctx, w, createErrorResponse(http.StatusBadRequest, message, i18n_err.ErrBadRequest, request.GetRequestID(ctx), request.GetLanguage(ctx), nil),
		http.StatusBadRequest)
}

func JSONUnprocessableEntity(ctx context.Context, w http.ResponseWriter, err i18n_err.I18nError, action *Action, message string) {
	JSONResponse(ctx, w, createErrorResponse(http.StatusUnprocessableEntity, message, err, request.GetRequestID(ctx), request.GetLanguage(ctx), action), http.StatusUnprocessableEntity)
}
