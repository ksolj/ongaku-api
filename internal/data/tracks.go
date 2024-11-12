package data

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ksolj/ongaku-api/internal/data/validator"
)

type TrackModel struct {
	Pool *pgxpool.Pool
}

func (m TrackModel) Insert(track *Track) error {
	query := `
        INSERT INTO tracks (name, duration, artists, album, tabs) 
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id, created_at, version`

	args := []any{track.Name, track.Duration, track.Artists, track.Album, track.Tabs}

	return m.Pool.QueryRow(context.Background(), query, args...).Scan(&track.ID, &track.CreatedAt, &track.Version)
}

func (m TrackModel) Get(id int64) (*Track, error) {
	return nil, nil
}

func (m TrackModel) Update(track *Track) error {
	return nil
}

func (m TrackModel) Delete(id int64) error {
	return nil
}

type Track struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"`
	Name      string    `json:"name"`
	Duration  int32     `json:"duration"` // In seconds
	Artists   []string  `json:"artists"`
	Album     string    `json:"album"`
	Tabs      []Tab     `json:"tabs,omitempty"` // TODO: add check for tabs not null
	Version   int32     `json:"version"`        // Keep track of how many times someone updated track info (this field may be deleted in the future)
}

func ValidateTrack(v *validator.Validator, track *Track) {
	v.Check(track.Name != "", "name", "must be provided")
	v.Check(len(track.Name) <= 500, "name", "must not be more than 500 bytes long")

	v.Check(track.Duration != 0, "duration", "must be provided")
	v.Check(track.Duration > 0, "duration", "must be a positive integer") // TODO: possible overflow???

	v.Check(track.Artists != nil, "artists", "must be provided")
	v.Check(len(track.Artists) >= 1, "artists", "must contain at least 1 artist")
	v.Check(validator.Unique(track.Artists), "artists", "must not contain duplicate values")
}
