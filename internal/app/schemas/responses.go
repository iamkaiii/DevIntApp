package schemas

import (
	_ "DevIntApp/docs"
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

type CreateMealResponse struct {
	ID              int
	MessageResponse string
}

type DeleteMealResponse struct {
	ID              int
	MessageResponse string
}

type UpdateMealResponse struct {
	ID              int
	MessageResponse string
}

///MILK REQUESTS MILK REQUESTS MILK REQUESTS MILK REQUESTS MILK REQUESTS MILK REQUESTS MILK REQUESTS MILK REQUESTS MILK REQUESTS///

type GetAllMilkRequestsWithParamsResponse struct {
	MilkRequests []ds.MilkRequests
}

type GetAllMilkRequestsResponse struct {
	Meals []ds.MilkRequests `json:"milk_requests"`
}

type GetMilkRequestResponse struct {
	MilkRequest      ds.MilkRequests `json:"milk_requests"`
	Count            int             `json:"count"`
	MilkRequestMeals []ds.Meals      `json:"milk_request_meals"`
}

type DeleteMealFromMilkReqResponse struct{}

type UpdateOrderMilkReqMealsResponse struct{}

type ResponseMessage struct {
}
