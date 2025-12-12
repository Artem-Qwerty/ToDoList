package quote_interface

type Quote struct {
	ID      int    `gorm:"column:id;primaryKey"`
	Text    string `gorm:"column:text"`
	Current bool   `gorm:"column:current"`
}

func (Quote) TableName() string {
	return "quotes"
}
