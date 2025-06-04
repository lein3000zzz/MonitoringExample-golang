package middleware

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/labstack/echo"
)

const RequestIDKey = "reqID"

func RequestIDMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(context echo.Context) error {
		b := make([]byte, 16)
		if _, err := rand.Read(b); err != nil {
			return err
		}

		reqID := hex.EncodeToString(b)
		context.Set("reqID", reqID)

		return next(context)
	}
}

func GetRequestID(c echo.Context) string {
	if v, ok := c.Get(RequestIDKey).(string); ok {
		return v
	}
	return ""
}
