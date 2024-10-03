package schemas

import (
	"DevIntApp/internal/app/ds"
)

type GetAllMealsResponse struct {
	ID    int        `json:"milk_req_ID"`
	Count int        `json:"count"`
	Meals []ds.Meals `json:"meals"`
}

type GetMealResponse struct {
	Meal ds.Meals `json:"meal"`
}

type CreateMealResponse struct{}

type DeleteMealResponse struct{}

///MILK REQUESTS MILK REQUESTS MILK REQUESTS MILK REQUESTS MILK REQUESTS MILK REQUESTS MILK REQUESTS MILK REQUESTS MILK REQUESTS///

type GetAllMilkRequestsWithParamsResponse struct {
	MilkRequests []ds.MilkRequests
}

type GetAllMilkRequestsResponse struct {
	Meals []ds.MilkRequests `json:"milk_requests"`
}

type GetMilkRequestResponse struct {
	MilkRequest      ds.MilkRequests `json:"milk_requests"`
	MilkRequestMeals []ds.Meals      `json:"milk_request_meals"`
}

type DeleteMealFromMilkReqResponse struct{}

type UpdateOrderMilkReqMealsResponse struct{}
