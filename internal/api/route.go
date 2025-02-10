package api

import (
	_ "github.com/EugeneKrivoshein/fin_service/docs"
	handler "github.com/EugeneKrivoshein/fin_service/internal/handlers"
	"github.com/gin-gonic/gin"
	files "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(h *handler.Handler) *gin.Engine {
	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(files.Handler))

	// Роут для пополнения баланса
	r.POST("/deposit", h.HandleDeposit)

	// Роут для перевода денег
	r.POST("/transfer", h.HandleTransfer)

	// Роут для получения последних 10 транзакций
	// Например: GET /transactions?user_id=1
	r.GET("/transactions", h.HandleGetTransactions)

	return r
}
