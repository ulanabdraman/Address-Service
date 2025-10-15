package http

import (
	"AddressService/internal/domains/message/model"
	"AddressService/internal/domains/message/usecase"
	"github.com/gin-gonic/gin"
	"net/http"
)

type MessageHandler struct {
	usecase usecase.MessageUseCase
}

func NewMessageHandler(uc usecase.MessageUseCase) *MessageHandler {
	return &MessageHandler{usecase: uc}
}

// 📥 Обработка одного сообщения (реалтайм)
func (h *MessageHandler) Handle(c *gin.Context) {
	var msg model.Message
	if err := c.ShouldBindJSON(&msg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	if err := h.usecase.ProcessMessage(c.Request.Context(), &msg); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "processing failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "message processed"})
}

// 📥 Обработка массива сообщений для отчёта (/report)
func (h *MessageHandler) HandleReport(c *gin.Context) {
	var msgs []*model.Message
	if err := c.ShouldBindJSON(&msgs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON array"})
		return
	}

	updated, err := h.usecase.ProcessMessages(c.Request.Context(), msgs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "report processing failed"})
		return
	}

	c.JSON(http.StatusOK, updated)
}
