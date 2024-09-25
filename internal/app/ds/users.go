package ds

type Users struct {
	ID          int    `gorm:"primaryKey"`
	Login       string `gorm:"type:varchar(255)"`
	Password    string `gorm:"type:varchar(255)"`
	IsModerator bool   `gorm:"type:boolean"`
}
