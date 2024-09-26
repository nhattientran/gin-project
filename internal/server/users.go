package server

import (
	"errors"
	"gin-project/internal/data"
	"gin-project/internal/validator"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *Server) registerUserHandler(c *gin.Context) {

	var input struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := c.ShouldBindJSON(&input)
	if err != nil {
		s.errorResponse(c, http.StatusBadRequest, err.Error())
		s.errorLog.PrintError(err, nil)
		return
	}
	user := data.User{
		Username:  input.Username,
		Email:     input.Email,
		Activated: false,
	}

	err = user.Password.Set(input.Password)
	if err != nil {
		s.errorResponse(c, http.StatusInternalServerError, err.Error())
		s.errorLog.PrintError(err, nil)
		return
	}

	v := validator.New()
	if data.ValidateUser(v, &user); !v.Valid() {
		s.errorResponse(c, http.StatusBadRequest, v.Errors)
		return
	}

	err = s.models.Users.Insert(c.Request.Context(), &user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateEmail):
			v.AddError("email", "a user with this email already exists")
			s.errorResponse(c, http.StatusConflict, v.Errors)
		default:
			s.serverErrorResponse(c, err)
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"user": user})
}
