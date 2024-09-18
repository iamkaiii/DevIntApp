package ds

type Products struct {
	ID          uint `gorm:"primaryKey"`
	CardText    string
	Description string
	Status      int
	ImageUrl    string
}
