package main

import (
	"fmt"
	"log"
	"server/internal/api/middleware"
	"server/internal/pkg/comment/handler"
	commentrepo "server/internal/pkg/comment/repository"
	commentsvc "server/internal/pkg/comment/service"
	"server/internal/pkg/session"
	threadhttp "server/internal/pkg/thread/handler"
	threadrepo "server/internal/pkg/thread/repository"
	threadsvc "server/internal/pkg/thread/service"

	"go.uber.org/zap"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/labstack/echo"
)

func main() {
	e := echo.New()
	sessionSvc := session.NewService()

	zapLogger, err := zap.NewProduction()
	if err != nil {
		log.Fatal("Error initializing zap logger:", err)
	}

	defer func(zapLogger *zap.Logger) {
		err := zapLogger.Sync()
		if err != nil {
			fmt.Println("Error syncing zap logger:", err)
		}
	}(zapLogger)

	logger := zapLogger.Sugar()

	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	threadGroup := e.Group("/thread")

	threadGroup.Use(middleware.RequestIDMiddleware)
	threadGroup.Use(middleware.AccessLogMiddleware(logger))
	threadGroup.Use(middleware.ErrorLogMiddleware(logger))
	threadGroup.Use(middleware.AuthEchoMiddleware(sessionSvc))
	threadGroup.Use(middleware.MetricsMiddleware)

	threadRepo := threadrepo.NewRepository(logger)
	threadSvc := threadsvc.NewService(threadRepo)
	threadHandler := threadhttp.Handler{ThreadSvc: threadSvc}

	commentRepo := commentrepo.NewRepository(logger)
	commentSvc := commentsvc.NewService(commentRepo, threadRepo)
	commentHandler := handler.Handler{CommentSvc: commentSvc}

	threadGroup.GET("/:tid", threadHandler.GetThread)
	threadGroup.POST("", threadHandler.CreateThread)
	threadGroup.POST("/:tid/comment", commentHandler.Create)
	threadGroup.POST("/:tid/comment/:cid/like", commentHandler.Like)

	fmt.Print(e.Start(":8000"))
}
