package util

import (
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"

	"github.com/EDDYCJY/go-gin-example/pkg/setting"
)

// GetPage get page parameters
func GetPage(c *gin.Context) int {
	result := 0
	page := com.StrTo(c.Query("page")).MustInt()
	if page > 0 {
		result = (page - 1) * setting.AppSetting.PageSize
	}

	return result
}

//GetSize get size parameters
func GetSize(c *gin.Context) int {
	result := setting.AppSetting.PageSize
	size := com.StrTo(c.Query("size")).MustInt()
	if size > 10 {
		result = size
	}
	return result
}
