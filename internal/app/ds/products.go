package ds

type Meals struct {
	ID         int `gorm:"primaryKey"`
	MealInfo   string
	MealWeight string
	MealBrand  string
	MealDetail string
	ImageUrl   string
	Status     bool
}
