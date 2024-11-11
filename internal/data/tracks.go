package data

import (
	"time"

	"github.com/ksolj/ongaku-api/internal/data/validator"
)

type Track struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"`
	Name      string    `json:"name"`
	Duration  int32     `json:"duration"` // in seconds
	Artists   []string  `json:"artists"`
	Album     string    `json:"album"`
	Tabs      []Tab     `json:"tabs,omitempty"`
	Version   int32     `json:"version"` // keep track of how many times someone updated track info (this field may be deleted in the future)
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
