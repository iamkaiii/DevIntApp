package ds

type MilkRequestsMeals struct {
	MilkRequestID int          `gorm:"primaryKey" json:"milk_request_id"`
	MealID        int          `gorm:"primaryKey" json:"meal_id"`
	MilkRequest   MilkRequests `gorm:"primaryKey;foreignKey:MilkRequestID"`
	Meal          Meals        `gorm:"primaryKey;foreignKey:MealID"`
	Amount        int          `json:"amount"`
}
