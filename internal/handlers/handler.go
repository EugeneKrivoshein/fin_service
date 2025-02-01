package handler

import (
	"net/http"
	"strconv"

	service "github.com/EugeneKrivoshein/fin_service/internal/services"
	"github.com/gin-gonic/gin"
)

// Handler содержит сервис для вызова бизнес-логики
type Handler struct {
	service *service.Service
}

// NewHandler создаёт новый экземпляр Handler
func NewHandler(s *service.Service) *Handler {
	return &Handler{
		service: s,
	}
}

// DepositRequest представляет входные данные для пополнения баланса
type DepositRequest struct {
	UserID int64   `json:"user_id" binding:"required"`
	Amount float64 `json:"amount" binding:"required,gt=0"`
}

// TransferRequest представляет входные данные для перевода
type TransferRequest struct {
	SenderID   int64   `json:"sender_id" binding:"required"`
	ReceiverID int64   `json:"receiver_id" binding:"required"`
	Amount     float64 `json:"amount" binding:"required,gt=0"`
}

// HandleDeposit обрабатывает запрос на пополнение баланса
func (h *Handler) HandleDeposit(c *gin.Context) {
	var req DepositRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.Deposit(c.Request.Context(), req.UserID, req.Amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Баланс успешно пополнен"})
}

// HandleTransfer обрабатывает запрос на перевод денег
func (h *Handler) HandleTransfer(c *gin.Context) {
	var req TransferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.Transfer(c.Request.Context(), req.SenderID, req.ReceiverID, req.Amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Перевод успешно выполнен"})
}

// HandleGetTransactions обрабатывает запрос на получение последних транзакций
func (h *Handler) HandleGetTransactions(c *gin.Context) {
	// Можно передавать user_id как параметр запроса
	userIDParam := c.Query("user_id")
	if userIDParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}
	userID, err := strconv.ParseInt(userIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}

	transactions, err := h.service.GetTransactions(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, transactions)
}
