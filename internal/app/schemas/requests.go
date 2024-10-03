package schemas

import (
	"DevIntApp/internal/app/ds"
	"time"
)

type GetAllMealsRequest struct{}

type GetMealRequest struct {
	ID string
}

type CreateMealRequest struct {
	ds.Meals
}

type DeleteMealRequest struct {
	ID string
}

type UpdateMealRequest struct {
	ID string
	ds.Meals
}

type AddMealToMilkReqRequest struct {
	ID string
}

type ChangePicRequest struct {
	ID       string `json:"id"`
	ImageUrl string `json:"image_link"`
}

///MILK REQUESTS MILK REQUESTS MILK REQUESTS MILK REQUESTS MILK REQUESTS MILK REQUESTS MILK REQUESTS MILK REQUESTS MILK REQUESTS///

type GetAllMilkRequestsWithParamsRequest struct {
	HavingStatus bool      `json:"is_status"`
	Status       int       `json:"status"`
	FromDate     time.Time `json:"from_date"`
	ToDate       time.Time `json:"to_date"`
}

type GetMilkRequestRequest struct {
	ID string
}

type DeleteMealFromMilkReqRequest struct {
	ID     int `json:"milk_req_id"`
	MealID int `json:"meal_id"`
}

type UpdateOrderMilkReqMealsRequest struct {
	ID     int `json:"milk_req_id"`
	MealID int `json:"meal_id"`
	OrderO int `json:"order_o"`
}

type UpdateFieldsMilkReqRequest struct {
	ID           string    `uri:"milk_request" json:"id"`
	Name         string    `json:"recipient_name"`
	Surname      string    `json:"recipient_surname"`
	Address      string    `json:"recipient_address"`
	DeliveryDate time.Time `json:"delivery_date"`
}

type DeleteMilkRequestRequest struct {
	ID int
}

type FormMilkRequestRequest struct {
	ID int
}

type FinishMilkRequestRequest struct {
	ID           int
	Status       int       `json:"status"`
	DeliveryDate time.Time `json:"delivery_date"`
}

type CreateUserRequest struct {
	ds.Users
}
