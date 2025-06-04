package utils

import (
	"go.uber.org/zap"
	"io"
)

func BodyCloserWithSugaredLogger(body io.ReadCloser, logger *zap.SugaredLogger) {
	err := body.Close()
	if err != nil {
		logger.Errorf("failed to close body: %v", err)
		return
	}
}
