package main

import (
	"context"
	"net/http"
	"os"

	"github.com/guruorgoru/newsguru/pkg/logs"
	"github.com/guruorgoru/newsguru/pkg/models"
	"github.com/guruorgoru/newsguru/pkg/router"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		logs.Error.Fatal("Error loading .env file")
	}

	portString := os.Getenv("PORT")
	if portString == "" {
		logs.Error.Fatal("Port is not set in the environment variables")
	}

	dsn := os.Getenv("DSN")

	db, err := openDB(dsn)
	if err != nil {
		logs.Error.Fatal(err)
	}
	if dsn == "" {
		logs.Error.Fatal("DSN is not set in the environment variables")
	}
	defer db.Close()

	app := &models.NewsModel{
		DB: db,
	}

	if err = db.Ping(context.Background()); err != nil {
		logs.Error.Fatal("Failed to connect to the database:", err)
	} else {
		logs.Info.Println("Connected to the database successfully")
	}
	router := router.NewsRouter(app)
	server := &http.Server{
		Addr:     ":" + portString,
		ErrorLog: logs.Error,
		Handler:  router,
	}
	logs.Info.Println("Server starting on port: ", portString)
	err1 := server.ListenAndServe()
	logs.Error.Fatal(err1)
}

func openDB(dsn string) (*pgxpool.Pool, error) {
	db, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, err
	}
	err = db.QueryRow(context.Background(), "SELECT 1").Scan(new(int))
	if err != nil {
		return nil, err
	}
	return db, nil
}
