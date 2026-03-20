package repository

import (
	"chatapp/internal/model"
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Conversation interface {
	CreateConversation(ctx context.Context, conversation model.Conversation) (uuid.UUID, error)
	CreateConversationKey(ctx context.Context, conversationID uuid.UUID, keys model.ConversationKey) error

	FindConversationByID(ctx context.Context, id uuid.UUID) (model.Conversation, error)

	DeleteConversationByID(ctx context.Context, id uuid.UUID) error

	AddMemberToConversation(ctx context.Context, conversationID uuid.UUID, userID uuid.UUID) error
}

type PostgresConversationRepository struct {
	db *pgxpool.Pool
}

func NewPostgresConversationRepository(db *pgxpool.Pool) Conversation {
	return &PostgresConversationRepository{
		db: db,
	}
}

// CreateConversation implements [Conversation].
func (p *PostgresConversationRepository) CreateConversation(ctx context.Context, conversation model.Conversation) (uuid.UUID, error) {
	query := `
		INSERT INTO conversations (name, is_group, created_by)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	var id uuid.UUID

	err := p.db.QueryRow(ctx, query, conversation.Name, conversation.IsGroup, conversation.CreatedBy).Scan(&id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("create conversation: %w", err)
	}

	return id, nil
}

// CreateConversationKey implements [Conversation].
func (p *PostgresConversationRepository) CreateConversationKey(ctx context.Context, conversationID uuid.UUID, keys model.ConversationKey) error {
	query := `
    	INSERT INTO conversation_keys (conversation_id, user_id, encrypted_key, key_version)
    	VALUES ($1, $2, $3, $4)
	`

	_, err := p.db.Exec(ctx, query,
		conversationID,
		keys.UserID,
		keys.EncryptedKey,
		keys.KeyVersion,
	)
	if err != nil {
		return fmt.Errorf("create conversation key: %w", err)
	}

	return nil
}

// DeleteConversationByID implements [Conversation].
func (p *PostgresConversationRepository) DeleteConversationByID(ctx context.Context, id uuid.UUID) error {
	query := `
		DELETE FROM conversations
		WHERE id = $1
	`

	_, err := p.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete conversation by id: %w", err)
	}

	return nil
}

// FindConversationByID implements [Conversation].
func (p *PostgresConversationRepository) FindConversationByID(ctx context.Context, id uuid.UUID) (model.Conversation, error) {
	query := `
		SELECT id, name, is_group, created_by, created_at
		FROM conversations
		WHERE id = $1
	`

	var conversation model.Conversation

	err := p.db.QueryRow(ctx, query, id).Scan(&conversation.ID, &conversation.Name, &conversation.IsGroup, &conversation.CreatedBy, &conversation.CreatedAt)
	if err != nil {
		return model.Conversation{}, fmt.Errorf("find conversation by id: %w", err)
	}

	return conversation, nil
}

// AddMemberToConversation implements [Conversation].
func (p *PostgresConversationRepository) AddMemberToConversation(ctx context.Context, conversationID uuid.UUID, userID uuid.UUID) error {
	query := `
		INSERT INTO conversation_members (conversation_id, user_id)
		VALUES ($1, $2) ON CONFLICT DO NOTHING
	`

	_, err := p.db.Exec(ctx, query, conversationID, userID)
	if err != nil {
		return fmt.Errorf("add member to conversation: %w", err)
	}

	return nil
}
