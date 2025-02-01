package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/EugeneKrivoshein/fin_service/config"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

func main() {
	cfg, err := config.LoadConfig("config.env")
	if err != nil {
		log.Fatalf("ошибка загрузки конфигурации: %v", err)
	}

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUser, cfg.DBPass, cfg.DBHost, cfg.DBPort, cfg.DBName,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}
	defer db.Close()

	if err := goose.Up(db, "migrations"); err != nil {
		log.Fatal("Migration failed:", err)
	}

	log.Println("Migrations applied successfully")
}
