package models

import "flyhorizons-userservice/models/enums"

type User struct {
	ID          int               `json:"id"`
	FullName    string            `json:"full_name"`
	Email       string            `json:"email"`
	AccountType enums.AccountType `json:"account_type"`
	Password    string            `json:"password"`
}
