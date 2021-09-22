package http

import (
	"belajar-bwa/domain"
	"belajar-bwa/helper"
	"belajar-bwa/transaction"
	_middleware "belajar-bwa/user/delivery/http/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
)

type TransactionHandler struct {
	transactionUsecase domain.TransactionUsecase
	UserUsecase        domain.UserUsecase
	AuthUsecase        domain.JWTService
}

func NewTransactionHandler(c *gin.RouterGroup, tu domain.TransactionUsecase, uu domain.UserUsecase, au domain.JWTService) {
	handler := &TransactionHandler{transactionUsecase: tu, AuthUsecase: au, UserUsecase: uu}

	newMiddleware := _middleware.NewUserMiddleware(au, uu)

	c.GET("/campaigns/:id/transactions", newMiddleware.Auth(), handler.GetCampaignTransactions)
	c.GET("/transactions", newMiddleware.Auth(), handler.GetUserTransactions)
	c.POST("/transactions", newMiddleware.Auth(), handler.CreateTransaction)
	c.POST("/transactions/notifications", handler.GetNotification)
}

func (h *TransactionHandler) GetCampaignTransactions(c *gin.Context) {
	var input domain.GetCampaignTransactionsInput

	err := c.ShouldBindUri(&input)
	if err != nil {
		response := helper.APIResponse("Error to get campaign's transactions", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	currentUser := c.MustGet("currentUser").(domain.User)

	input.User = currentUser

	transactions, err := h.transactionUsecase.GetTransactionsByCampaignID(input)
	if err != nil {
		errors := err.Error()
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Error to get campaign's transactions", http.StatusBadRequest, "error", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Campaign's Transactions", http.StatusOK, "success", transaction.FormatCampaignTransactions(transactions))
	c.JSON(http.StatusOK, response)
}

func (h *TransactionHandler) GetUserTransactions(c *gin.Context) {
	currentUser := c.MustGet("currentUser").(domain.User)
	userID := currentUser.ID

	transactions, err := h.transactionUsecase.GetTransactionsByUserID(userID)
	if err != nil {
		response := helper.APIResponse("Error to get user's transactions", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("User's Transactions", http.StatusOK, "success", transaction.FormatUserTransactions(transactions))
	c.JSON(http.StatusOK, response)
}

func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
	var input domain.CreateTransactionInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Error to create transactions", http.StatusBadRequest, "error", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	currentUser := c.MustGet("currentUser").(domain.User)
	input.User = currentUser

	newTransaction, err := h.transactionUsecase.CreateTransaction(input)
	if err != nil {
		response := helper.APIResponse("Error to create transactions", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Success to create campaign", http.StatusOK, "success", transaction.FormatTransaction(newTransaction))
	c.JSON(http.StatusOK, response)
}

func (h TransactionHandler) GetNotification(c *gin.Context) {
	var input domain.TransactionNotificationInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		response := helper.APIResponse("Failed to process notification", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	err = h.transactionUsecase.ProcessPayment(input)
	if err != nil {
		response := helper.APIResponse("Failed to process notification", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	c.JSON(http.StatusOK, input)
}
