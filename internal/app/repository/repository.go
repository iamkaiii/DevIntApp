package repository

import (
	"DevIntApp/internal/app/ds"
	"fmt"
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
	err := r.db.Where("status=true").Order("id ASC").Find(&prods).Error
	if err != nil {
		return nil, err
	}
	return prods, nil
}

func (r *Repository) GetMealByID(mealID string) ([]ds.Meals, error) {
	var meal []ds.Meals
	err := r.db.Where("id = ?", mealID).Find(&meal).Error
	if err != nil {
		return nil, err
	}
	return meal, nil
}

func (r *Repository) GetMealByMealInfo(cardText string) ([]ds.Meals, error) {
	var milkmeal []ds.Meals
	err := r.db.Where("meal_info LIKE ?", "%"+cardText+"%").First(&milkmeal).Error
	if err != nil {
		return nil, err
	}
	return milkmeal, nil
}

func (r *Repository) GetWorkingMilkRequest() ([]ds.MilkRequests, error) {
	var milkrequest []ds.MilkRequests
	err := r.db.Where("status=0").Find(&milkrequest).Error
	if err != nil {
		return nil, err
	}
	return milkrequest, nil
}

func (r *Repository) GetLastMilkRequest() (ds.MilkRequests, error) {
	var milkrequest ds.MilkRequests
	err := r.db.Order("date_create DESC").Find(&milkrequest).Error
	if err != nil {
		return ds.MilkRequests{}, err
	}
	return milkrequest, nil
}

func (r *Repository) CreateMilkRequest() (ds.MilkRequests, error) {
	newMilkRequest := ds.MilkRequests{
		Status:     0,
		DateCreate: time.Now(),
		DateUpdate: time.Now(),
		CreatorID:  1,
	}
	err := r.db.Create(&newMilkRequest).Error
	if err != nil {
		return ds.MilkRequests{}, err
	}
	milkrequest, err := r.GetLastMilkRequest()
	return milkrequest, err
}

func (r *Repository) GetMilkRequestByID(id int) (*ds.Meals, error) { // ?
	milkrequest := &ds.Meals{}
	err := r.db.Where("id = ?", id).First(milkrequest).Error
	if err != nil {
		return nil, err
	}

	return milkrequest, nil
}

func (r *Repository) AddToMilkRequest(milreq_ID int, milkmeal_ID int) error {
	milkReqMeal := ds.MilkRequestsMeals{
		MilkRequestID: milreq_ID,
		MealID:        milkmeal_ID,
	}
	err := r.db.Create(&milkReqMeal).Error
	if err != nil {
		return fmt.Errorf("failed to add to milk request: %w", err)
	}

	return nil
}

func (r *Repository) GetMealsIDsByMilkRequestID(cartID int) ([]ds.MilkRequestsMeals, error) {
	var MealsIDs []ds.MilkRequestsMeals

	err := r.db.
		Where("milk_requests_meals.milk_request_id = ?", cartID).Order("amount ASC").Find(&MealsIDs).Error

	if err != nil {
		return nil, err
	}

	return MealsIDs, nil
}

func (r *Repository) DeleteMilkRequest(id int) error {
	err := r.db.Exec("UPDATE milk_requests SET status = ? WHERE id = ?", 3, id).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetMilkRequestStatusByID(id int) (int, error) {
	milkrequest := &ds.MilkRequests{}
	err := r.db.Where("id = ?", id).First(milkrequest).Error
	if err != nil {
		return -1, err
	}
	return milkrequest.Status, nil
}

func (r *Repository) CreateMeal(meal ds.Meals) error {
	if err := r.db.Create(&meal).Error; err != nil {
		return err
	}
	return nil
}

func (r *Repository) DeleteMealByID(id string) error {
	var meal ds.Meals
	if err := r.db.First(&meal, id).Error; err != nil {
		return err // Если запись не найдена или произошла ошибка, возвращаем её
	}
	meal.ImageUrl = ""
	meal.Status = false // Предполагается, что у вас есть поле Status в структуре Meal
	if err := r.db.Save(&meal).Error; err != nil {
		return err // Возвращаем ошибку, если обновление не удалось
	}
	return nil // Возвращаем nil, если всё прошло успешно
}

func (r *Repository) UpdateMealByID(id string, meal ds.Meals) error {
	var existingMeal ds.Meals
	if err := r.db.First(&existingMeal, id).Error; err != nil {
		return err
	}
	if err := r.db.Model(&existingMeal).Updates(meal).Error; err != nil {
		return err
	}
	return nil
}

func (r *Repository) ChangePicByID(id string, image string) error {
	// 1. Поиск записи по ID
	meal := ds.Meals{}
	result := r.db.First(&meal, "ID = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("запись с ID %s не найдена", id)
	}
	meal.ImageUrl = image
	err := r.db.Save(&meal).Error
	if err != nil {
		return err
	}
	return nil
}
