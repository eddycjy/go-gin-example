package models

import (
	"github.com/EDDYCJY/go-gin-example/pkg/logging"
	"github.com/EDDYCJY/go-gin-example/pkg/util"
	"github.com/jinzhu/gorm"
	"time"
)

type Info struct {
	ProgramId        string  `gorm:"primary_key" json:"programId"`
	ProgramsetId     string  `json:"programsetId"`
	ProgramsetName   string  `json:"programsetName"`
	ProgramClassCode string  `json:"programClassCode"`
	StatusCode       string  `json:"statusCode"`
	Score            float32 `json:"score"`
	Status           string  `json:"status"`
	ProgramName      string  `json:"programName"`
	PlayUrl          string  `json:"playUrl"`
	VideoType        string  `json:"videoType"`
	SelectedId       string  `json:"selectedId"`
}

type InfoRes struct {
	Info
	IsSelected        int    `json:"isSelected,omitempty"`
	CacheValidityTime string `json:"cacheValidityTime"`
}

type SelectedInfo struct {
	ID             string `gorm:"primary_key" json:"id"`
	ProgramId      string `json:"programId"`
	ProgramsetId   string `json:"programsetId"`
	ProgramsetName string `json:"programsetName"`
	//ProgramClassCode string  `json:"programClassCode"`
	//StatusCode       string  `json:"statusCode"`
	//Score            float32 `json:"score"`
	//Status           string  `json:"status"`
	ProgramName       string `json:"programName"`
	PlayUrl           string `json:"playUrl"`
	CacheSize         int    `json:"cacheSize"`
	CacheValidity     int    `json:"cacheValidity"`
	CacheValidityTime string `json:"cacheValidityTime"`
	TerminalCount     int    `json:"terminalCount"`
	PlayCount         int    `json:"playCount"`
	SaveBandwidth     int    `json:"saveBandwidth"`
	VideoType         string `json:"videoType"`
}

type OuterSelectedInfo struct {
	ID            string `json:"id"`
	ProgramId     string `json:"programId"`
	ProgramName   string `json:"programName"`
	PlayUrl       string `json:"playUrl"`
	CacheSize     int    `json:"cacheSize"`
	CacheValidity int    `json:"cacheValidity"`
	VideoType     string `json:"videoType"`
}

// GetInfos gets a list of infos based on paging and constraints
func GetInfos(pageNum int, pageSize int, isSelected int, maps map[string]string) ([]*InfoRes, error) {
	var (
		infos   []*InfoRes
		infoRes []*InfoRes
		dbres   *gorm.DB
		err     error
	)

	currentTimestamp := time.Now().Unix()
	datetime := time.Unix(int64(currentTimestamp), 0).Format("2006-01-02 15:04:05")

	if pageSize > 0 && pageNum >= 0 {
		if isSelected == 0 {
			dbres = db.Table("operation_info").Select("operation_info.*, selected_operation_info.cache_validity_time").Joins("left join selected_operation_info on operation_info.selected_id = selected_operation_info.id ").Where("operation_info.program_name like ? AND operation_info.programset_name like ? AND operation_info.program_id like ? AND operation_info.programset_id like ? AND selected_operation_info.cache_validity_time < ?",
				maps["program_name"], maps["programset_name"], maps["program_id"], maps["programset_id"], datetime).Or("operation_info.program_name like ? AND operation_info.programset_name like ? AND operation_info.program_id like ? AND operation_info.programset_id like ? AND selected_operation_info.cache_validity_time IS NULL", maps["program_name"], maps["programset_name"], maps["program_id"], maps["programset_id"]).Offset(pageNum).Limit(pageSize).Find(&infoRes)
		} else if isSelected == 1 {
			dbres = db.Table("operation_info").Select("operation_info.*, selected_operation_info.cache_validity_time").Joins("left join selected_operation_info on operation_info.selected_id = selected_operation_info.id ").Where(" operation_info.program_name like ? AND operation_info.programset_name like ? AND operation_info.program_id like ? AND operation_info.programset_id like ? AND selected_operation_info.cache_validity_time > ?", maps["program_name"], maps["programset_name"], maps["program_id"], maps["programset_id"], datetime).Offset(pageNum).Limit(pageSize).Find(&infoRes)
		} else {
			dbres = db.Table("operation_info").Select("operation_info.*, selected_operation_info.cache_validity_time").Joins("left join selected_operation_info on operation_info.selected_id = selected_operation_info.id ").Where(" operation_info.program_name like ? AND operation_info.programset_name like ? AND operation_info.program_id like ? AND operation_info.programset_id like ? ", maps["program_name"], maps["programset_name"], maps["program_id"], maps["programset_id"]).Offset(pageNum).Limit(pageSize).Find(&infoRes)
		}
		err = dbres.Error
	} else {
		err = db.Table("operation_info").Where("").Find(&infos).Error
	}

	for _, info := range infoRes {
		timestamp, err := util.FormatDatetimeToTimestamp(info.CacheValidityTime)
		if err != nil {
			logging.Error(err)
		} else {
			if timestamp > currentTimestamp {
				info.IsSelected = 1
			} else {
				info.IsSelected = 0
			}
		}
		infos = append(infos, info)
	}

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return infos, nil
}

