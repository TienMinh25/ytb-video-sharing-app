package utils

import (
	"net/http"
	"ytb-video-sharing-app-be/internal/dto"

	"github.com/gin-gonic/gin"
)

func SuccessResponse[T any](ctx *gin.Context, statusCode int, data T) {
	ctx.JSON(statusCode, dto.ResponseSuccess[T]{
		Data: data,
		Metadata: dto.Metadata{
			Code: statusCode,
		},
	})
}

func PaginatedResponse[T any](ctx *gin.Context, data T, currentPage, limit, totalPages, totalItems int, isNext, isPrevious bool) {
	ctx.JSON(http.StatusOK, dto.ResponseSuccessPagingation[T]{
		Data: data,
		Metadata: dto.MetadataWithPagination{
			Code: http.StatusOK,
			Pagination: &dto.Pagination{
				Page:       currentPage,
				Limit:      limit,
				TotalItems: totalItems,
				TotalPages: totalPages,
				IsNext:     isNext,
				IsPrevious: isPrevious,
			},
		},
	})
}

func ErrorResponse(ctx *gin.Context, statusCode int, errDetail interface{}) {
	ctx.JSON(statusCode, dto.ResponseError{
		Metadata: dto.Metadata{
			Code: statusCode,
		},
		Error: errDetail,
	})
}
