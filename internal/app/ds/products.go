package ds

type Meals struct {
	ID         int `gorm:"primaryKey"`
	MealInfo   string
	MealWeight string
	MealDetail string
	ImageUrl   string
	Status     bool
}
