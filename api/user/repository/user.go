package userRepository

import (
	"context"

	userStructs "github.com/theguarantors/tiger/api/structs"

	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Create(ctx context.Context, user *userStructs.User) error {
	query := `INSERT INTO users (id, name, email) VALUES ($1, $2, $3)`
	_, err := r.db.Exec(query, user.ID, user.Name, user.Email)
	return err
}

func (r *UserRepository) Get(ctx context.Context, id string) (*userStructs.User, error) {
	user := &userStructs.User{}
	query := `SELECT id, name, email FROM users WHERE id = $1`

	if err := r.db.Get(user, query, id); err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) Update(ctx context.Context, user *userStructs.User) error {
	return nil
}

func (r *UserRepository) Delete(ctx context.Context, id string) error {
	return nil
}
