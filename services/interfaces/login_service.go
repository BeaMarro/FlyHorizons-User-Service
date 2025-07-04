package interfaces

import (
	"flyhorizons-userservice/models/request"
	"flyhorizons-userservice/models/response"
)

type LoginService interface {
	Login(request.LoginRequest, string) (*response.LoginResponse, error)
}
