package schemas

import (
	_ "DevIntApp/docs"
	"mime/multipart"
	"time"
)

type GetAllMealsRequest struct {
}

type GetMealRequest struct {
	ID string
}

type CreateMealRequest struct {
	MealInfo   string `json:"meal_info"`
	MealWeight string `json:"meal_weight"`
	MealBrand  string `json:"meal_brand"`
	MealDetail string `json:"meal_detail"`
	ImageUrl   string `json:"image_url"`
	Status     bool   `json:"status"`
}

type DeleteMealRequest struct {
	ID string
}

type UpdateMealRequest struct {
	MealInfo   string `json:"meal_info"`
	MealWeight string `json:"meal_weight"`
	MealBrand  string `json:"meal_brand"`
	MealDetail string `json:"meal_detail"`
	ImageUrl   string `json:"image_url"`
	Status     bool   `json:"status"`
}

type AddMealToMilkReqRequest struct {
}

type ChangeImgRequest struct {
	ID    string
	Image *multipart.FileHeader `form:"image" json:"image"`
}

type DeleteImgRequest struct{}

///MILK REQUESTS MILK REQUESTS MILK REQUESTS MILK REQUESTS MILK REQUESTS MILK REQUESTS MILK REQUESTS MILK REQUESTS MILK REQUESTS///

type GetAllMilkRequestsWithParamsRequest struct {
	HavingStatus bool `json:"is_status"`
	Status       int  `json:"status"`
}

type GetMilkRequestRequest struct {
	ID string
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
	ID string
}

type FormMilkRequestRequest struct {
	ID string
}

type FinishMilkRequestRequest struct {
	ID           string
	Status       int       `json:"status"`
	DeliveryDate time.Time `json:"delivery_date"`
}

type DeleteMealFromMilkReqRequest struct {
	ID     string
	MealID int `json:"meal_id"`
}
type UpdateAmountMilkReqMealRequest struct {
	ID     string
	MealID int `json:"meal_id"`
	Amount int `json:"amount"`
}

type RegisterUserRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type LoginUserRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
