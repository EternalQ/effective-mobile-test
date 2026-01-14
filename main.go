package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	_ "github.com/EternalQ/effective-mobile-test/docs"

	"github.com/EternalQ/effective-mobile-test/pkg/api"
	"github.com/EternalQ/effective-mobile-test/pkg/service"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
	"github.com/spf13/viper"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

var (
	dbUser string
	dbPass string
	dbName string
	dbHost string
	logLvl int
)

func readEnv() {
	viper.AutomaticEnv()

	viper.SetDefault("DB_USER", "admin")
	dbUser = viper.GetString("DB_USER")

	viper.SetDefault("DB_PASS", "admin")
	dbPass = viper.GetString("DB_PASS")

	viper.SetDefault("DB_NAME", "billing")
	dbName = viper.GetString("DB_NAME")

	viper.SetDefault("DB_HOST", "localhost:5432")
	dbHost = viper.GetString("DB_HOST")

	viper.SetDefault("LOG_LVL", -4)
	logLvl = viper.GetInt("LOG_LVL")
}

// @title Effective Mobile Test API
// @version 1.0
// @description This is a sample server for the Effective Mobile Test.
// @host localhost:8080
// @BasePath /api
func main() {
	readEnv()

	logHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.Level(logLvl),
	})
	log := slog.New(logHandler)

	log.Info("App started")

	pgsStr := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", dbUser, dbPass, dbHost, dbName)
	pgs, err := sqlx.Connect("postgres", pgsStr)
	if err != nil {
		log.Error("Can't connect to PostgreSQL, check .env")
		os.Exit(0)
	}
	log.Info("PotgreSQL connected")

	subServ := service.NewSubscriptionService(pgs, log)
	log.Info("Subscription service created")

	router := mux.NewRouter()

	c := cors.AllowAll()
	router.Use(c.Handler)

	api.StartServer(log, subServ, router)
	
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
	
	log.Info("Server started, listening on 8080")
	http.ListenAndServe(":8080", router)
}