// GetInfoTotal counts the total number of infos based on the constraint
func GetInfoTotal(isSelected int, maps map[string]string) (int, error) {
	var (
		count int
		dbres *gorm.DB
		err   error
	)

	currentTimestamp := time.Now().Unix()
	datetime := time.Unix(int64(currentTimestamp), 0).Format("2006-01-02 15:04:05")

	if isSelected == 0 {
		dbres = db.Table("operation_info").
			Select("operation_info.*, selected_operation_info.cache_validity_time").
			Joins("left join selected_operation_info on operation_info.selected_id = selected_operation_info.id ").
			Where("operation_info.program_name like ? AND operation_info.programset_name like ? AND operation_info.program_id like ? AND operation_info.programset_id like ? AND selected_operation_info.cache_validity_time < ?",
				maps["program_name"], maps["programset_name"], maps["program_id"], maps["programset_id"], datetime).
			Or("operation_info.program_name like ? AND operation_info.programset_name like ? AND operation_info.program_id like ? AND operation_info.programset_id like ? AND selected_operation_info.cache_validity_time IS NULL",
				maps["program_name"], maps["programset_name"], maps["program_id"], maps["programset_id"]).
			Count(&count)

	} else if isSelected == 1 {
		dbres = db.Table("operation_info").
			Select("operation_info.*, selected_operation_info.cache_validity_time").
			Joins("left join selected_operation_info on operation_info.selected_id = selected_operation_info.id ").
			Where(" operation_info.program_name like ? AND operation_info.programset_name like ? AND operation_info.program_id like ? AND operation_info.programset_id like ? AND selected_operation_info.cache_validity_time > ?",
				maps["program_name"], maps["programset_name"], maps["program_id"], maps["programset_id"], datetime).
			Count(&count)
	} else {
		dbres = db.Table("operation_info").
			Select("operation_info.*, selected_operation_info.cache_validity_time").
			Joins("left join selected_operation_info on operation_info.selected_id = selected_operation_info.id ").
			Where(" operation_info.program_name like ? AND operation_info.programset_name like ? AND operation_info.program_id like ? AND operation_info.programset_id like ? ",
				maps["program_name"], maps["programset_name"], maps["program_id"], maps["programset_id"]).
			Count(&count)
	}
	err = dbres.Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return 0, err
	}

	return count, nil
}

// GetSelectedInfos gets a list of infos based on paging and constraints
func GetSelectedInfos(pageNum int, pageSize int, isValid int, maps map[string]string) ([]*SelectedInfo, error) {
	var (
		infos []*SelectedInfo
		err   error
	)

	currentTimestamp := time.Now().Unix()
	datetime := time.Unix(int64(currentTimestamp), 0).Format("2006-01-02 15:04:05")

	if pageSize > 0 && pageNum >= 0 {
		if isValid == 0 {
			err = db.Table("selected_operation_info").
				Where("program_name like ? AND programset_name like ? AND program_id like ? AND programset_id like ? AND cache_validity_time < ?",
					maps["program_name"], maps["programset_name"], maps["program_id"], maps["programset_id"], datetime).
				Offset(pageNum).Limit(pageSize).Find(&infos).Error
		} else if isValid == 1 {
			err = db.Table("selected_operation_info").
				Where("program_name like ? AND programset_name like ? AND program_id like ? AND programset_id like ? AND cache_validity_time > ?", maps["program_name"], maps["programset_name"], maps["program_id"], maps["programset_id"], datetime).Offset(pageNum).Limit(pageSize).Find(&infos).Error
		} else {
			err = db.Table("selected_operation_info").
				Where("program_name like ? AND programset_name like ? AND program_id like ? AND programset_id like ? ",
					maps["program_name"], maps["programset_name"], maps["program_id"], maps["programset_id"]).
				Offset(pageNum).Limit(pageSize).Find(&infos).Error
		}
	} else {
		err = db.Table("selected_operation_info").Where("").Find(&infos).Error
	}

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return infos, nil
}

