package converter

import (
	"flyhorizons-userservice/models"
	"flyhorizons-userservice/models/enums"
	entities "flyhorizons-userservice/repositories/entity"
	"time"
)

type UserConverter struct{}

func (userConverter *UserConverter) ConvertUserEntityToUser(entity entities.UserEntity) models.User {
	return models.User{
		ID:          entity.ID,
		FullName:    entity.FullName,
		Email:       entity.Email,
		AccountType: enums.AccountTypeFromInt(entity.AccountType),
		Password:    entity.Password,
	}
}

func (userConverter *UserConverter) ConvertUserToUserEntity(user models.User) entities.UserEntity {
	return entities.UserEntity{
		ID:          user.ID,
		FullName:    user.FullName,
		Email:       user.Email,
		AccountType: int(user.AccountType),
		Password:    user.Password,
		// Set the current time of creation/update
		CreatedAt: time.Now(),
	}
}
