package repository

import (
	"chatapp/internal/model"
	"context"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user model.User) (int, error)

	FindUserByEmail(ctx context.Context, email string) (model.User, error)
	FindUserByID(ctx context.Context, id int64) (model.User, error)
}

type PostgresUserRepository struct{
	
}
