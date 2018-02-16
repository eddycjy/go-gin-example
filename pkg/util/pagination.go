package util

import (
	"github.com/gin-gonic/gin"
	"github.com/Unknwon/com"

	"gin-blog/pkg/setting"
)

func GetPage(c *gin.Context) int {
	result := 0
	page, _ := com.StrTo(c.Query("page")).Int()
    if page > 0 {
        result = (page - 1) * setting.PageSize
    }

    return result
}