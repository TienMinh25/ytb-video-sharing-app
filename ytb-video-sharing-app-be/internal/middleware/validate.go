package middleware

import (
	"net/http"
	"reflect"
	"ytb-video-sharing-app-be/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func getJSONTag[T any](fieldName string) string {
	var t T
	typ := reflect.TypeOf(t)

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)

		if field.Name == fieldName {
			jsonTag := field.Tag.Get("json")

			if jsonTag != "" {
				return jsonTag
			}

			break
		}
	}

	return fieldName
}

func ValidateRequest[T any]() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req T

		if err := ctx.ShouldBindJSON(&req); err != nil {
			if validationErrors, ok := err.(validator.ValidationErrors); ok {
				errorsMap := make(map[string]string)

				for _, fieldErr := range validationErrors {
					jsonField := getJSONTag[T](fieldErr.StructField()) // get JSON field
					errorsMap[jsonField] = fieldErr.Error()
				}

				utils.ErrorResponse(ctx, http.StatusBadRequest, errorsMap)
			} else {
				utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
			}
			ctx.Abort()
			return
		}

		ctx.Set("data", req)
		ctx.Next()
	}
}
