package util

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type OrderBy int32

const (
	OrderByNone OrderBy = 0 // ソートなし
	OrderByASC  OrderBy = 1 // 昇順
	OrderByDesc OrderBy = 2 // 降順
)

type Order struct {
	Key       string
	Direction OrderBy
}

func GetParam(ctx *gin.Context, param string) string {
	return ctx.Param(param)
}

func GetParamInt64(ctx *gin.Context, param string) (int64, error) {
	return strconv.ParseInt(ctx.Param(param), 10, 64)
}

func GetQuery(ctx *gin.Context, query string, defaultValue string) string {
	return ctx.DefaultQuery(query, defaultValue)
}

func GetQueryInt32(ctx *gin.Context, query string, defaultValue int32) (int32, error) {
	str := strconv.FormatInt(int64(defaultValue), 10)
	num, err := strconv.ParseInt(ctx.DefaultQuery(query, str), 10, 32)
	return int32(num), err
}

func GetQueryInt64(ctx *gin.Context, query string, defaultValue int64) (int64, error) {
	str := strconv.FormatInt(defaultValue, 10)
	return strconv.ParseInt(ctx.DefaultQuery(query, str), 10, 64)
}

func GetQueryFloat64(ctx *gin.Context, query string, defaultValue float64) (float64, error) {
	str := strconv.FormatFloat(defaultValue, 'f', -1, 64)
	return strconv.ParseFloat(ctx.DefaultQuery(query, str), 64)
}

func GetQueryStrings(ctx *gin.Context, query string) []string {
	str := GetQuery(ctx, query, "")
	if str == "" {
		return []string{}
	}
	return strings.Split(str, ",")
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

func GetOrders(ctx *gin.Context) []*Order {
	strs := GetQueryStrings(ctx, "orders")
	orders := make([]*Order, len(strs))
	for i := range strs {
		order := &Order{}
		if strings.HasPrefix(strs[i], "-") {
			order.Key = strings.TrimPrefix(strs[i], "-")
			order.Direction = OrderByDesc
		} else {
			order.Key = strs[i]
			order.Direction = OrderByASC
		}
		orders[i] = order
	}
	return orders
}
