package handlers

import (
	"fmt"
	internalManager "github.com/DipandaAser/tg-bot-storage/internal/manager"
	v1 "github.com/DipandaAser/tg-bot-storage/pkg/models/v1"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func UploadFile(ctx *gin.Context) {

	var err error
	var query v1.UploadFileQuery
	query.FileName = ctx.Query("file_name")
	if query.ChatId, err = strconv.ParseInt(ctx.Query("chat_id"), 10, 64); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "bad query parameters"})
		return
	}

	fileIdentifier, err := internalManager.
		GetDefaultManager().
		UploadFileReader(query.ChatId, query.FileName, ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("could not upload file. err: %s", err.Error()),
		})
		return
	}

	ctx.JSON(http.StatusOK, fileIdentifier)
}

func DownloadFile(ctx *gin.Context) {
	var query v1.DownloadFileQuery
	var err error

	if query.MessageIdentifier.MessageId, err = strconv.Atoi(ctx.Query("msg_id")); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "bad query parameters"})
		return
	}

	if query.MessageIdentifier.ChatId, err = strconv.ParseInt(ctx.Query("chat_id"), 10, 64); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "bad query parameters"})
		return
	}

	if query.DraftChatId, err = strconv.ParseInt(ctx.Query("draft_chat_id"), 10, 64); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "bad query parameters"})
		return
	}

	result, err := internalManager.
		GetDefaultManager().
		DownloadFileReader(query.MessageIdentifier, query.DraftChatId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("could not download file. err: %s", err.Error()),
		})
		return
	}

	extraHeaders := map[string]string{
		"Content-Disposition": fmt.Sprintf("attachment; filename=%s", result.FileInfo.Name),
	}

	ctx.DataFromReader(
		http.StatusOK,
		result.FileInfo.Size,
		result.FileInfo.ContentType,
		result.Data,
		extraHeaders)
}
