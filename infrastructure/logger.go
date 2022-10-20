package infrastructure

import plLogger "github.com/punk-link/logger"

func NewLoggerWithoutInjection() *plLogger.Logger {
	return &plLogger.Logger{}
}
