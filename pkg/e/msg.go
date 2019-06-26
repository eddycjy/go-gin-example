package e

var MsgFlags = map[int]string{
	SUCCESS:                     "ok",
	ERROR:                       "fail",
	InvalidParams:               "请求参数错误",
	ErrorExistTag:               "已存在该标签名称",
	ErrorExistTagFail:           "获取已存在标签失败",
	ErrorNotExistTag:            "该标签不存在",
	ErrorGetTagsFail:            "获取所有标签失败",
	ErrorCountTagFail:           "统计标签失败",
	ErrorAddTagFail:             "新增标签失败",
	ErrorEditTagFail:            "修改标签失败",
	ErrorDeleteTagFail:          "删除标签失败",
	ErrorExportTagFail:          "导出标签失败",
	ErrorImportTagFail:          "导入标签失败",
	ErrorNotExistArticle:        "该文章不存在",
	ErrorAddArticleFail:         "新增文章失败",
	ErrorDeleteArticleFail:      "删除文章失败",
	ErrorCheckExistArticleFail:  "检查文章是否存在失败",
	ErrorEditArticleFail:        "修改文章失败",
	ErrorCountArticleFail:       "统计文章失败",
	ErrorGetArticlesFail:        "获取多个文章失败",
	ErrorGetArticleFail:         "获取单个文章失败",
	ErrorGenArticlePosterFail:   "生成文章海报失败",
	ErrorAuthCheckTokenFail:     "Token鉴权失败",
	ErrorAuthCheckTokenTimeout:  "Token已超时",
	ErrorAuthToken:              "Token生成失败",
	ErrorAuth:                   "Token错误",
	ErrorUploadSaveImageFail:    "保存图片失败",
	ErrorUploadCheckImageFail:   "检查图片失败",
	ErrorUploadCheckImageFormat: "校验图片错误，图片格式或大小有问题",
}

// GetMsg get error information based on Code
func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}

	return MsgFlags[ERROR]
}
