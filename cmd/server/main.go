package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/braxtonkin/blogapi/internal/data"
	"log/slog"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

type application struct {
	logger       *slog.Logger
	blogModel    data.BlogModel
	commentModel data.CommentModel
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	/*	err := godotenv.Load()
		if err != nil {
			logger.Error("Unable to load .env file", "error", err)
			os.Exit(1)
		}
	*/
	serverPort := os.Getenv("SERVER_PORT")
	dbUser := os.Getenv("POSTGRES_USER")
	dbPass := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")
	dsn := fmt.Sprintf("host=database port=5432 user=%s password=%s dbname=%s sslmode=disable",
		dbUser,
		dbPass,
		dbName)

	db, err := openDB(dsn)
	if err != nil {
		logger.Error("Unable to connect to postgres", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	app := &application{
		logger:       logger,
		blogModel:    data.BlogModel{DB: db},
		commentModel: data.CommentModel{DB: db},
	}

	mux := app.routes()

	app.logger.Info(fmt.Sprintf("Starting server on port %s", serverPort))
	if err := http.ListenAndServe(fmt.Sprintf(":%s", serverPort), mux); err != nil {
		app.logger.Error("Error starting server", "error", err)
	}

}

// Opens the database
func openDB(dsn string) (*sql.DB, error) {
	time.Sleep(2 * time.Second)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	// Cancels if db doesn't respond in 5 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Ping to establish the connection to the database
	err = db.PingContext(ctx)
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
