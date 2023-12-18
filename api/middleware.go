package api

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	log "github.com/slamchillz/xchange/logger"
	"github.com/rs/zerolog"
)


const (
	AUTHENTICATIONHEADER = "authorization"
	AUTHENTICATIONSCHEME = "bearer"
	AUTHENTICATIONCONTEXTKEY = "authenticated_customer"
)

type correlationId string

type logResponseWriter struct {
	gin.ResponseWriter
	statusCode int
}

func newlogResponseWriter(w gin.ResponseWriter) *logResponseWriter {
	return &logResponseWriter{w, http.StatusOK}
}

func (w *logResponseWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}


func (server *Server) Authenticate(ctx *gin.Context) {
	authHeader := ctx.GetHeader(AUTHENTICATIONHEADER)
	if len(authHeader) <= len(AUTHENTICATIONSCHEME) {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing or incomplete authentication header"})
		return
	}
	authHeaderValues := strings.Fields(authHeader)
	if len(authHeaderValues) != 2 {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "unrecognized authentication header format"})
		return
	}
	if strings.ToLower(authHeaderValues[0]) != AUTHENTICATIONSCHEME {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unsupported authentication scheme"})
		return
	}
	accessToken := authHeaderValues[1]
	customer, err := server.token.VerifyToken(accessToken)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid access token"})
		return
	}
	ctx.Set(AUTHENTICATIONCONTEXTKEY, customer)
	ctx.Next()
}

func (server *Server) HTTPLogger(ctx *gin.Context) {
	start := time.Now()
	logger := log.GetLogger()
	correlation_id := uuid.New().String()
	correlationContext := context.WithValue(ctx.Request.Context(), correlationId("correlation_id"), correlation_id)
	ctx.Request = ctx.Request.WithContext(correlationContext)
	logger.UpdateContext(func(c zerolog.Context) zerolog.Context {
		return c.Str("correlation_id", correlation_id)
	})
	ctx.Request.Header.Add("correlation_id", correlation_id)
	lrw := newlogResponseWriter(ctx.Writer)
	ctx.Writer = lrw
	ctx.Request = ctx.Request.WithContext(logger.WithContext(ctx.Request.Context()))
	defer func() {
		panicValue := recover()
		if panicValue != nil {
			logger.Error().
				Str("method", ctx.Request.Method).
				Str("path", ctx.Request.URL.Path).
				Str("remote_addr", ctx.Request.RemoteAddr).
				Int("status_code", lrw.statusCode).
				Dur("latency", time.Since(start)).
				Msg("request handled with panic")
			panic(panicValue)
		}
		logger.Info().
			Str("method", ctx.Request.Method).
			Str("path", ctx.Request.URL.Path).
			Str("remote_addr", ctx.Request.RemoteAddr).
			Int("status_code", lrw.statusCode).
			Dur("latency", time.Since(start)).
			Msg("request handled")
	}()
	ctx.Next()
}
