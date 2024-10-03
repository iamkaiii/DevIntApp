package repository

import (
	"DevIntApp/internal/app/ds"
	"fmt"
	log "github.com/sirupsen/logrus"
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

func (r *Repository) GetMealByID(mealID string) (ds.Meals, error) {
	var meal ds.Meals
	err := r.db.Where("id = ?", mealID).Find(&meal).Error
	if err != nil {
		return ds.Meals{}, err
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
		Status:      0,
		DateCreate:  time.Now(),
		DateUpdate:  time.Now(),
		CreatorID:   1,
		ModeratorID: 2,
	}
	err := r.db.Create(&newMilkRequest).Error
	if err != nil {
		return ds.MilkRequests{}, err
	}
	err = r.db.Model(&newMilkRequest).Update("moderator_id", nil).Error
	if err != nil {
		return ds.MilkRequests{}, err
	}

	milkrequest, err := r.GetLastMilkRequest()
	return milkrequest, err
}

func (r *Repository) GetMilkRequestByID(id int) (ds.MilkRequests, error) { // ?
	milkrequest := ds.MilkRequests{}
	err := r.db.Where("id = ?", id).First(&milkrequest).Error
	if err != nil {
		return ds.MilkRequests{}, err
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

func (r *Repository) GetMealsIDsByMilkRequestID(milk_request_ID int) ([]int, error) {
	var MealsIds []int
	err := r.db.
		Model(&ds.MilkRequestsMeals{}).
		Where("milk_request_id = ?", milk_request_ID).
		Pluck("meal_id", &MealsIds).
		Error

	if err != nil {
		return nil, err
	}
	return MealsIds, nil
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
	if err := r.db.First(&existingMeal, "ID = ?", id).Error; err != nil {
		return err
	}

	existingMeal.MealInfo = meal.MealInfo
	existingMeal.MealWeight = meal.MealWeight
	existingMeal.MealBrand = meal.MealBrand
	existingMeal.MealDetail = meal.MealDetail
	existingMeal.ImageUrl = meal.ImageUrl
	existingMeal.Status = meal.Status

	err := r.db.Save(&existingMeal).Error
	if err != nil {
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

func (r *Repository) GetAllMilkRequestsWithFilters(status int, having_status bool) ([]ds.MilkRequests, error) {
	var milkRequests []ds.MilkRequests
	log.Println(status, having_status)
	db := r.db // Инициализируем db без фильтра по дате
	if having_status {
		db = db.Where("Status = ?", status) // Фильтр по статусу
	}
	err := db.Find(&milkRequests).Error // Выборка записей из базы данных
	if err != nil {
		return nil, fmt.Errorf("failed to get milk requests: %w", err)
	}
	return milkRequests, nil
}
