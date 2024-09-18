package repository

import (
	"DevIntApp/internal/app/ds"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

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

func (r *Repository) GetAllProducts() ([]ds.Products, error) { //FIO ?
	var prods []ds.Products
	err := r.db.Where("status=1").Find(&prods).Error
	if err != nil {
		return nil, err
	}
	return prods, nil
}

func (r *Repository) GetProductByID(id int) (*ds.Products, error) { // ?
	product := &ds.Products{}
	err := r.db.Where("id = ?", id).First(product).Error
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (r *Repository) GetProductByCardText(cardText string) ([]ds.Products, error) {
	var product []ds.Products
	err := r.db.Where("card_text LIKE ?", "%"+cardText+"%").First(&product).Error
	if err != nil {
		return nil, err
	}
	return product, nil
}
