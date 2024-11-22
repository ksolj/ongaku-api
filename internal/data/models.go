package data

import (
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Models struct {
	Tracks TrackModel
	Tokens TokenModel
	Users  UserModel
}

func NewModels(pool *pgxpool.Pool) Models {
	return Models{
		Tracks: TrackModel{Pool: pool},
		Tokens: TokenModel{Pool: pool},
		Users:  UserModel{Pool: pool},
	}
}
