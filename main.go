package main

import (
	"context"
	"github.com/aaanger/todo/pkg/repository"
	"github.com/aaanger/todo/pkg/routes"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type Server struct {
	httpServer *http.Server
}

func (srv *Server) run(port string, handler http.Handler) error {
	srv.httpServer = &http.Server{
		Addr:    port,
		Handler: handler,
	}
	return srv.httpServer.ListenAndServe()
}

func (srv *Server) shutdown(ctx context.Context) error {
	return srv.httpServer.Shutdown(ctx)
}

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))
	err := godotenv.Load()
	if err != nil {
		logrus.Fatalf("Error init config: %s", err.Error())
	}
	db, err := repository.NewPostgresConfig(repository.PostgresConfig{
		Host:     os.Getenv("PSQL_HOST"),
		Port:     os.Getenv("PSQL_PORT"),
		Username: os.Getenv("PSQL_USER"),
		Password: os.Getenv("PSQL_PASSWORD"),
		DBName:   os.Getenv("PSQL_DBNAME"),
		SSLMode:  os.Getenv("PSQL_SSLMODE"),
	})
	if err != nil {
		logrus.Fatalf("Error creating database: %s", err.Error())
	}

	srv := new(Server)
	go func() {
		err = srv.run(":3000", routes.PathHandler(db))
		if err != nil {
			logrus.Fatalf("Error run server: %s", err.Error())
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	err = srv.shutdown(context.Background())
	if err != nil {
		logrus.Errorf("error shutting down server: %w", err)
	}

	err = db.Close()
	if err != nil {
		logrus.Errorf("error closing database: %w", err)
	}
}
