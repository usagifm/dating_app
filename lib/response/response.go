package response

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/usagifm/dating-app/lib/i18n"
	i18n_err "github.com/usagifm/dating-app/lib/i18n/errors"
)

type Response struct {
	Data     interface{} `json:"data"`
	Error    *Error      `json:"error"`
	Success  bool        `json:"success"`
	Metadata Meta        `json:"metadata"`
	Message  string      `json:"message"`
	Status   int         `json:"status"`
}

type Meta struct {
	RequestId string `json:"request_id"`
}

type Error struct {
	Code     string  `json:"code"`
	Title    string  `json:"message_title"`
	Message  string  `json:"message"`
	Severity string  `json:"message_severity"`
	Action   *Action `json:"action"`
}

type Action struct {
	NextState string `json:"next_state"`
}

const (
	NextStateVerification  = "verification"
	NextStateIdpLogin      = "idp_login"
	NextStateResetPassword = "reset_password"
	NextStateLogin         = "login"
)

func JSONResponse(ctx context.Context, w http.ResponseWriter, data Response, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func JSONResponseRaw(ctx context.Context, w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func createSuccessResponse(data interface{}, reqId string, statusCode int, message string) Response {
	return Response{
		Status:  statusCode,
		Message: message,
		Data:    data,
		Success: true,
		Metadata: Meta{
			RequestId: reqId,
		},
	}
}

func createErrorResponse(statusCode int, message string, err i18n_err.I18nError, reqId, lang string, action *Action, args ...interface{}) Response {
	return Response{
		Status:  statusCode,
		Message: message,
		Error: &Error{
			Code:     err.Error(),
			Title:    i18n.Title(lang, err.Error()),
			Message:  i18n.Message(lang, err.Error()),
			Severity: "error",
			Action:   action,
		},
		Metadata: Meta{
			RequestId: reqId,
		},
	}
}
