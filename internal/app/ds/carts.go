package ds

import "time"

type Milk_Requests struct {
	ID               int `gorm:"primaryKey"`
	Status           int
	DateCreate       time.Time
	DateUpdate       time.Time
	DateFinish       time.Time
	CreatorID        int
	ModeratorID      int
	Creator          Users `gorm:"foreignKey:CreatorID"`
	Moderator        Users `gorm:"foreignKey:ModeratorID"`
	RecipientName    string
	RecipientSurname string
	Address          string
	DeliveryDate     time.Time
}
