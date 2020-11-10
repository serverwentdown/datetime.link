// +build !production

package main

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func zapConfig() zap.Config {
	dev := zap.NewDevelopmentConfig()
	dev.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	return dev
}
