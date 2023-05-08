package e

var MsgFlags = map[int]string{
	SUCCESS:                         "ok",
	ERROR:                           "fail",
	INVALID_PARAMS:                  "invalid params",
	ERROR_EXIST_TAG:                 "tag exists",
	ERROR_EXIST_TAG_FAIL:            "failed to check if tag exist",
	ERROR_NOT_EXIST_TAG:             "tag doesn't exist",
	ERROR_GET_TAGS_FAIL:             "failed to get tags",
	ERROR_COUNT_TAG_FAIL:            "count tag failed",
	ERROR_ADD_TAG_FAIL:              "failed to add tag",
	ERROR_EDIT_TAG_FAIL:             "failed to edit tag",
	ERROR_DELETE_TAG_FAIL:           "failed to delete tag",
	ERROR_EXPORT_TAG_FAIL:           "failed to export tag",
	ERROR_IMPORT_TAG_FAIL:           "failed to import tag",
	ERROR_NOT_EXIST_ARTICLE:         "article doesn't exist",
	ERROR_ADD_ARTICLE_FAIL:          "failed to add article",
	ERROR_DELETE_ARTICLE_FAIL:       "failed to delete article",
	ERROR_CHECK_EXIST_ARTICLE_FAIL:  "failed to check if article exist",
	ERROR_EDIT_ARTICLE_FAIL:         "failed to edit article",
	ERROR_COUNT_ARTICLE_FAIL:        "failed to count article",
	ERROR_GET_ARTICLES_FAIL:         "failed to get articles",
	ERROR_GET_ARTICLE_FAIL:          "failed to get article",
	ERROR_GEN_ARTICLE_POSTER_FAIL:   "failed to generate article poster",
	ERROR_AUTH_CHECK_TOKEN_FAIL:     "token check failed",
	ERROR_AUTH_CHECK_TOKEN_TIMEOUT:  "token check timeout",
	ERROR_AUTH_TOKEN:                "auth token error",
	ERROR_AUTH:                      "auth error",
	ERROR_UPLOAD_SAVE_IMAGE_FAIL:    "failed to save uploaded image",
	ERROR_UPLOAD_CHECK_IMAGE_FAIL:   "failed to check image",
	ERROR_UPLOAD_CHECK_IMAGE_FORMAT: "failed to to check uploaded image",
}

// GetMsg get error information based on Code
func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}

	return MsgFlags[ERROR]
}
