package main

import (
	"log"
	"net/http"

	_ "github.com/EugeneKrivoshein/fin_service/docs"
	route "github.com/EugeneKrivoshein/fin_service/internal/api"
	handler "github.com/EugeneKrivoshein/fin_service/internal/handlers"
	"github.com/EugeneKrivoshein/fin_service/internal/postgres"
	repo "github.com/EugeneKrivoshein/fin_service/internal/postgres"
	service "github.com/EugeneKrivoshein/fin_service/internal/services"
)

// @title Финансовый сервис API
// @version 1.0
// @description API для управления балансом и переводами денег
// @host localhost:8080
// @BasePath /
func main() {
	pgxProvider, err := postgres.NewPGXProvider()
	if err != nil {
		log.Fatalf("Ошибка создания подключения к БД: %v", err)
	}
	defer pgxProvider.Close()

	repository := repo.NewRepository(pgxProvider.Pool)
	serviceLayer := service.NewService(repository)
	handlerLayer := handler.NewHandler(serviceLayer)
	router := route.SetupRouter(handlerLayer)

	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}
}
