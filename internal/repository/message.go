package repository

import (
	"chatapp/internal/model"
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Message interface {
	CreateMessage(ctx context.Context, message model.Message) (uuid.UUID, error)
	FindMessagesByConversationID(ctx context.Context, conversationID uuid.UUID) ([]model.Message, error)
}

type PostgresMessageRepository struct {
	db *pgxpool.Pool
}

func NewPostgresMessageRepository(db *pgxpool.Pool) Message {
	return &PostgresMessageRepository{
		db: db,
	}
}

// CreateMessage implements [Message].
func (p *PostgresMessageRepository) CreateMessage(ctx context.Context, message model.Message) (uuid.UUID, error) {
	query := `
		INSERT INTO messages (conversation_id, sender_id, ciphertext, nonce, signature, key_version)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	var id uuid.UUID

	err := p.db.QueryRow(ctx, query, message.ConversationID, message.SenderID, message.Ciphertext, message.Nonce, message.Signature, message.KeyVersion).Scan(&id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("create message: %w", err)
	}

	return id, nil
}

// FindMessagesByConversationID implements [Message].
func (p *PostgresMessageRepository) FindMessagesByConversationID(ctx context.Context, conversationID uuid.UUID) ([]model.Message, error) {
	query := `
		SELECT id, conversation_id, sender_id, ciphertext, nonce, signature, key_version, created_at
		FROM messages
		WHERE conversation_id = $1
		ORDER BY created_at ASC
		LIMIT 50
	`

	rows, err := p.db.Query(ctx, query, conversationID)
	if err != nil {
		return nil, fmt.Errorf("find messages by conversation id: %w", err)
	}
	defer rows.Close()

	var messages []model.Message
	for rows.Next() {
		var message model.Message
		err := rows.Scan(&message.ID, &message.ConversationID, &message.SenderID, &message.Ciphertext, &message.Nonce, &message.Signature, &message.KeyVersion, &message.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("find messages by conversation id: %w", err)
		}
		messages = append(messages, message)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return messages, nil
}
