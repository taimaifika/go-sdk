package common

type successPaginationRes struct {
	Data   interface{} `json:"data"`
	Paging interface{} `json:"paging,omitempty"`
	Filter interface{} `json:"filter,omitempty"`
}

func NewSuccessPaginationResponse(data, paging, filter interface{}) *successPaginationRes {
	return &successPaginationRes{Data: data, Paging: paging, Filter: filter}
}

func SimpleSuccessPaginationResponse(data interface{}) *successPaginationRes {
	return NewSuccessPaginationResponse(data, nil, nil)
}

func SimpleSuccessResponse(data interface{}) interface{} {
	return data
}
