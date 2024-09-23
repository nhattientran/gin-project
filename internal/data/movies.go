package data

import (
	"context"
	"database/sql"
	"errors"
	"gin-project/internal/database"
	"gin-project/internal/validator"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"time"
)

// MovieModel define a model db for movie
type MovieModel struct {
	service database.Service
}

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Movie struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"-"`
	Title     string    `json:"title"`
	Year      int32     `json:"year,omitempty"`
	Runtime   Runtime   `json:"runtime"`
	Genres    []string  `json:"genres,omitempty"`
	Version   int32     `json:"version"`
}

func (m *MovieModel) Insert(ctx context.Context, movie *Movie) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	query := `
			INSERT INTO movies (title, year, runtime, genres) 
			VALUES ($1, $2, $3, $4) 
			RETURNING id, created_at, version`
	args := []interface{}{movie.Title, movie.Year, movie.Runtime, movie.Genres}
	return m.service.DB().QueryRowContext(ctx, query, args...).Scan(&movie.ID, &movie.CreatedAt, &movie.Version)
}

func (m *MovieModel) Get(ctx context.Context, id string) (*Movie, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	var movie Movie
	query := `
			SELECT id, created_at, title, year, runtime, genres, version
			FROM movies
			WHERE id = $1`
	err := m.service.DB().QueryRowContext(ctx, query, id).
		Scan(&movie.ID,
			&movie.CreatedAt,
			&movie.Title,
			&movie.Year,
			&movie.Runtime,
			pq.Array(&movie.Genres),
			&movie.Version)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}
	return &movie, nil
}

func (m *MovieModel) Update(ctx context.Context, movie *Movie) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	query := `
			UPDATE movies
			SET title = $1, year = $2, runtime = $3, genres = $4, version = version + 1
			WHERE id = $5
			RETURNING version`
	args := []interface{}{movie.Title, movie.Year, movie.Runtime, pq.Array(movie.Genres), movie.ID}
	err := m.service.DB().QueryRowContext(ctx, query, args...).
		Scan(&movie.Version)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrEditConflict
		}
		return err
	}
	return nil
}

func (m *MovieModel) Delete(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	result, err := m.service.DB().ExecContext(ctx, `
			DELETE FROM movies
			WHERE id = $1`,
		id)
	rowsAffected, err := result.RowsAffected()
	// if no row affected, return error
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return err
}

func (m *MovieModel) List(ctx context.Context, title string, genres []string, filters *Filters) ([]*Movie, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	var movies []*Movie
	query := `
	
			SELECT id, created_at, title, year, runtime, genres, version
			FROM movies
			WHERE (LOWER(title) ILike LOWER($1) OR $1 = '') 
			AND (genres @> $2 OR $2 IS NULL)
			ORDER BY
				CASE WHEN $3 IN ('title', '-title') THEN title END ASC,
				CASE WHEN $3 IN ('year', '-year') THEN year END DESC,
				CASE WHEN $3 IN ('runtime', '-runtime') THEN runtime END DESC,
				CASE WHEN $3 IN ('created_at', '-created_at') THEN created_at END DESC
			LIMIT $4 OFFSET $5`

	rows, err := m.service.DB().QueryContext(ctx, query, title, pq.Array(genres), filters.Sort, filters.PageSize, filters.Page)
	if err != nil {
		return nil, err
	}

	// close rows
	defer rows.Close()
	// iterate over rows
	for rows.Next() {
		var movie Movie
		err = rows.Scan(&movie.ID,
			&movie.CreatedAt,
			&movie.Title,
			&movie.Year,
			&movie.Runtime,
			pq.Array(&movie.Genres),
			&movie.Version)
		if err != nil {
			return nil, err
		}
		movies = append(movies, &movie)
	}
	// check for errors from iterating over rows
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return movies, nil
}

func ValidateMovie(v *validator.Validator, movie *Movie) {
	v.Check(movie.Title != "", "title", "must be provided")
	v.Check(len(movie.Title) <= 500, "title", "must not be more than 500 bytes long")
	v.Check(movie.Year != 0, "year", "must be provided")
	v.Check(movie.Year >= 1888, "year", "must be greater than 1888")
	v.Check(movie.Year <= int32(time.Now().Year()), "year", "must not be in the future")
	v.Check(movie.Runtime != 0, "runtime", "must be provided")
	v.Check(movie.Runtime > 0, "runtime", "must be a positive integer")
	v.Check(movie.Genres != nil, "genres", "must be provided")
	v.Check(len(movie.Genres) >= 1, "genres", "must contain at least 1 genre")
	v.Check(len(movie.Genres) <= 5, "genres", "must not contain more than 5 genres")
	v.Check(validator.Unique(movie.Genres), "genres", "must not contain duplicate values")
}
