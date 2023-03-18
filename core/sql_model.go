package core

import "time"

type SQLModel struct {
	Id        int        `json:"-" gorm:"column:id;" db:"id"`
	FakeId    *UID       `json:"id" gorm:"-"`
	CreatedAt *time.Time `json:"created_at" gorm:"column:created_at;"  db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty" gorm:"column:updated_at;"  db:"updated_at"`
}

func (sqlModel *SQLModel) Mask(dbType int) {
	uid := NewUID(uint32(sqlModel.Id), dbType, 1)
	sqlModel.FakeId = &uid
}
