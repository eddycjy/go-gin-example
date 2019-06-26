package e

const (
	SUCCESS       = 200
	ERROR         = 500
	InvalidParams = 400

	ErrorExistTag      = 10001
	ErrorExistTagFail  = 10002
	ErrorNotExistTag   = 10003
	ErrorGetTagsFail   = 10004
	ErrorCountTagFail  = 10005
	ErrorAddTagFail    = 10006
	ErrorEditTagFail   = 10007
	ErrorDeleteTagFail = 10008
	ErrorExportTagFail = 10009
	ErrorImportTagFail = 10010

	ErrorNotExistArticle       = 10011
	ErrorCheckExistArticleFail = 10012
	ErrorAddArticleFail        = 10013
	ErrorDeleteArticleFail     = 10014
	ErrorEditArticleFail       = 10015
	ErrorCountArticleFail      = 10016
	ErrorGetArticlesFail       = 10017
	ErrorGetArticleFail        = 10018
	ErrorGenArticlePosterFail  = 10019

	ErrorAuthCheckTokenFail    = 20001
	ErrorAuthCheckTokenTimeout = 20002
	ErrorAuthToken             = 20003
	ErrorAuth                  = 20004

	ErrorUploadSaveImageFail    = 30001
	ErrorUploadCheckImageFail   = 30002
	ErrorUploadCheckImageFormat = 30003
)
