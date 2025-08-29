package policyRepository

import (
	"context"

	"github.com/theguarantors/tiger/entities"

	"github.com/jmoiron/sqlx"
)

type PolicyRepository struct {
	db *sqlx.DB
}

func NewPolicyRepository(db *sqlx.DB) *PolicyRepository {
	return &PolicyRepository{
		db: db,
	}
}

func (r *PolicyRepository) Get(ctx context.Context, id string) (*entities.Policy, error) {
	policy := &entities.Policy{}
	query := `SELECT id, name, email FROM policy WHERE id = $1`

	if err := r.db.Get(policy, query, id); err != nil {
		return nil, err
	}

	return policy, nil
}
