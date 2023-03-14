package v1

import "go.uber.org/zap"

// Core can hold middleware, context, and logger, etc.
// This struct will bind the API associated methods
// found in apps.go and dashboard.go and is instantiated
// in routes.go
type CoreHandler struct {
	Log *zap.SugaredLogger
}
