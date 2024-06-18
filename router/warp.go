package router

import (
	"encoding/json"
	"fastApi/app/http/response"
	"fastApi/core/global"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
)

type HandlerFunc func(ctx *gin.Context) (res interface{}, err error)

func wrap(handler HandlerFunc) func(c *gin.Context) {
	return func(ctx *gin.Context) {
		res, err := handler(ctx)
		if err != nil {
			ctx.JSON(http.StatusOK, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusOK, res)
	}
}

// errorResponse 返回错误消息
func errorResponse(err error) response.Response {
	ve, ok := err.(validator.ValidationErrors)

	if ok {
		return response.ParamErr("参数错误", validationErrorsFormat(ve.Translate(global.Trans)), err)
	}

	if _, ok := err.(*json.UnmarshalTypeError); ok {
		return response.ParamErr("JSON类型不匹配", nil, err)
	}

	return response.ParamErr("参数错误", nil, err)
}

func validationErrorsFormat(fields map[string]string) map[string][]string {
	res := map[string][]string{}
	var errs []string
	for _, err := range fields {
		errs = append(errs, err)
	}

	res["errors"] = errs
	return res
}
