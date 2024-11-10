package data

import (
	"time"
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
