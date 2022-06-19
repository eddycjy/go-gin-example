package models

type TerminalRecord struct {
	ID             string `gorm:"primary_key" json:"id"`
	ProgramId      int    `json:"programId"`
	SelectedInfoId string `json:"selectedInfoId"`
}
