package ds

type MilkRequestsNMeals struct {
	CartID      int           `gorm:"primaryKey"`
	ChildMealID int           `gorm:"primaryKey"`
	Cart        Milk_Requests `gorm:"primaryKey;foreignKey:CartID"`
	ChildMeal   Meals         `gorm:"primaryKey;foreignKey:ChildMealID"`
	OrderO      int
}
