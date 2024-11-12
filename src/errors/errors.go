package errors

import (
	i18n_err "github.com/usagifm/dating-app/lib/i18n/errors"
)

var (
	ErrQuotaExceeded     = i18n_err.NewI18nError("err_quota_exceeded")
	ErrOtpUsed           = i18n_err.NewI18nError("err_otp_used")
	ErrInvalidOtp        = i18n_err.NewI18nError("err_invalid_otp")
	ErrExpiredOtp        = i18n_err.NewI18nError("err_expired_otp")
	ErrSessionIdNotFound = i18n_err.NewI18nError("err_session_id_not_found")
)
