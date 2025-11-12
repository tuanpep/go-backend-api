package logger

import (
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Logger wraps logrus.Logger
type Logger struct {
	*logrus.Logger
}

// NewLogger creates a new logger instance
func NewLogger(level string) *Logger {
	log := logrus.New()

	// Set log level
	switch level {
	case "debug":
		log.SetLevel(logrus.DebugLevel)
	case "info":
		log.SetLevel(logrus.InfoLevel)
	case "warn":
		log.SetLevel(logrus.WarnLevel)
	case "error":
		log.SetLevel(logrus.ErrorLevel)
	default:
		log.SetLevel(logrus.InfoLevel)
	}

	// Set JSON formatter for structured logging
	log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	})

	// Set output to stdout
	log.SetOutput(os.Stdout)

	return &Logger{Logger: log}
}

// GinLogger returns a gin.HandlerFunc for logging HTTP requests
func (l *Logger) GinLogger() gin.HandlerFunc {
	return gin.LoggerWithConfig(gin.LoggerConfig{
		Formatter: func(param gin.LogFormatterParams) string {
			l.WithFields(logrus.Fields{
				"timestamp":  param.TimeStamp.Format(time.RFC3339),
				"method":     param.Method,
				"path":       param.Path,
				"status":     param.StatusCode,
				"latency":    param.Latency.String(),
				"client_ip":  param.ClientIP,
				"user_agent": param.Request.UserAgent(),
				"error":      param.ErrorMessage,
			}).Info("HTTP Request")
			return ""
		},
		Output: os.Stdout,
	})
}

// GinRecovery returns a gin.HandlerFunc for recovering from panics
func (l *Logger) GinRecovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		l.WithFields(logrus.Fields{
			"error":  recovered,
			"path":   c.Request.URL.Path,
			"method": c.Request.Method,
		}).Error("Panic recovered")

		c.AbortWithStatus(500)
	})
}

// WithContext creates a logger with context
func (l *Logger) WithContext(ctx *gin.Context) *logrus.Entry {
	fields := logrus.Fields{}

	if userID, exists := ctx.Get("user_id"); exists {
		fields["user_id"] = userID
	}

	if username, exists := ctx.Get("username"); exists {
		fields["username"] = username
	}

	if requestID := ctx.GetHeader("X-Request-ID"); requestID != "" {
		fields["request_id"] = requestID
	}

	return l.Logger.WithFields(fields)
}
