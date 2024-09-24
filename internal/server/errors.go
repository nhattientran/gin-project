package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *Server) errorResponse(c *gin.Context, status int, message interface{}) {
	c.JSON(status, gin.H{"error": message})
	c.AbortWithStatus(status)
}

func (s *Server) serverErrorResponse(c *gin.Context, err error) {
	s.errorLog.PrintError(err, nil)
	message := "the server encountered a problem and could not process your request"
	s.errorResponse(c, http.StatusInternalServerError, message)
}
func (s *Server) notFoundResponse(c *gin.Context) {
	s.errorResponse(c, http.StatusNotFound, "resource not found")
}

func (s *Server) rateLimitExceededResponse(c *gin.Context) {
	s.errorResponse(c, http.StatusTooManyRequests, "rate limit exceeded")
}
