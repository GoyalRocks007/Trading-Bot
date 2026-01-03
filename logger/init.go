package logger

import (
	"os"

	"go.uber.org/zap"
)

var Log *zap.Logger

func Init() {
	env := os.Getenv("APP_ENV")

	var cfg zap.Config
	if env == "prod" {
		cfg = zap.NewProductionConfig()
		cfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	} else {
		cfg = zap.NewDevelopmentConfig()
		cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	}

	var err error
	Log, err = cfg.Build()
	if err != nil {
		panic(err)
	}
}
