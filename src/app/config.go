package app

import (
	"context"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/usagifm/dating-app/lib/logger"

	"github.com/spf13/viper"
)

/*
	All config should be required.
	Optional only allowed if zero value of the type is expected being the default value.
	time.Duration units are “ns”, “us” (or “µs”), “ms”, “s”, “m”, “h”. as in time.ParseDuration().
*/

type (
	Postgres struct {
		ConnURI            string        `mapstructure:"PG_CONN_URI" validate:"required"`
		MaxPoolSize        int           `mapstructure:"PG_MAX_POOL_SZE"` //Optional, default to 0 (zero value of int)
		MaxIdleConnections int           `mapstructure:"PG_MAX_IDLE_CONNECTIONS"`
		MaxIdleTime        time.Duration `mapstructure:"PG_MAX_IDLE_TIME"` //Optional, default to '0s' (zero value of time.Duration)
		MaxLifeTime        time.Duration `mapstructure:"PG_MAX_IDLE_TIME"` //Optional, default to '0s' (zero value of time.Duration)
	}

	Redis struct {
		Host           string        `mapstructure:"REDIS_HOST" validate:"required"`
		Password       string        `mapstructure:"REDIS_PASSWORD" validate:"required"`
		DBPrefix       string        `mapstructure:"REDIS_DB_PREFIX" validate:"required"`
		InvalidateTime time.Duration `mapstructure:"REDIS_INVALIDATE_TIME" validate:"required"`
	}
	Xendit struct {
		APIKey                 string   `mapstructure:"XENDIT_API_KEY" validate:"required"`
		CallbackToken          string   `mapstructure:"XENDIT_CALLBACK_TOKEN" validate:"required"`
		InvoiceAPIURL          string   `mapstructure:"XENDIT_INVOICE_API_URL" validate:"required"`
		AvailablePaymentMethod []string `mapstructure:"XENDIT_AVAILABLE_PAYMENT_METHOD" validate:"required"`
	}

	Configuration struct {
		ServiceName   string      `mapstructure:"SERVICE_NAME"`
		Translation   Translation `mapstructure:",squash"`
		Postgres      Postgres    `mapstructure:",squash"`
		Redis         Redis       `mapstructure:",squash"`
		BindAddress   int         `mapstructure:"BIND_ADDRESS" validate:"required"`
		LogLevel      int         `mapstructure:"LOG_LEVEL" validate:"required"`
		JWTSecret     string      `mapstructure:"JWT_SECRET" validate:"required"`
		Xendit        Xendit      `mapstructure:",squash"`
		TraceEndpoint string      `mapstructure:"TRACE_ENDPOINT"`
		TraceRate     float64     `mapstructure:"TRACE_RATE"`
	}
)

func InitConfig(ctx context.Context) (*Configuration, error) {
	var cfg Configuration

	_, err := os.Stat(".env")
	if !os.IsNotExist(err) {
		viper.SetConfigFile(".env")

		if err := viper.ReadInConfig(); err != nil {
			logger.GetLogger(ctx).Errorf("failed to read config:%v", err)
			return nil, err
		}
	}

	viper.AutomaticEnv()

	if err := viper.Unmarshal(&cfg); err != nil {
		logger.GetLogger(ctx).Errorf("failed to bind config:%v", err)
		return nil, err
	}

	validate := validator.New()
	if err := validate.Struct(cfg); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			logger.GetLogger(ctx).Errorf("invalid config:%v", err)
		}
		logger.GetLogger(ctx).Errorf("failed to load config")
		return nil, err
	}

	logger.GetLogger(ctx).Infof("Config loaded: %+v", cfg)
	return &cfg, nil
}
