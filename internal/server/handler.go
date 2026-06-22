package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oojoseph67/ecommerce/internal/dto"
	"github.com/oojoseph67/ecommerce/internal/services"
	"github.com/oojoseph67/ecommerce/internal/utils/responses"
)

func (s *Server) healthCheck(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (s *Server) signup(ctx *gin.Context) {

}
