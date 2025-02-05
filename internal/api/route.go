package api

import (
	handler "github.com/EugeneKrivoshein/fin_service/internal/handlers"
	"github.com/gin-gonic/gin"
)

func SetupRouter(h *handler.Handler) *gin.Engine {
	r := gin.Default()

	// Роут для пополнения баланса
	r.POST("/deposit", h.HandleDeposit)

	// Роут для перевода денег
	r.POST("/transfer", h.HandleTransfer)

	// Роут для получения последних 10 транзакций
	// Например: GET /transactions?user_id=1
	r.GET("/transactions", h.HandleGetTransactions)

	return r
}
