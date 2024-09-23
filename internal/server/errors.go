package server

import "github.com/gin-gonic/gin"

func (s *Server) errorResponse(message interface{}) map[string]interface{} {
	return gin.H{"error": message}
}

func (s *Server) notFoundResponse() map[string]interface{} {
	return s.errorResponse("Not found")
}
