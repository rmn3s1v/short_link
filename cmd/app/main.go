package main

import (
	"database/sql"
	"log"
	"net/http"
	"short-link/cmd/internal/config"
	"short-link/cmd/internal/handler"
	"short-link/cmd/internal/repository"
	"short-link/cmd/internal/service"

	_ "github.com/lib/pq"
)

func main() {
	cfg := config.Load()

	var repo repository.Repository

	if cfg.StorageType == "postgres" {
		db, err := sql.Open("postgres", cfg.PostgresDNS)
		if err != nil {
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
