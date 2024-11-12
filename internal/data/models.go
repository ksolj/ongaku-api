package data

import (
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

type Models struct {
	Tracks TrackModel
}

func NewModels(pool *pgxpool.Pool) Models {
	return Models{
		Tracks: TrackModel{Pool: pool},
	}
}
