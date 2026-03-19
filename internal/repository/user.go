package repository

import (
	"chatapp/internal/model"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user model.User) (int, error)

	FindUserByEmail(ctx context.Context, email string) (model.User, error)
	FindUserByID(ctx context.Context, id int64) (model.User, error)
}

type PostgresUserRepository struct {
	db *pgxpool.Pool
}

func NewPostgresUserRepository(db *pgxpool.Pool) UserRepository {
	return &PostgresUserRepository{
		db: db,
	}
}

// CreateUser implements [UserRepository].
func (p *PostgresUserRepository) CreateUser(ctx context.Context, user model.User) (int, error) {
	panic("unimplemented")
}

// FindUserByEmail implements [UserRepository].
func (p *PostgresUserRepository) FindUserByEmail(ctx context.Context, email string) (model.User, error) {
	panic("unimplemented")
}

// FindUserByID implements [UserRepository].
func (p *PostgresUserRepository) FindUserByID(ctx context.Context, id int64) (model.User, error) {
	panic("unimplemented")
}
