// Enchanted Garden/model/branch.go
package model

import "time"

//структуры

type Branch struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"type:varchar(200);not null" json:"name"`
	ParentID  *uint     `gorm:"parent_id" json:"ParentID"`
	CreatedAt time.Time `gorm:"created_at" json:"created_at"`
	Flora     []Flora   `gorm:"foreignKey:BranchID;constraint:OnDelete:CASCADE" json:"flora,omitempty"`
	Children  []Branch  `gorm:"foreignKey:ParentID" json:"children,omitempty"`
}

type CreateBranchReq struct {
	Name     string `json:"name"`
	ParentID *uint  `json:"parent_id"`
}

type UpdateBranchReq struct {
	Name     *string `json:"name"`
	ParentID *uint   `json:"parent_id"`
}
