package data

import (
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Models struct {
	Tracks      TrackModel
	Permissions PermissionModel
	Tokens      TokenModel
	Users       UserModel
}

func NewModels(pool *pgxpool.Pool, redis *redis.Client) Models {
	return Models{
		Tracks:      TrackModel{Pool: pool, Redis: redis},
		Permissions: PermissionModel{Pool: pool},
		Tokens:      TokenModel{Pool: pool},
		Users:       UserModel{Pool: pool},
	}
}
