package repository

import (
	"chatapp/internal/model"
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type User interface {
	CreateUser(ctx context.Context, user model.User) (uuid.UUID, error)

	FindUserByEmail(ctx context.Context, email string) (model.User, error)
	FindUserByID(ctx context.Context, id uuid.UUID) (model.User, error)
}

type PostgresUserRepository struct {
	db *pgxpool.Pool
}

func NewPostgresUserRepository(db *pgxpool.Pool) User {
	return &PostgresUserRepository{
		db: db,
	}
}

// CreateUser implements [UserRepository].
func (p *PostgresUserRepository) CreateUser(ctx context.Context, user model.User) (uuid.UUID, error) {
	query := `
		INSERT INTO users (email, signing_public_key, encryption_public_key, encrypted_private_key)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	var id uuid.UUID

	err := p.db.QueryRow(ctx, query, user.Email, user.SigningPublicKey, user.EncryptionPublicKey, user.EncryptedPrivateKey).Scan(&id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("create user: %w", err)
	}

	return id, nil
}

// FindUserByEmail implements [UserRepository].
func (p *PostgresUserRepository) FindUserByEmail(ctx context.Context, email string) (model.User, error) {
	query := `
		SELECT id, email, signing_public_key, encryption_public_key, encrypted_private_key, created_at
		FROM users
		WHERE email = $1
	`

	var user model.User

	err := p.db.QueryRow(ctx, query, email).Scan(&user.ID, &user.Email, &user.SigningPublicKey, &user.EncryptionPublicKey, &user.EncryptedPrivateKey, &user.CreatedAt)
	if err != nil {
		return model.User{}, fmt.Errorf("find user by email: %w", err)
	}

	return user, nil
}

// FindUserByID implements [UserRepository].
func (p *PostgresUserRepository) FindUserByID(ctx context.Context, id uuid.UUID) (model.User, error) {
	query := `
		SELECT id, email, signing_public_key, encryption_public_key, encrypted_private_key, created_at
		FROM users
		WHERE id = $1
	`

	var user model.User

	err := p.db.QueryRow(ctx, query, id).Scan(&user.ID, &user.Email, &user.SigningPublicKey, &user.EncryptionPublicKey, &user.EncryptedPrivateKey, &user.CreatedAt)
	if err != nil {
		return model.User{}, fmt.Errorf("find user by id: %w", err)
	}

	return user, nil
}
