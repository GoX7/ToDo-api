package main

import (
	"fmt"
	"log"
	"net/http"
	"to-do/internal/config"
	"to-do/internal/http/controls/handlers"
	"to-do/internal/logger"
	"to-do/internal/sqlite"

	"github.com/go-chi/chi/v5"
)

func main() {
	cfg, err := config.Load() // load config
	if err != nil {
		log.Fatal(err)
	}

	logs, err := logger.New(cfg) // load logger
	if err != nil {
		log.Fatal(err)
	}

	db, err := sqlite.New() // New database
	if err != nil {
		logs.Server.Warn("*Error connect to database")
		logs.Sqlite.Error(fmt.Sprint("*Error connect to database:", err))
		log.Fatal(err)
	}

	logs.Server.Info("Start 1/2")
	logs.Server.Debug("Connect to database")
	logs.Sqlite.Info("Connect to database")

	Start(cfg, logs, db) // Function start server
}

func Start(cfg *config.Config, logger *logger.Logger, db *sqlite.Database) {
	logger.Server.Debug("Function start...")

	router := chi.NewRouter()                             // router
	handlers.NewHandler(db, cfg, logger).Register(router) // register handlers
	logger.Server.Info("Start 2/2")

	server := &http.Server{
		Addr:         cfg.Server.Addr,
		Handler:      router,
		WriteTimeout: cfg.Server.Wto,
		ReadTimeout:  cfg.Server.Rto,
	}

	log.Print("Start")
	logger.Server.Info("Start server")
	err := server.ListenAndServe() // start server
	if err != nil {
		logger.Server.Error(fmt.Sprint("Stop server, error:", err))
	}
}
