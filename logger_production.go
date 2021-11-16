//go:build production
// +build production

package main

import (
	"go.uber.org/zap"
)

func zapConfig() zap.Config {
	return zap.NewProductionConfig()
}
