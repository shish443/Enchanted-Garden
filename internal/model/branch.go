// Enchanted-Garden/internal/model/branch.go
package model

import "time"

type Branch struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"type:varchar(200);not null" json:"name"`
	ParentID  *uint     `gorm:"column:parent_id" json:"parent_id"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	Flora     []Flora   `gorm:"foreignKey:BranchID;constraint:OnDelete:CASCADE" json:"employees,omitempty"`
	Children  []Branch  `gorm:"foreignKey:ParentID;constraint:OnDelete:CASCADE" json:"children,omitempty"`
}

type CreateBranchReq struct {
	Name     string `json:"name"`
	ParentID *uint  `json:"parent_id"`
}

type UpdateBranchReq struct {
	Name     *string `json:"name"`
	ParentID **uint  `json:"parent_id"`
}
