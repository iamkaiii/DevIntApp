package ds

type Meals struct {
	ID         int    `gorm:"primaryKey;autoIncrement" json:"id"`
	MealInfo   string `json:"meal_info"`
	MealWeight string `json:"meal_weight"`
	MealBrand  string `json:"meal_brand"`
	MealDetail string `json:"meal_detail"`
	ImageUrl   string `json:"image_url"`
	Status     bool   `json:"status"`
}
