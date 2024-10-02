package schemas

import (
	"DevIntApp/internal/app/ds"
)

type GetAllMealsResponse struct {
	ID    int        `json:"milk_req_ID"`
	Meals []ds.Meals `json:"meals"`
	Count int        `json:"count"`
}

type GetMealResponse struct {
	Meal []ds.Meals `json:"meal"`
}

type CreateMealResponse struct{}

type DeleteMealResponse struct{}

type GetAllMilkRequestsResponse struct {
	Meals []ds.MilkRequests `json:"milk_requests"`
}

type GetMilkRequestResponse struct {
	MilkRequest      *ds.MilkRequests `json:"milk_requests"`
	MilkRequestMeals []ds.Meals       `json:"milk_request_meals"`
}

type DeleteMealFromMilkReqResponse struct{}

type UpdateOrderMilkReqMealsResponse struct{}
