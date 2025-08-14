package utils

import "github.com/gin-gonic/gin"

func JSONError(c *gin.Context, status int, msg string) { c.JSON(status, gin.H{"error": msg}) }
func JSONSuccess(c *gin.Context, status int, data interface{}) { c.JSON(status, gin.H{"data": data}) }