package gin

import (
	"context"

	"github.com/NidzamuddinMuzakki/test-sharing-session/go-lib-common/constant"
	"github.com/NidzamuddinMuzakki/test-sharing-session/go-lib-common/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader(constant.XRequestIdHeader)
		if requestID == "" {
			requestID = uuid.New().String()
		}
		c.Writer.Header().Add(constant.XRequestIdHeader, requestID)
		c.Request = c.Request.WithContext(logger.AddRequestID(c.Request.Context(), requestID))
		c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), constant.XRequestIdHeader, requestID))
		c.Next()
	}
}
