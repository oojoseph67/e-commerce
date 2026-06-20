package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) healthCheck(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
}
