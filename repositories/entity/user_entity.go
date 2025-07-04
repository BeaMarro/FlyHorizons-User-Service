package entities

import "time"

type UserEntity struct {
	ID          int       `gorm:"column:ID;primaryKey"`
	FullName    string    `gorm:"column:FullName"`
	Email       string    `gorm:"column:Email;unique"`
	AccountType int       `gorm:"column:AccountType"`
	Password    string    `gorm:"column:Password"`
	CreatedAt   time.Time `gorm:"column:CreatedAt"`
}

// Override the default table name
func (UserEntity) TableName() string {
	return "Account"
}
