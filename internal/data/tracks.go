package data

import (
	"time"
)

type Track struct {
	ID        int64
	CreatedAt time.Time
	Name      string
	Duration  int32 // in ms
	Artists   []string
	Album     string
	Tabs      string // for now tabs' type is string but later it'll be changed
	Version   int32  // keep track of how many times someone updated track info (this field may be deleted in the future)
}
