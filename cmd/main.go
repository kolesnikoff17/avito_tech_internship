package main

import (
	"balance_api/config"
	v1 "balance_api/internal/controller/http/v1"
	"balance_api/internal/usecase"
	"balance_api/internal/usecase/report"
	"balance_api/internal/usecase/repository"
	"balance_api/pkg/httpserver"
	"balance_api/pkg/logger"
	"balance_api/pkg/postgres"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"os/signal"
	"syscall"
	// _ "github.com/swaggo/files"       // swagger embed files
	// _ "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

// @title           Balance API
// @version         1.0
// @description     Service for interactions with user's money accounts
// @contact.name   Developer
// @contact.email  vladkolesnikofff@gmail.com
// @license.name  MIT
// @license.url   https://github.com/kolesnikoff17/avito_tech_internship/blob/main/LICENSE
// @host      localhost:8080
// @BasePath  /v1
func main() {
	cfg := config.NewConfig()
	uri := config.DbParams(cfg)

	l, err := logger.New(cfg.Logger.Level)
	if err != nil {
		log.Fatalf("failed to build logger: %s", err)
	}

	db, err := postgres.New(uri, postgres.MaxConn(cfg.PG.MaxConn))
	if err != nil {
		l.Fatalf("failed to connect to db: %s", err)
	}

	r, err := report.New("reports/")
	if err != nil {
		l.Fatalf("failed to create report folder: %s", err)
	}

	useCase := usecase.New(repository.New(db), r)

	handler := gin.New()
	v1.NewRouter(handler, useCase, l)
	server := httpserver.New(handler)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case sig := <-interrupt:
		l.Infof("shutting down with signal: %s", sig)
	case err = <-server.Notify():
		l.Infof("server err: %s", err)
	}

	err = server.Shutdown()
	if err != nil {
		l.Infof("server shutdown err: %s", err)
	}
}
