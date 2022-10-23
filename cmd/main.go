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
)

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
