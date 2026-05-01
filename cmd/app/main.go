package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"short-link/cmd/internal/config"
	"short-link/cmd/internal/handler"
	"short-link/cmd/internal/repository"
	"short-link/cmd/internal/service"
	"time"

	_ "github.com/lib/pq"
)

func main() {
	cfg := config.Load()

	var repo repository.Repository

	if cfg.StorageType == "postgres" {
		if cfg.PostgresDSN == "" {
			log.Fatal("POSTGRES_DSN is required when STORAGE=postgres")
		}

		db, err := sql.Open("postgres", cfg.PostgresDSN)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := db.PingContext(ctx); err != nil {
			log.Fatal(err)
		}

		if err := repository.InitPostgres(ctx, db); err != nil {
			log.Fatal(err)
		}

		repo = repository.NewPostgresRepo(db)
	} else {
		repo = repository.NewMemoryRepo()
	}

	svc := service.New(repo)
	h := handler.New(svc)

	http.HandleFunc("/shorten", h.Shorten)
	http.HandleFunc("/", h.Redirect)

	log.Println("Server started on :" + cfg.ServerPort)
	log.Fatal(http.ListenAndServe(":"+cfg.ServerPort, nil))
}
