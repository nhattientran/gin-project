package server

import (
	"errors"
	"gin-project/internal/data"
	"gin-project/internal/validator"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"strings"
)

func (s *Server) health(c *gin.Context) {
	resp := make(map[string]string)
	resp["message"] = "healthy"
	resp["environment"] = s.config.env
	resp["version"] = version

	c.JSON(http.StatusOK, resp)
}

func (s *Server) healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, s.db.Health())
}

func (s *Server) createMovieHandler(c *gin.Context) {
	var movie data.Movie
	err := c.ShouldBindJSON(&movie)
	if err != nil {
		s.errorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	v := validator.New()
	if data.ValidateMovie(v, &movie); !v.Valid() {
		s.errorResponse(c, http.StatusBadRequest, v.Errors)
		return
	}

	// save to db
	s.models.Movies.Insert(c.Request.Context(), &movie)

	c.JSON(http.StatusCreated, movie)
}

func (s *Server) showMovieHandler(c *gin.Context) {
	id := c.Param("id")
	_, err := uuid.Parse(id)
	if err != nil {
		s.errorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	movie, err := s.models.Movies.Get(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			s.notFoundResponse(c)
			return
		}
		s.errorLog.PrintError(err, nil)
		s.errorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, movie)
}

func (s *Server) updateMovieHandler(c *gin.Context) {
	id := c.Param("id")
	var input struct {
		Title   *string       `json:"title"`
		Year    *int32        `json:"year"`
		Runtime *data.Runtime `json:"runtime"`
		Genres  []string      `json:"genres"`
	}
	err := c.ShouldBindJSON(&input)
	if err != nil {
		s.errorLog.PrintError(err, nil)
		s.errorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	movie, err := s.models.Movies.Get(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			s.notFoundResponse(c)
			return
		}
		s.errorLog.PrintError(err, nil)
		s.errorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	movie.ID, err = uuid.Parse(id)
	if err != nil {
		s.errorLog.PrintError(err, nil)
		s.errorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if input.Title != nil {
		movie.Title = *input.Title
	}
	// year
	if input.Year != nil {
		movie.Year = *input.Year
	}
	// runtime
	if input.Runtime != nil {
		movie.Runtime = *input.Runtime
	}
	// genres
	if len(input.Genres) > 0 {
		movie.Genres = input.Genres
	}
	// validate
	v := validator.New()
	if data.ValidateMovie(v, movie); !v.Valid() {
		s.errorLog.PrintError(err, nil)
		s.errorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// update
	err = s.models.Movies.Update(c.Request.Context(), movie)
	if err != nil {
		s.errorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, movie)
}

func (s *Server) deleteMovieHandler(c *gin.Context) {
	// get id
	id := c.Param("id")
	_, err := uuid.Parse(id)
	if err != nil {
		s.errorLog.PrintError(err, nil)
		s.errorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	// delete
	err = s.models.Movies.Delete(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			s.notFoundResponse(c)
			return
		}
		s.errorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "movie deleted"})
}

func (s *Server) listMoviesHandler(c *gin.Context) {
	var input struct {
		Title  string   `form:"title" default:""`
		Genres []string `form:"genres"`
		data.Filters
	}
	input.Filters = data.NewFilters()
	input.Filters.SortSafelist = []string{"id", "title", "year", "runtime", "-id", "-title", "-year", "-runtime"}
	// bind
	err := c.ShouldBindQuery(&input)
	if err != nil {
		s.errorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if len(input.Genres) > 0 {
		input.Genres = strings.Split(input.Genres[0], ",")
	}

	v := validator.New()
	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		s.errorResponse(c, http.StatusBadRequest, v.Errors)
		return
	}

	// get movies
	movies, err := s.models.Movies.List(c.Request.Context(), input.Title, input.Genres, &input.Filters)
	if err != nil {
		s.errorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	response := map[string]interface{}{
		"movies":   movies,
		"metadata": input,
	}

	c.JSON(http.StatusOK, response)
}
