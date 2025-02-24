package dto

type ResponseSuccess[T any] struct {
	Data     T        `json:"data,omitempty"`
	Metadata Metadata `json:"metadata"`
}

type ResponseSuccessPagingation[T any] struct {
	Data     T                      `json:"data,omitempty"`
	Metadata MetadataWithPagination `json:"metadata"`
}

type ResponseError struct {
	Metadata Metadata    `json:"metadata"`
	Error    interface{} `json:"error,omitempty"`
}

type Metadata struct {
	Code int `json:"code"`
}

type MetadataWithPagination struct {
	Code       int         `json:"code"`
	Pagination *Pagination `json:"pagination"`
}

type Pagination struct {
	Page       int  `json:"page"`
	Limit      int  `json:"limit"`
	TotalItems int  `json:"total_items"`
	TotalPages int  `json:"total_pages"`
	IsNext     bool `json:"is_next"`
	IsPrevious bool `json:"is_previous"`
}

type ErrorResponse struct {
	Code    int    `json:"-"`
	Message string `json:"message"`
}

func (e *ErrorResponse) Error() string {
	return e.Message
}
