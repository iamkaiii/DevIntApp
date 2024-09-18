package ds

type CartsNProducts struct {
	CartID    uint
	ProductID uint
	Cart      Carts    `gorm:"primaryKey;foreignKey:CartID"`
	Product   Products `gorm:"primaryKey;foreignKey:ProductID"`
	Order     int
}
