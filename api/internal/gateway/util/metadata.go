package util

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
)

var AuthTokenType = "Bearer"

var errNotExistsAuthorizationHeader = errors.New("util: authorization header is not contain")

func GetAuthToken(c *gin.Context) (string, error) {
	authorization := c.GetHeader("Authorization")
	if authorization == "" {
		return "", errNotExistsAuthorizationHeader
	}

	token := strings.TrimPrefix(authorization, AuthTokenType+" ")
	return token, nil
}
