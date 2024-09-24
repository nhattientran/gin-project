package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()
	r.Use(s.recoverPanic())
	r.Use(s.rateLimit())

	r.GET("/", s.health)
	r.GET("/v1/health", s.healthHandler)

	r.POST("/v1/movies", s.createMovieHandler)
	r.GET("/v1/movies/:id", s.showMovieHandler)
	r.PUT("/v1/movies/:id", s.updateMovieHandler)
	r.DELETE("/v1/movies/:id", s.deleteMovieHandler)
	r.GET("/v1/movies", s.listMoviesHandler)

	return r
}
