package repository

import (
	"DevIntApp/internal/app/ds"
	"DevIntApp/internal/app/schemas"
	token2 "DevIntApp/internal/token"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"strconv"
	"time"
)

//cart status: 0 - черновик, 1 - сформирован, 2 - завершен, 3 - удалён, 4 - отклонён

type Repository struct {
	db *gorm.DB
	rd *redis.Client
}

func New(dsn string) (*Repository, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ENDPOINT"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	return &Repository{
		db: db,
		rd: redisClient,
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

func (r *Repository) GetMealByMealInfo(info string) ([]ds.Meals, error) {
	var milkMeal []ds.Meals
	err := r.db.Where("meal_info LIKE ?", "%"+info+"%").First(&milkMeal).Error
	if err != nil {
		return nil, err
	}
	return milkMeal, nil
}

func (r *Repository) GetWorkingMilkRequest() ([]ds.MilkRequests, error) {
	var milkRequest []ds.MilkRequests
	err := r.db.Where("status=0").Find(&milkRequest).Error
	if err != nil {
		return nil, err
	}
	return milkRequest, nil
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
	creatorID := 1
	newMilkRequest := ds.MilkRequests{
		Status:     0,
		DateCreate: time.Now(),
		DateUpdate: time.Now(),
		CreatorID:  &creatorID,
	}
	err := r.db.Create(&newMilkRequest).Error
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

func (r *Repository) AddToMilkRequest(milkRequestID int, milkMealID int) error {
	milkReqMeal := ds.MilkRequestsMeals{
		MilkRequestID: milkRequestID,
		MealID:        milkMealID,
	}
	err := r.db.Create(&milkReqMeal).Error
	if err != nil {
		return fmt.Errorf("failed to add to milk request: %w", err)
	}

	return nil
}

func (r *Repository) GetMealsIDsByMilkRequestID(milkRequestID int) ([]int, error) {
	var MealsIds []int
	err := r.db.
		Model(&ds.MilkRequestsMeals{}).
		Where("milk_request_id = ?", milkRequestID).
		Pluck("meal_id", &MealsIds).
		Error

	if err != nil {
		return nil, err
	}
	return MealsIds, nil
}

func (r *Repository) DeleteMilkRequest(id int) error {
	var milkRequest ds.MilkRequests
	if err := r.db.First(&milkRequest, "id = ?", id).Error; err != nil {
		return err
	}
	milkRequest.Status = 3 // Устанавливаем статус удаления
	if err := r.db.Save(&milkRequest).Error; err != nil {
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

func (r *Repository) UpdateFieldsMilkReq(request schemas.UpdateFieldsMilkReqRequest) error {
	var milkRequest ds.MilkRequests
	// Загрузка записи из базы данных по ID
	if err := r.db.First(&milkRequest, "id = ?", request.ID).Error; err != nil {
		return err
	}
	if request.Name != "" {
		milkRequest.RecipientName = request.Name
	}
	if request.Surname != "" {
		milkRequest.RecipientSurname = request.Surname
	}
	if request.Address != "" {
		milkRequest.Address = request.Address
	}
	if !request.DeliveryDate.IsZero() {
		milkRequest.DeliveryDate = request.DeliveryDate
	}
	if err := r.db.Save(&milkRequest).Error; err != nil {
		return err
	}
	return nil // Возвращаем nil, если все прошло успешно
}

func (r *Repository) FormMilkRequest(id string) error {
	var milkRequest ds.MilkRequests
	if err := r.db.First(&milkRequest, "id = ?", id).Error; err != nil {
		return err
	}
	if milkRequest.CreatorID == nil {
		err := fmt.Errorf("Unable to finish request. Probably some fields are empty")
		return err
	}
	milkRequest.Status = 1
	if err := r.db.Save(&milkRequest).Error; err != nil {
		return err
	}
	return nil
}

func (r *Repository) FinishMilkRequest(id string, status int, deliveryDate time.Time) error {
	var milkRequest ds.MilkRequests
	if err := r.db.First(&milkRequest, "id = ?", id).Error; err != nil {
		return err
	}
	if milkRequest.CreatorID == nil {
		err := fmt.Errorf("Unable to finish request. Probably some fields are empty")
		return err
	}
	moderatorID := 2
	milkRequest.Status = status
	milkRequest.DeliveryDate = deliveryDate
	milkRequest.DateFinish = time.Now()
	milkRequest.ModeratorID = &moderatorID
	if err := r.db.Save(&milkRequest).Error; err != nil {
		return err
	}
	return nil
}

func (r *Repository) DeleteMealFromMilkRequest(id string, mealID int) error {
	var milkRequestMeal ds.MilkRequestsMeals
	if err := r.db.Where("milk_request_id = ? AND meal_id = ?", id, mealID).First(&milkRequestMeal).Error; err != nil {
		return err
	}
	if err := r.db.Delete(&milkRequestMeal).Error; err != nil {
		return err
	}
	return nil
}

func (r *Repository) UpdateAmountMilkReqMeal(id string, mealID int, amount int) error {
	var milkRequestMeal ds.MilkRequestsMeals
	if err := r.db.Where("milk_request_id = ? AND meal_id = ?", id, mealID).First(&milkRequestMeal).Error; err != nil {
		return err
	}
	milkRequestMeal.Amount = amount
	if err := r.db.Save(&milkRequestMeal).Error; err != nil {
		return err
	}
	return nil
}

func (r *Repository) RegisterUser(user ds.Users) (error, int) {
	err := r.db.First(&user, "login = ?", user.Login).Error
	if err == nil {
		return err, 0 //пользователь с таким логином уже есть
	}
	if err = r.db.Create(&user).Error; err != nil {
		return err, 1 //пользователь создан
	}
	return nil, 2 // произошла ошибка с создание пользователя
}

func (r *Repository) LoginUser(user ds.Users) (error, string) {
	var user_in_db ds.Users
	err := r.db.First(&user_in_db, "login = ?", user.Login).Error
	if err != nil {
		return errors.New("Такой пользователь не существует"), ""
	}
	if user.Password != user_in_db.Password {
		return errors.New("Пароли не совпадают"), ""
	}

	token, err := token2.GenerateJWTToken(user_in_db)
	if err != nil {
		return err, ""
	}

	err = r.SaveJWTToken(user_in_db.ID, token)
	if err != nil {
		return err, ""
	}
	return nil, token

}

func (r *Repository) SaveJWTToken(userID int, token string) error {
	exp := 1 * time.Hour
	userID_str := strconv.Itoa(userID)
	err := r.rd.Set(userID_str, token, exp).Err()
	if err != nil {
		return err
	}
	return nil
}
