package v1

import "go.uber.org/zap"

// Core can hold middleware, context, and logger
// In order to pass from main down to the api.
type CoreHandler struct {
	Log *zap.SugaredLogger
}
