package services

import (
	"encoding/json"
	"flyhorizons-userservice/config"
	"flyhorizons-userservice/models"
	"flyhorizons-userservice/models/enums"
	"flyhorizons-userservice/services/authentication"
	"flyhorizons-userservice/services/converter"
	"flyhorizons-userservice/services/errors"
	"flyhorizons-userservice/services/interfaces"
	"flyhorizons-userservice/services/validation"
	"log"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

type UserService struct {
	userRepo           interfaces.UserRepository
	accountHashing     *authentication.AccountHashing
	passwordValidation validation.PasswordValidator
	userConverter      converter.UserConverter
}

func NewUserService(repo interfaces.UserRepository, accountHashing *authentication.AccountHashing, passwordValidator validation.PasswordValidator, userConverter converter.UserConverter) *UserService {
	return &UserService{
		userRepo:       repo,
		accountHashing: accountHashing,
		userConverter:  userConverter,
	}
}

func (userService *UserService) GetAll() []models.User {
	userEntities := userService.userRepo.GetAll()

	var users []models.User
	for _, userEntity := range userEntities {
		user := userService.userConverter.ConvertUserEntityToUser(userEntity)
		users = append(users, user)
	}

	return users
}

func (userService *UserService) GetByID(id int) (*models.User, error) {
	userEntity := userService.userRepo.GetByID(id)
	// User is not found
	if userEntity.ID == 0 {
		return nil, errors.NewUserNotFoundError(id, 404)
	}

	// User is found
	user := userService.userConverter.ConvertUserEntityToUser(userEntity)
	return &user, nil
}

func (UserService *UserService) UserExists(id int) bool {
	for _, user := range UserService.GetAll() {
		if user.ID == id {
			return true
		}
	}
	return false
}

func (userService *UserService) Create(user models.User) (*models.User, error) {
	if userService.UserExists(user.ID) {
		return nil, errors.NewUserExistsError(user.ID, 409)
	}

	// Validate password
	err := userService.passwordValidation.Validate(user.Password)
	if err != nil {
		return nil, err
	}

	// Encode password
	hashedPassword, err := userService.accountHashing.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}
	user.Password = hashedPassword

	var userEntity = userService.userConverter.ConvertUserToUserEntity(user)
	var postUserEntity = userService.userRepo.Create(userEntity)
	var postUser = userService.userConverter.ConvertUserEntityToUser(postUserEntity)

	// Successful account creation
	accountTypeInt := enums.AccountTypeFromInt(int(postUser.AccountType))
	var accountType string

	switch accountTypeInt {
	case 0:
		accountType = "Admin"
	case 1:
		accountType = "User"
	}

	log.Printf(
		"Successfully created account:\n  User ID: %v\n  Account Type: %v\n  Timestamp: %s",
		postUser.ID,
		accountType,
		time.Now().Format(time.RFC3339),
	)

	return &postUser, nil
}

func (userService *UserService) DeleteByID(id int) (bool, error) {
	if !userService.UserExists(id) {
		return false, errors.NewUserNotFoundError(id, 404)
	}

	// Delete data from the user database
	isDeleted := userService.userRepo.DeleteByID(id)

	// Delete data from any other databases (containing user data)
	if isDeleted {
		// Post user_deleted event to RabbitMQ
		channel := config.RabbitMQClient.Channel

		body, err := json.Marshal(struct {
			UserID int `json:"userId"`
		}{
			UserID: id,
		})

		if err != nil {
			log.Printf("Failed to marshal user deletion event: %v", err)
		}

		err = channel.Publish(
			"",
			"user_deleted",
			false,
			false,
			amqp091.Publishing{
				ContentType: "application/json",
				Body:        body,
			},
		)

		if err != nil {
			log.Printf("An error occurred while posting the messaging to RabbitMQ %v\n", err)
		}
	}

	// Successful account deletion
	log.Printf(
		"Successfully deleted account:\n  User ID: %v\n Timestamp: %s",
		id,
		time.Now().Format(time.RFC3339),
	)

	return isDeleted, nil
}

func (userService *UserService) Update(user models.User) (*models.User, error) {
	if !userService.UserExists(user.ID) {
		return nil, errors.NewUserNotFoundError(user.ID, 404)
	}

	// Validate password
	err := userService.passwordValidation.Validate(user.Password)
	if err != nil {
		return nil, err
	}

	// Encode password
	hashedPassword, err := userService.accountHashing.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}
	user.Password = hashedPassword

	var userEntity = userService.userConverter.ConvertUserToUserEntity(user)
	var putUserEntity = userService.userRepo.Update(userEntity)
	var putUser = userService.userConverter.ConvertUserEntityToUser(putUserEntity)

	// Successful account updating
	log.Printf(
		"Successfully updated account:\n  User ID: %v\n  Account Type: %v\n  Timestamp: %s",
		putUser.ID,
		enums.AccountTypeFromInt(int(putUser.AccountType)),
		time.Now().Format(time.RFC3339),
	)

	return &putUser, nil
}
