package model

import "time"

// CoderModel 计数器模型
type CoderModel struct {
	Id        int32     `gorm:"column:id" json:"id"`
	Code      int32     `gorm:"column:code" json:"code"`
	CreatedAt time.Time `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updatedAt" json:"updatedAt"`
}
