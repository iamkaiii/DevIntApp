package ds

type MilkRequestsMeals struct {
	MilkRequestID int          `gorm:"primaryKey"`
	MealID        int          `gorm:"primaryKey"`
	MilkRequest   MilkRequests `gorm:"primaryKey;foreignKey:MilkRequestID"`
	Meal          Meals        `gorm:"primaryKey;foreignKey:MealID"`
	Amount        int
}
