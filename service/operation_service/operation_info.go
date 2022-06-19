package operation_service

import (
	"fmt"
	"github.com/EDDYCJY/go-gin-example/models"
	"github.com/EDDYCJY/go-gin-example/pkg/logging"
	"github.com/EDDYCJY/go-gin-example/pkg/util"
)

type Info struct {
	ProgramId      string
	ProgramsetId   string
	ProgramName    string
	ProgramsetName string

	IsSelected int
	PageNum    int
	PageSize   int
}

type SelectedInfo struct {
	ID             string
	ProgramId      string
	ProgramsetId   string
	ProgramsetName string
	ProgramName    string
	PlayUrl        string
	VideoType      string

	TerminalCount int
	PlayCount     int
	SaveBandwidth int

	IsValid       int
	CacheSize     int
	CacheValidity int

	PageNum  int
	PageSize int
}

type SelectedInfoPost struct {
	UserID string
	Count  int
}

func (t *Info) Count() (int, error) {
	return models.GetInfoTotal(t.IsSelected, t.getMaps())
}

func (t *Info) GetAll() ([]*models.InfoRes, error) {
	var (
		infos []*models.InfoRes
	)

	infos, err := models.GetInfos(t.PageNum, t.PageSize, t.IsSelected, t.getMaps())
	if err != nil {
		return nil, err
	}

	return infos, nil
}

func (t *Info) getMaps() map[string]string {
	maps := make(map[string]string)

	maps["programset_name"] = "%" + t.ProgramsetName + "%"
	maps["program_name"] = "%" + t.ProgramName + "%"
	maps["programset_id"] = "%" + t.ProgramsetId + "%"
	maps["program_id"] = "%" + t.ProgramId + "%"

	return maps
}

func (st *SelectedInfo) SelectedCount() (int, error) {
	return models.GetSelectedInfoTotal(st.IsValid, st.getMaps())
}

func (st *SelectedInfo) GetSelectedAll() ([]*models.SelectedInfo, error) {
	var (
		infos []*models.SelectedInfo
	)

	infos, err := models.GetSelectedInfos(st.PageNum, st.PageSize, st.IsValid, st.getMaps())
	if err != nil {
		return nil, err
	}

	return infos, nil
}

func (st *SelectedInfo) GetSelectedAllValid() ([]*models.OuterSelectedInfo, error) {
	infos, err := models.GetValidSelectedInfos()
	if err != nil {
		return nil, err
	}

	return infos, nil
}

func (st *SelectedInfo) AddSelectedInfo() error {
	selectInfo := map[string]interface{}{
		"id":             util.RandomString(16),
		"programId":      st.ProgramId,
		"programsetId":   st.ProgramsetId,
		"programsetName": st.ProgramsetName,
		"programName":    st.ProgramName,
		"playUrl":        st.PlayUrl,
		"cacheSize":      st.CacheSize,
		"cacheValidity":  st.CacheValidity,
		"videoType":      st.VideoType,
	}

	if err := models.AddSelectedInfo(selectInfo); err != nil {
		return err
	}

	return nil
}

func (st *SelectedInfo) EditSelectedInfo() error {

	originInfo, err := models.GetSelectInfoById(st.ProgramId, st.ID)

	if err != nil {
		return err
	}

	if originInfo.ID != "" {
		playCount := st.PlayCount + originInfo.PlayCount
		saveBandwidth := st.SaveBandwidth + originInfo.SaveBandwidth
		return models.EditSelectedInfo(st.ProgramId, st.ID, map[string]interface{}{
			"program_id":     st.ProgramId,
			"play_count":     playCount,
			"save_bandwidth": saveBandwidth,
		})
	} else {
		logging.Warn(fmt.Printf("programId=%d, 可能不在缓存有效期内, 无法修改相关统计信息", st.ProgramId))
		return nil
	}

}

func (st *SelectedInfo) DeleteSelectedInfo() error {
	return models.DeleteSelectedInfo(st.ID)
}

func (st *SelectedInfo) getMaps() map[string]string {
	maps := make(map[string]string)

	maps["programset_name"] = "%" + st.ProgramsetName + "%"
	maps["program_name"] = "%" + st.ProgramName + "%"
	maps["programset_id"] = "%" + st.ProgramsetId + "%"
	maps["program_id"] = "%" + st.ProgramId + "%"

	//if st.IsValid != -1 {
	//	maps["is_valid"] = "%" + strconv.Itoa(int(st.IsValid)) + "%"
	//} else {
	//	maps["is_valid"] = "%%"
	//}

	return maps
}
