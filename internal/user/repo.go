package user

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/pmoieni/auth-server/internal/user/db/sqlc"
)

type Repo interface {
	CreateUser(context.Context, User) (*User, error)
	GetUsers(context.Context) ([]*User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
}

type repo struct {
	q *sqlc.Queries
}

func New(conn sqlc.DBTX) Repo {
	return &repo{q: sqlc.New(conn)}
}

func (r *repo) CreateUser(ctx context.Context, u User) (*User, error) {
	created, err := r.q.CreateUser(ctx, &sqlc.CreateUserParams{
		Username: u.Username,
		Email:    u.Email,
	})

	return &User{
		ID:       created.ID,
		Username: created.Username,
		Email:    created.Email.(string),
	}, err
}

func (r *repo) GetUsers(ctx context.Context) ([]*User, error) {
	res := make([]*User, 0)
	users, err := r.q.ListUsers(ctx, &sqlc.ListUsersParams{
		// TODO: Paginate results.
		Limit: 20,
	})
	if err != nil {
		return res, fmt.Errorf("listUsers: %w", err)
	}

	for _, u := range users {
		res = append(res, &User{
			ID:       u.ID,
			Username: u.Username,
			Email:    u.Email.(string),
		})
	}
	return res, nil
}

func (r *repo) GetUserByID(ctx context.Context, id uuid.UUID) (*User, error) {
	found, err := r.q.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &User{
		ID:       found.ID,
		Username: found.Username,
		Email:    found.Email.(string),
	}, nil
}

func (r *repo) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	found, err := r.q.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return &User{
		ID:       found.ID,
		Username: found.Username,
		Email:    found.Email.(string),
	}, nil
}
