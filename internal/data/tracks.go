package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ksolj/ongaku-api/internal/data/validator"
)

type TrackModel struct {
	Pool *pgxpool.Pool
}

func (t TrackModel) Insert(track *Track) error {
	query := `
        INSERT INTO tracks (name, duration, artists, album, tabs) 
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id, created_at, version`

	args := []any{track.Name, track.Duration, track.Artists, track.Album, track.Tabs}

	return t.Pool.QueryRow(context.Background(), query, args...).Scan(&track.ID, &track.CreatedAt, &track.Version)
}

func (t TrackModel) Get(id int64) (*Track, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
        SELECT id, created_at, name, duration, artists, album, tabs, version
        FROM tracks
        WHERE id = $1`

	var track Track

	err := t.Pool.QueryRow(context.Background(), query, id).Scan(
		&track.ID,
		&track.CreatedAt,
		&track.Name,
		&track.Duration,
		&track.Artists,
		&track.Album,
		&track.Tabs,
		&track.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &track, nil
}

func (t TrackModel) Update(track *Track) error {
	query := `
        UPDATE tracks 
        SET name = $1, duration = $2, artists = $3, album = $4, tabs = $5, version = version + 1
        WHERE id = $6
        RETURNING version`

	args := []any{
		track.Name,
		track.Duration,
		track.Artists,
		track.Album,
		track.Tabs,
		track.ID,
	}

	return t.Pool.QueryRow(context.Background(), query, args...).Scan(&track.Version)
}

func (t TrackModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
        DELETE FROM tracks
        WHERE id = $1`

	result, err := t.Pool.Exec(context.Background(), query, id)
	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()

	// If no rows were affected, we know that the tracks table didn't contain a record
	// with the provided ID at the moment we tried to delete it. In that case we
	// return an ErrRecordNotFound error.
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

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
