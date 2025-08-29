package userService

import (
	"context"
	"errors"
	"net/http"

	logger "github.com/TheGuarantors/tg-logger/pkg"
	"github.com/theguarantors/tiger/internal/entities"
	utils "github.com/theguarantors/tiger/utils"
)

//go:generate mockery --name UserApp --with-expecter --srcpkg github.com/theguarantors/tiger/internal/user/service
type UserApp interface {
	CreateUser(ctx context.Context, user *entities.User) (*entities.User, error)
	GetUser(ctx context.Context, id string) (*entities.User, error)
	UpdateUser(ctx context.Context, user *entities.User) error
	DeleteUser(ctx context.Context, id string) error
}

type UserService struct {
	userApp UserApp
	logger  *logger.Logger
}

func NewUserService(userApp UserApp, logger *logger.Logger) *UserService {
	return &UserService{
		userApp: userApp,
		logger:  logger,
	}
}

func (s *UserService) CreateUser(ctx context.Context, request *http.Request) *http.Response {
	user, err := s.validateUserCreationRequest(request)
	if err != nil {
		// log error
		return utils.ServerResponse(ctx, nil, err, s.logger)
	}

	user, err = s.userApp.CreateUser(ctx, user)

	return utils.ServerResponse(ctx, user, err, s.logger)
}

func (s *UserService) GetUser(ctx context.Context, request *http.Request) *http.Response {
	id := utils.GetPathParam(request, "id")
	if id == "" {
		// log error
		return utils.ServerResponse(ctx, nil, errors.New("id is required"), s.logger)
	}

	user, err := s.userApp.GetUser(ctx, id)

	return utils.ServerResponse(ctx, user, err, s.logger)
}

// private methods always in the bottom of the file
func (s *UserService) validateUserCreationRequest(request *http.Request) (*entities.User, error) {
	var user *entities.User

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
