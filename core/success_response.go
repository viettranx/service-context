package core

type successResponse struct {
	Data   interface{} `json:"data"`
	Paging interface{} `json:"paging,omitempty"`
	Extra  interface{} `json:"extra,omitempty"`
}

func SuccessResponse(data, paging, extra interface{}) *successResponse {
	return &successResponse{Data: data, Paging: paging, Extra: extra}
}

func ResponseData(data interface{}) *successResponse {
	return SuccessResponse(data, nil, nil)
}