// GetSelectedInfoTotal counts the total number of infos based on the constraint
func GetSelectedInfoTotal(isValid int, maps map[string]string) (int, error) {
	var (
		count int
		err   error
	)

	currentTimestamp := time.Now().Unix()
	datetime := time.Unix(int64(currentTimestamp), 0).Format("2006-01-02 15:04:05")

	if isValid == 0 {
		err = db.Table("selected_operation_info").
			Where("program_name like ? AND programset_name like ? AND program_id like ? AND programset_id like ? AND cache_validity_time < ?",
				maps["program_name"], maps["programset_name"], maps["program_id"], maps["programset_id"], datetime).
			Count(&count).Error
	} else if isValid == 1 {
		err = db.Table("selected_operation_info").
			Where("program_name like ? AND programset_name like ? AND program_id like ? AND programset_id like ? AND cache_validity_time > ?", maps["program_name"], maps["programset_name"], maps["program_id"], maps["programset_id"], datetime).Count(&count).Error
	} else {
		err = db.Table("selected_operation_info").
			Where("program_name like ? AND programset_name like ? AND program_id like ? AND programset_id like ?", maps["program_name"], maps["programset_name"], maps["program_id"], maps["programset_id"]).Count(&count).Error
	}

	if err != nil && err != gorm.ErrRecordNotFound {
		return 0, err
	}

	return count, nil
}

func GetValidSelectedInfos() ([]*OuterSelectedInfo, error) {
	var (
		infos []*OuterSelectedInfo
		err   error
	)

	timestamp := time.Now().Unix()
	datetime := time.Unix(int64(timestamp), 0).Format("2006-01-02 15:04:05")

	err = db.Table("selected_operation_info").Where("cache_validity_time > ?", datetime).Find(&infos).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return infos, nil
}

func GetSelectInfoById(programId string, recordId string) (*SelectedInfo, error) {

	var (
		info SelectedInfo
		err  error
	)

	//timestamp := time.Now().Unix()
	//datetime := time.Unix(int64(timestamp), 0).Format("2006-01-02 15:04:05")
	err = db.Table("selected_operation_info").Where("program_id = ? AND id = ?", programId, recordId).First(&info).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return &info, err
}

func AddSelectedInfo(data map[string]interface{}) error {

	day := data["cacheValidity"].(int)
	timestamp := time.Now().Unix() + int64(day*24*3600)
	datetime := time.Unix(timestamp, 0).Format("2006-01-02 15:04:05")

	selectedInfo := SelectedInfo{
		ID:                data["id"].(string),
		ProgramId:         data["programId"].(string),
		ProgramsetId:      data["programsetId"].(string),
		ProgramsetName:    data["programsetName"].(string),
		ProgramName:       data["programName"].(string),
		PlayUrl:           data["playUrl"].(string),
		CacheSize:         data["cacheSize"].(int),
		CacheValidity:     data["cacheValidity"].(int),
		CacheValidityTime: datetime,
		TerminalCount:     0,
		PlayCount:         0,
		SaveBandwidth:     0,
		VideoType:         data["videoType"].(string),
	}

	dbres := db.Table("selected_operation_info").Create(&selectedInfo)

	if dbres.Error != nil {
		return dbres.Error
	}

	if err := db.Table("operation_info").Model(&Info{}).Where("program_id = ?", selectedInfo.ProgramId).Updates(map[string]interface{}{
		"selected_id": selectedInfo.ID,
	}).Error; err != nil {
		return err
	}
	return nil
}

func EditSelectedInfo(programId string, recordId string, data map[string]interface{}) error {
	if err := db.Table("selected_operation_info").Model(&SelectedInfo{}).Where("program_id = ? AND id = ?", programId, recordId).Update(data).Error; err != nil {
		return err
	}

	return nil
}

func DeleteSelectedInfo(id string) error {
	if err := db.Table("selected_operation_info").Where("id = ?", id).Delete(&SelectedInfo{}).Error; err != nil {
		return err
	}

	return nil
}
