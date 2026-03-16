package Database

import "time"

type FavoriteFlat struct {
	Id        uint64
	UserTgID  int64
	FlatCode  string
	CreatedAt time.Time
}
