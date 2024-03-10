package models

import "github.com/VATUSA/primary-api/pkg/database"

type DivisionStaffPosition struct {
	PositionId uint   `gorm:"primaryKey" json:"position_id" example:"1"`
	CID        uint   `json:"cid" example:"1293257"`
	User       *User  `json:"-" gorm:"foreignKey:CID"`
	Title      string `json:"title" example:"Division Director" gorm:"size:120"`
}

func (p *DivisionStaffPosition) Create() error {
	return database.DB.Create(p).Error
}

func (p *DivisionStaffPosition) Update() error {
	return database.DB.Save(p).Error
}

func (p *DivisionStaffPosition) Delete() error {
	return database.DB.Delete(p).Error
}

func (p *DivisionStaffPosition) Get() error {
	return database.DB.Where("position_id = ?", p.PositionId).Preload("User").First(p).Error
}

func GetAllDivisionStaffPositions() ([]DivisionStaffPosition, error) {
	var positions []DivisionStaffPosition
	return positions, database.DB.Find(&positions).Error
}
