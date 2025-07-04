package services

import (
	"flyhorizons-userservice/models"
	"flyhorizons-userservice/models/enums"
	"flyhorizons-userservice/models/request"
	"flyhorizons-userservice/models/response"
	"flyhorizons-userservice/services/authentication"
	"flyhorizons-userservice/services/converter"
	"flyhorizons-userservice/services/errors"
	"flyhorizons-userservice/services/interfaces"
	"log"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type OAuthTokenSigner struct {
	jwtSigner *authentication.JwtTokenSigner
}

func NewOAuthTokenSigner(jwtSigner *authentication.JwtTokenSigner) *OAuthTokenSigner {
	return &OAuthTokenSigner{
		jwtSigner: jwtSigner,
	}
}

func (signer *OAuthTokenSigner) SignToken(claims jwt.Claims) (string, error) {
	return signer.jwtSigner.SignToken(claims)
}

type LoginService struct {
	repo          interfaces.UserRepository
	userConverter converter.UserConverter
	tokenSigner   interfaces.TokenSigner
}

func NewLoginService(repo interfaces.UserRepository, userConverter converter.UserConverter, tokenSigner interfaces.TokenSigner) *LoginService {
	return &LoginService{
		repo:          repo,
		userConverter: userConverter,
		tokenSigner:   tokenSigner,
	}
}

func (service *LoginService) Login(loginRequest request.LoginRequest, ip string) (*response.LoginResponse, error) {
	accountEntity := service.repo.GetByEmail(loginRequest.Email)
	account := service.userConverter.ConvertUserEntityToUser(accountEntity)

	// Unsuccessful login attempt
	if !service.matchesPassword(loginRequest.Password, account.Password) {
		log.Printf(
			"Unsuccessful login attempt:\n  User ID: %v\n  Timestamp: %s\n  IP Address: %s",
			account.ID,
			time.Now().Format(time.RFC3339),
			ip,
		)
		return nil, errors.NewInvalidCredentialsError(400)
	}

	// Generate OAuth Token
	accessToken, err := service.generateOAuthToken(account)

	if err != nil {
		return nil, err
	}

	// Successful login attempt
	log.Printf(
		"Successful login attempt:\n  User ID: %v\n  Timestamp: %s\n  IP Address: %s",
		account.ID,
		time.Now().Format(time.RFC3339),
		ip,
	)

	return &response.LoginResponse{AccessToken: accessToken}, nil
}

func (service *LoginService) matchesPassword(rawPassword, encodedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(encodedPassword), []byte(rawPassword))
	return err == nil
}

func (service *LoginService) generateOAuthToken(account models.User) (string, error) {
	var role = ""

	if account.AccountType == enums.User {
		role = "user"
	} else if account.AccountType == enums.Admin {
		role = "admin"
	} else {
		return "", errors.NewInvalidAccountTypeError(401)
	}

	// OAuth compliant claims
	claims := jwt.MapClaims{
		"sub":        account.ID,                            // Subject (user ID)
		"email":      account.Email,                         // User email
		"account_id": account.ID,                            // User ID (kept for not crashing the frontend)
		"role":       role,                                  // User role
		"iss":        "flyhorizons-user-service",            // Issuer
		"aud":        "flyhorizons-api",                     // Audience
		"iat":        time.Now().Unix(),                     // Issued at
		"exp":        time.Now().Add(72 * time.Hour).Unix(), // Expiration
	}

	// Save last login time
	service.repo.SaveLastLoginTime(account.ID)

	return service.tokenSigner.SignToken(claims)
}
