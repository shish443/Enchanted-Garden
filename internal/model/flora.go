// Enchanted Garden/model/flora.go
package model

import "time"

type Flora struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	BranchID  uint       `gorm:"not null;index" json:"branch_id"`
	FullName  string     `gorm:"type:varchar(200);not null" json:"full_name"`
	Position  string     `gorm:"type:varchar(200);not null" json:"position"`
	HiredAt   *time.Time `gorm:"type:date" json:"hired_at"`
	CreatedAt time.Time  `json:"created_at"`
}
type PlantFloraReq struct {
	FullName string     `json:"full_name"`
	HiredAt  *time.Time `json:"hired_at"`
	Position string     `json:"position"`
}
