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

type GetAllMilkRequestsRequest struct {
}

type GetMilkRequestRequest struct {
	ID string `uri:"meal_id"`
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

type UpdateAdditionalFieldsMilkReqRequest struct {
	ID           string    `uri:"milk_request" json:"id"`
	Name         string    `json:"recipient_name"`
	Surname      string    `json:"recipient_surname"`
	Address      string    `json:"recipient_address"`
	DeliveryDate time.Time `json:"delivery_date"`
}

type DeleteMilkRequestRequest struct {
	ID int `json:"id"`
}

type CreateUserRequest struct {
	ds.Users
}

type FinishMilkRequest struct {
	ID int `json:"id"`
}
