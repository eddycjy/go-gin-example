package v1

import (
	"github.com/EDDYCJY/go-gin-example/pkg/app"
	"github.com/EDDYCJY/go-gin-example/pkg/e"
	"github.com/EDDYCJY/go-gin-example/pkg/util"
	"github.com/EDDYCJY/go-gin-example/service/operation_service"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"net/http"
	"strings"
)

func ListOperationInfos(c *gin.Context) {
	var isSelected = -1
	appG := app.Gin{C: c}
	programName := c.Query("programName")
	programsetName := c.Query("programsetName")
	programId := c.Query("programId")
	programsetId := c.Query("programsetId")
	if c.Query("isSelected") != "" {
		isSelected = com.StrTo(c.Query("isSelected")).MustInt()
	}

	infoService := operation_service.Info{
		ProgramName:    programName,
		ProgramsetName: programsetName,
		ProgramId:      programId,
		ProgramsetId:   programsetId,
		IsSelected:     isSelected,
		PageNum:        util.GetPage(c),
		PageSize:       util.GetSize(c),
	}
	infos, err := infoService.GetAll()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_TAGS_FAIL, nil)
		return
	}

	count, err := infoService.Count()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_COUNT_TAG_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, map[string]interface{}{
		"lists": infos,
		"total": count,
	})
}

func ListSelectedOperationInfos(c *gin.Context) {
	appG := app.Gin{C: c}

	var isValid = -1

	programName := c.Query("programName")
	programsetName := c.Query("programsetName")
	programId := c.Query("programId")
	programsetId := c.Query("programsetId")
	if c.Query("isValid") != "" {
		isValid = com.StrTo(c.Query("isValid")).MustInt()
	}

	infoService := operation_service.SelectedInfo{
		ProgramName:    programName,
		ProgramsetName: programsetName,
		ProgramId:      programId,
		ProgramsetId:   programsetId,
		IsValid:        isValid,
		PageNum:        util.GetPage(c),
		PageSize:       util.GetSize(c),
	}
	infos, err := infoService.GetSelectedAll()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_TAGS_FAIL, nil)
		return
	}

	count, err := infoService.SelectedCount()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_COUNT_TAG_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, map[string]interface{}{
		"lists": infos,
		"total": count,
	})
}

func ListValidSelectedOperationInfos(c *gin.Context) {
	appG := app.Gin{C: c}
	infoService := operation_service.SelectedInfo{}
	infos, err := infoService.GetSelectedAllValid()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_TAGS_FAIL, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, map[string]interface{}{
		"lists": infos,
	})
}

func AddSelectedOperationInfos(c *gin.Context) {
	appG := app.Gin{C: c}
	var selectedInfos []*operation_service.SelectedInfo

	if err := c.ShouldBindJSON(&selectedInfos); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_TAGS_FAIL, nil)
		return
	}

	for _, selectedInfo := range selectedInfos {
		selectedInfoService := operation_service.SelectedInfo{
			ProgramId:      selectedInfo.ProgramId,
			ProgramsetId:   selectedInfo.ProgramsetId,
			ProgramsetName: selectedInfo.ProgramsetName,
			ProgramName:    selectedInfo.ProgramName,
			PlayUrl:        selectedInfo.PlayUrl,
			CacheSize:      selectedInfo.CacheSize,
			CacheValidity:  selectedInfo.CacheValidity,
			VideoType:      selectedInfo.VideoType,
		}

		if err := selectedInfoService.AddSelectedInfo(); err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR_ADD_ARTICLE_FAIL, nil)
			return
		}

	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

func EditSelectedOperationInfos(c *gin.Context) {
	appG := app.Gin{C: c}

	var selectedInfos []*operation_service.SelectedInfo

	if err := c.ShouldBindJSON(&selectedInfos); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_TAGS_FAIL, nil)
		return
	}

	for _, selectedInfo := range selectedInfos {
		selectedInfoService := operation_service.SelectedInfo{
			ID:            selectedInfo.ID,
			ProgramId:     selectedInfo.ProgramId,
			PlayCount:     selectedInfo.PlayCount,
			SaveBandwidth: selectedInfo.SaveBandwidth,
		}

		if err := selectedInfoService.EditSelectedInfo(); err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR_ADD_ARTICLE_FAIL, nil)
			return
		}
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

func DeleteSelectedOperationInfos(c *gin.Context) {
	appG := app.Gin{C: c}
	ids := c.Query("ids")
	idArr := strings.Split(ids, ",")
	for _, id := range idArr {
		selectedInfoService := operation_service.SelectedInfo{ID: id}
		if err := selectedInfoService.DeleteSelectedInfo(); err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR_ADD_ARTICLE_FAIL, nil)
			return
		}
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}
