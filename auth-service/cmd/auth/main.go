package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"xm-auth/internal/database"
	"xm-auth/server"

	"github.com/subosito/gotenv"
)

const DEV_ENV string = "dev"

func init() {
	err := gotenv.Load()
	if err != nil {
		log.Fatalf("gotenv.Load() error: %s\n", err.Error())
	}
}

func main() {
	log.Println("Starting auth service...")

	env := os.Getenv("ENV")
	if env == "" {
		env = DEV_ENV
	}
	isDevMode := env == DEV_ENV

	port := os.Getenv("PORT")

	connStr := os.Getenv("MONGO_CONN_STR")
	db, err := database.New(connStr, context.Background())
	if err != nil {
		log.Fatalf("DB initialization error: %s", err.Error())
	}

	server := &http.Server{
		Addr:              port,
		Handler:           server.SetupRouter(db, isDevMode),
		ReadHeaderTimeout: 1 * time.Second,
	}

	go func() {
		log.Printf("Server listening on port %s\n", port)
		err := server.ListenAndServe()
		if err != nil {
			log.Printf("server.ListenAndServe() error: %s\n", err.Error())
		}
	}()

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT, os.Interrupt)

	<-stop

	stopCtx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	log.Println("Shutting down auth service")

	if err = db.MdbClient.Disconnect(stopCtx); err != nil {
		log.Printf("DB client disconnect error: %s\n", err.Error())
	}

	err = server.Shutdown(stopCtx)
	if err != nil {
		log.Printf("server.Shutdown() error: %s\n", err.Error())
	}
}
