package util

import (
	"errors"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

var errNotFoundQuery = errors.New("this query is not found")

func GetParam(ctx *gin.Context, param string) string {
	return ctx.Param(param)
}

func GetParamInt64(ctx *gin.Context, param string) (int64, error) {
	return strconv.ParseInt(ctx.Param(param), 10, 64)
}

func GetQuery(ctx *gin.Context, query string, defaultValue string) string {
	return ctx.DefaultQuery(query, defaultValue)
}

func GetQueryInt64(ctx *gin.Context, query string, defaultValue int64) (int64, error) {
	str := strconv.FormatInt(defaultValue, 10)
	return strconv.ParseInt(ctx.DefaultQuery(query, str), 10, 64)
}

func GetQueryInt32s(ctx *gin.Context, query string) ([]int32, error) {
	str := GetQuery(ctx, query, "")
	if str == "" {
		return []int32{}, nil
	}
	strs := strings.Split(str, ",")

	res := make([]int32, len(strs))
	for i := range strs {
		num, err := strconv.ParseInt(strs[i], 10, 32)
		if err != nil {
			return nil, err
		}
		res[i] = int32(num)
	}
	return res, nil
}
