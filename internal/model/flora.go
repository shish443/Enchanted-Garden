// Enchanted-Garden/internal/model/flora.go
package model

import (
	"strings"
	"time"
)

type Flora struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	BranchID  uint       `gorm:"not null;index" json:"department_id"`
	FullName  string     `gorm:"type:varchar(200);not null" json:"full_name"`
	Position  string     `gorm:"type:varchar(200);not null" json:"position"`
	HiredAt   *time.Time `gorm:"type:date" json:"hired_at"`
	CreatedAt time.Time  `json:"created_at"`
}

type PlantFloraReq struct {
	FullName string    `json:"full_name"`
	HiredAt  *DateOnly `json:"hired_at"`
	Position string    `json:"position"`
}

type DateOnly time.Time

func (d *DateOnly) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	if s == "null" || s == "" {
		return nil
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}
	*d = DateOnly(t)
	return nil
}

func (d *DateOnly) MarshalJSON() ([]byte, error) {
	if d == nil {
		return []byte("null"), nil
	}
	return []byte(`"` + time.Time(*d).Format("2006-01-02") + `"`), nil
}
