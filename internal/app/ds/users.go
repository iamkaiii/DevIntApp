package ds

type Users struct {
	ID          int    `gorm:"primaryKey;autoIncrement" json:"id"`
	Login       string `gorm:"type:varchar(255)" json:"login"`
	Password    string `gorm:"type:varchar(255)" json:"password"`
	IsModerator bool   `gorm:"type:boolean" json:"is_moderator"`
}
