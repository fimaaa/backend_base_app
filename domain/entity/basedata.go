package entity

import "math"

type BaseReqFind struct {
	Page   int                    `json:"page" form:"page"`
	Size   int                    `json:"size" form:"size"`
	Value  interface{}            `json:"value" form:"value"`
	SortBy map[string]interface{} `json:"sort_by" form:"sort_by"`
}

type BaseResponsePagination struct {
	List       interface{}     `json:"list"`
	Pagination *PaginationData `json:"metadata"`
}

type PaginationData struct {
	CurrentPage   int  `json:"current_page"`
	PerPage       int  `json:"per_page"`
	TotalPage     int  `json:"Count"`
	TotalRecords  int  `json:"total_records"`
	LinkParameter Link `json:"link_parameter"`
	Links         Link `json:"links"`
}

type Link struct {
	First    string `json:"first"`
	Last     string `json:"last"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
}

type MyError struct {
	Message string
}

type AuthRefreshToken struct {
	Id       string `json:"id" bson:"id"`
	DeviceId string `json:"id_device" bson:"id_device"`
}

// Implement the Error method for MyError
func (e MyError) Error() string {
	return e.Message
}

// Define a function that returns a MyError
func NewMyError(message string) MyError {
	return MyError{Message: message}
}

func calculateTotalPage(totalRecords int64, perPage int64) int {
	return int(math.Ceil(float64(totalRecords) / float64(perPage)))
}
func (req BaseReqFind) ToResponse(list interface{}, totalRecords int64) BaseResponsePagination {

	return BaseResponsePagination{
		List: list,
		Pagination: &PaginationData{
			CurrentPage:  req.Page,
			PerPage:      req.Size,
			TotalPage:    calculateTotalPage(totalRecords, int64(req.Size)),
			TotalRecords: int(totalRecords),
		},
	}
}
