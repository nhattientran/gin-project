package data

import (
	"gin-project/internal/database"
)

// Models create models struct base on MovieModel
type Models struct {
	Movies MovieModel
}

func NewModels(db database.Service) Models {
	return Models{
		Movies: MovieModel{
			service: db,
		},
	}
}
