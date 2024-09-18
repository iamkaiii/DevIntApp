package ds

import "time"

type Carts struct {
	ID          uint `gorm:"primaryKey"`
	Status      int
	DateCreate  time.Time
	DateUpdate  time.Time
	DateFinish  time.Time
	CreatorID   uint
	ModeratorID uint
	Creator     Users `gorm:"foreignKey:CreatorID"`
	Moderator   Users `gorm:"foreignKey:ModeratorID"`
}
