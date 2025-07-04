package interfaces

import (
	entities "flyhorizons-userservice/repositories/entity"
)

type UserRepository interface {
	GetAll() []entities.UserEntity
	GetByID(id int) entities.UserEntity
	GetByEmail(email string) entities.UserEntity
	Create(entities.UserEntity) entities.UserEntity
	DeleteByID(id int) bool
	Update(entities.UserEntity) entities.UserEntity
	SaveLastLoginTime(id int)
}
