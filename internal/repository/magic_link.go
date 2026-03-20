package repository

import (
	"chatapp/internal/model"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrMagicLinkAlreadyUsed = fmt.Errorf("magic link already used")
	ErrMagicLinkNotFound    = fmt.Errorf("magic link not found")
)

type MagicLink interface {
	CreateMagicLink(ctx context.Context, magicLink model.MagicLink) (uuid.UUID, error)

	// FindMagicLinkByToken returns a magic link by its token.
	//
	// It returns ErrMagicLinkNotFound if no record exists.
	//
	// It does not validate whether the magic link is expired or already used.
	FindMagicLinkByToken(ctx context.Context, token string) (model.MagicLink, error)

	DeleteMagicLinkByID(ctx context.Context, id uuid.UUID) error

	// UseMagicLink marks a magic link as used.
	//
	// It returns ErrMagicLinkAlreadyUsed if the magic link is already used.
	//
	// It does not validate whether the magic link is expired.
	UseMagicLink(ctx context.Context, id uuid.UUID) error
}

type PostgresMagicLinkRepository struct {
	db *pgxpool.Pool
}

func NewPostgresMagicLinkRepository(db *pgxpool.Pool) MagicLink {
	return &PostgresMagicLinkRepository{
		db: db,
	}
}

// CreateMagicLink implements [MagicLink].
func (p *PostgresMagicLinkRepository) CreateMagicLink(ctx context.Context, magicLink model.MagicLink) (uuid.UUID, error) {
	query := `
		INSERT INTO magic_links (user_id, token, expires_at)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	var id uuid.UUID

	err := p.db.QueryRow(ctx, query, magicLink.UserID, magicLink.Token, magicLink.ExpiresAt).Scan(&id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("create magic link: %w", err)
	}

	return id, nil
}

// DeleteMagicLinkByID implements [MagicLink].
func (p *PostgresMagicLinkRepository) DeleteMagicLinkByID(ctx context.Context, id uuid.UUID) error {
	query := `
		DELETE FROM magic_links
		WHERE id = $1
	`

	_, err := p.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete magic link by id: %w", err)
	}

	return nil
}

// FindMagicLinkByToken implements [MagicLink].
func (p *PostgresMagicLinkRepository) FindMagicLinkByToken(ctx context.Context, token string) (model.MagicLink, error) {
	query := `
		SELECT id, user_id, token, expires_at, used, created_at
		FROM magic_links
		WHERE token = $1
	`

	var magicLink model.MagicLink

	err := p.db.QueryRow(ctx, query, token).Scan(&magicLink.ID, &magicLink.UserID, &magicLink.Token, &magicLink.ExpiresAt, &magicLink.Used, &magicLink.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.MagicLink{}, ErrMagicLinkNotFound
		}

		return model.MagicLink{}, fmt.Errorf("find magic link by token: %w", err)
	}

	return magicLink, nil
}

// UseMagicLink implements [MagicLink].
func (p *PostgresMagicLinkRepository) UseMagicLink(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE magic_links
		SET used = true
		WHERE id = $1 AND used = false
	`

	cmd, err := p.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("use magic link: %w", err)
	}

	if cmd.RowsAffected() == 0 {
		return ErrMagicLinkAlreadyUsed
	}

	return nil
}
