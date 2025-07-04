package repositories

import (
	entities "flyhorizons-userservice/repositories/entity"
	"flyhorizons-userservice/services/interfaces"
	"time"
)

type UserRepository struct {
	*BaseRepository
}

var _ interfaces.UserRepository = (*UserRepository)(nil)

func NewUserRepository(baseRepo *BaseRepository) *UserRepository {
	return &UserRepository{
		BaseRepository: baseRepo,
	}
}

func (repo *UserRepository) GetAll() []entities.UserEntity {
	db, _ := repo.CreateConnection()

	var users []entities.UserEntity
	db.Find(&users)

	return users
}

func (repo *UserRepository) GetByID(id int) entities.UserEntity {
	db, _ := repo.CreateConnection()

	var user entities.UserEntity
	db.First(&user, id)

	return user
}

func (repo *UserRepository) GetByEmail(email string) entities.UserEntity {
	db, _ := repo.CreateConnection()

	var user entities.UserEntity
	db.Where("email = ?", email).First(&user)

	return user
}

func (repo *UserRepository) Create(userEntity entities.UserEntity) entities.UserEntity {
	db, _ := repo.CreateConnection()

	db.Create(&userEntity)

	return userEntity
}

func (repo *UserRepository) DeleteByID(id int) bool {
	db, _ := repo.CreateConnection()

	result := db.Delete(&entities.UserEntity{}, id)

	if result.Error != nil || result.RowsAffected == 0 {
		return false
	}

	return true
}

func (repo *UserRepository) Update(userEntity entities.UserEntity) entities.UserEntity {
	db, _ := repo.CreateConnection()

	db.Save(&userEntity)

	return userEntity
}

func (repo *UserRepository) SaveLastLoginTime(userID int) {
	db, _ := repo.CreateConnection()

	db.Model(&entities.UserEntity{}).
		Where("id = ?", userID).
		Update("LastLogin", time.Now())
}
