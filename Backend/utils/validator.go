package utils

import "github.com/gin-gonic/gin"

func ValidateJSON(c *gin.Context, obj interface{}) bool {

	if err := c.ShouldBindJSON(obj); err != nil {

		BadRequest(c, "Request tidak valid", err.Error())

		return false
	}

	return true
}