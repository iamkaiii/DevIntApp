package ds

type CartsNMeals struct {
	CartID      int
	ChildMealID int
	Cart        Carts `gorm:"primaryKey;foreignKey:CartID"`
	ChildMeal   Meals `gorm:"primaryKey;foreignKey:ChildMealID"`
	OrderO      int
}
