package handler

import (
	"net/http"
	"strconv"

	"github.com/EugeneKrivoshein/fin_service/internal/postgres"
	service "github.com/EugeneKrivoshein/fin_service/internal/services"
	"github.com/gin-gonic/gin"
)

var _ postgres.Transaction

type Handler struct {
	service *service.Service
}

func NewHandler(s *service.Service) *Handler {
	return &Handler{
		service: s,
	}
}

type DepositRequest struct {
	UserID int64   `json:"user_id" binding:"required"`
	Amount float64 `json:"amount" binding:"required,gt=0"`
}

type TransferRequest struct {
	SenderID   int64   `json:"sender_id" binding:"required"`
	ReceiverID int64   `json:"receiver_id" binding:"required"`
	Amount     float64 `json:"amount" binding:"required,gt=0"`
}

// HandleDeposit godoc
// @Summary Пополнение баланса
// @Description Позволяет пользователю пополнить свой баланс
// @Tags Баланс
// @Accept json
// @Produce json
// @Param input body DepositRequest true "Данные для пополнения"
// @Success 200 {object} map[string]string "Баланс успешно пополнен"
// @Failure 400 {object} map[string]string "Ошибка валидации"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /deposit [post]
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

// HandleTransfer godoc
// @Summary Перевод денег
// @Description Позволяет пользователю перевести деньги другому пользователю
// @Tags Транзакции
// @Accept json
// @Produce json
// @Param input body TransferRequest true "Данные для перевода"
// @Success 200 {object} map[string]string "Перевод успешно выполнен"
// @Failure 400 {object} map[string]string "Ошибка валидации"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /transfer [post]
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

// HandleGetTransactions godoc
// @Summary Получение последних 10 транзакций
// @Description Возвращает список последних 10 транзакций пользователя
// @Tags Транзакции
// @Accept json
// @Produce json
// @Param user_id query int true "ID пользователя"
// @Success 200 {array} postgres.Transaction "Список транзакций"
// @Failure 400 {object} map[string]string "Ошибка валидации"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /transactions [get]
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
