package repository

import (
	"DevIntApp/internal/app/ds"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
)

//cart status: 0 - черновик, 1 - сформирован, 2 - завершен, 3 - удалён, 4 - отклонён

type Repository struct {
	db *gorm.DB
}

func New(dsn string) (*Repository, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return &Repository{
		db: db,
	}, nil
}

func (r *Repository) GetAllMeals() ([]ds.Meals, error) {
	var prods []ds.Meals
	err := r.db.Where("status=true").Find(&prods).Error
	if err != nil {
		return nil, err
	}
	return prods, nil
}
func (r *Repository) GetMealByID(mealID int) ([]ds.Meals, error) {
	var meal []ds.Meals
	err := r.db.Where("id = ?", mealID).Find(&meal).Error
	if err != nil {
		return nil, err
	}
	return meal, nil
}

func (r *Repository) GetMealByMealInfo(cardText string) ([]ds.Meals, error) {
	var product []ds.Meals
	err := r.db.Where("meal_info LIKE ?", "%"+cardText+"%").First(&product).Error
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (r *Repository) GetWorkingCart() ([]ds.Carts, error) {
	var cart []ds.Carts
	err := r.db.Where("status=0").Find(&cart).Error
	if err != nil {
		return nil, err
	}
	return cart, nil
}

func (r *Repository) GetLastCart() ([]ds.Carts, error) {
	var cart []ds.Carts
	err := r.db.Order("date_create DESC").Find(&cart).Error
	if err != nil {
		return nil, err
	}
	return cart, nil
}

func (r *Repository) CreateCart() ([]ds.Carts, error) {
	newCart := ds.Carts{
		Status:      0,
		DateCreate:  time.Now(),
		DateUpdate:  time.Now(),
		DateFinish:  time.Now(),
		CreatorID:   1,
		ModeratorID: 2,
	}
	err := r.db.Create(&newCart).Error
	if err != nil {
		return nil, err
	}
	cart, err := r.GetLastCart()
	return cart, err
}

func (r *Repository) GetCardByID(id int) (*ds.Meals, error) { // ?
	card := &ds.Meals{}
	err := r.db.Where("id = ?", id).First(card).Error
	if err != nil {
		return nil, err
	}

	return card, nil
}

func (r *Repository) AddToCart(id int) (uint, error) {
	cart, err := r.GetWorkingCart()
	if err != nil {
		return 0, err
	}
	cart_n_meal := &ds.CartsNMeals{
		CartID:      cart[0].ID,
		ChildMealID: id,
	}
	err = r.db.Create(&cart_n_meal).Error
	if err != nil {
		return 0, err
	}
	return 0, err
}

func (r *Repository) GetMealsIDsByCartID(cartID int) ([]ds.CartsNMeals, error) {
	var MealsIDs []ds.CartsNMeals

	err := r.db.
		Where("carts_n_meals.cart_id = ?", cartID).Order("order_o ASC").Find(&MealsIDs).Error

	if err != nil {
		return nil, err
	}

	return MealsIDs, nil
}

func (r *Repository) DeleteCart(id int) error {
	err := r.db.Exec("UPDATE carts SET status = ? WHERE id = ?", 3, id).Error
	if err != nil {
		return err
	}
	return nil
}
