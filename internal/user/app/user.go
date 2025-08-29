package userApp

import (
	"context"

	"github.com/google/uuid"
	"github.com/theguarantors/tiger/entities"
)

//go:generate mockery --name UserRepository --with-expecter --srcpkg github.com/theguarantors/tiger/internal/user/app
type UserRepository interface {
	Create(ctx context.Context, user *entities.User) error
	Get(ctx context.Context, id string) (*entities.User, error)
	Update(ctx context.Context, user *entities.User) error
	Delete(ctx context.Context, id string) error
}

type UserApp struct {
	userRepo UserRepository
	// outside services would be injected here
}

func NewUserApp(userRepo UserRepository) *UserApp {
	return &UserApp{
		userRepo: userRepo,
	}
}

func (u *UserApp) CreateUser(ctx context.Context, user *entities.User) (*entities.User, error) {
	user.ID = uuid.New().String()

	if err := u.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserApp) GetUser(ctx context.Context, id string) (*entities.User, error) {
	return u.userRepo.Get(ctx, id)
}

func (u *UserApp) UpdateUser(ctx context.Context, user *entities.User) error {
	return u.userRepo.Update(ctx, user)
}

func (u *UserApp) DeleteUser(ctx context.Context, id string) error {
	return u.userRepo.Delete(ctx, id)
}
