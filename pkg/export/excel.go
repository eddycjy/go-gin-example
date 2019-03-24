package export

import "github.com/EDDYCJY/go-gin-example/pkg/setting"

const EXT = ".xlsx"

// GetExcelFullUrl get the full access path of the Excel file
func GetExcelFullUrl(name string) string {
	return setting.AppSetting.PrefixUrl + "/" + GetExcelPath() + name
}

// GetExcelPath get the relative save path of the Excel file
func GetExcelPath() string {
	return setting.AppSetting.ExportSavePath
}

// GetExcelFullPath Get the full save path of the Excel file
func GetExcelFullPath() string {
	return setting.AppSetting.RuntimeRootPath + GetExcelPath()
}
