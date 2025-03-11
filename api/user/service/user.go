package userService

import (
	"context"
	"errors"
	"net/http"

	"github.com/theguarantors/tiger/api/structs"
	utils "github.com/theguarantors/tiger/utils"
)

//go:generate mockery --name UserApp --with-expecter --srcpkg github.com/theguarantors/tiger/api/user/service
type UserApp interface {
	CreateUser(ctx context.Context, user *structs.User) (*structs.User, error)
	GetUser(ctx context.Context, id string) (*structs.User, error)
	UpdateUser(ctx context.Context, user *structs.User) error
	DeleteUser(ctx context.Context, id string) error
}

type UserService struct {
	userApp UserApp
}

func NewUserService(userApp UserApp) *UserService {
	return &UserService{
		userApp: userApp,
	}
}

func (s *UserService) CreateUser(ctx context.Context, request *http.Request) *http.Response {
	user, err := s.validateUserCreationRequest(request)
	if err != nil {
		// log error
		return utils.ServerResponse(nil, err)
	}

	return utils.ServerResponse(s.userApp.CreateUser(ctx, user))
}

func (s *UserService) GetUser(ctx context.Context, request *http.Request) *http.Response {
	id := utils.GetPathParam(request, "id")
	if id == "" {
		// log error
		return utils.ServerResponse(nil, errors.New("id is required"))
	}

	return utils.ServerResponse(s.userApp.GetUser(ctx, id))
}

// private methods always in the bottom of the file
func (s *UserService) validateUserCreationRequest(request *http.Request) (*structs.User, error) {
	var user *structs.User

	if err := utils.UnmarshalJSON(request, &user); err != nil {
		return nil, err
	}

	if user.Email == "" {
		return nil, errors.New("email is required")
	}

	if user.Name == "" {
		return nil, errors.New("name is required")
	}

	return user, nil
}
