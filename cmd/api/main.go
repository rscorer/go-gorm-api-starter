package main

import (
	"api/db"
	"context"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type Api struct {
	db  *db.Repository
	log *logrus.Logger
}

func main() {
	log := logrus.New()
	log.SetReportCaller(true)
	log.SetLevel(logrus.DebugLevel)
	if err := godotenv.Load(); err != nil {
		log.WithError(err).Error("No .env file found")
	}
	appAddr := "0.0.0.0:" + os.Getenv("PORT")

	storage := db.CreateStorage(log)

	database, err := storage.InitDB(os.Getenv("DSN"))
	if err != nil {
		log.WithError(err).Fatal("opening db")
	}
	api := Api{db: database, log: log}

	router := api.setupRouter()

	srv := &http.Server{
		Addr:         appAddr,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      router,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.WithError(err).Fatal("listen")
		}
	}()
	//Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	<-c
	log.Println("Shutting down server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.WithError(err).Fatal("Server Shutdown")
	}
	log.Println("Server exited")

}
