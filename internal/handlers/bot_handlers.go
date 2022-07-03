package handlers

import (
	"fmt"
	internalManager "github.com/DipandaAser/tg-bot-storage/internal/manager"
	v1 "github.com/DipandaAser/tg-bot-storage/pkg/models/v1"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func AddBotToManager(ctx *gin.Context) {

	var query v1.AddBotQuery
	if err := ctx.ShouldBind(&query); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "bad query parameters"})
		return
	}

	if strings.TrimSpace(query.Token) == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "token is required"})
		return
	}

	err := internalManager.GetDefaultManager().AddBot(query.Token)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("could not add bot. err: %s", err.Error()),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "bot added successfully",
	})
}
