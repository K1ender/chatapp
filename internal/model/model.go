package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID uuid.UUID

	Email               string
	SigningPublicKey    []byte
	EncryptionPublicKey []byte
	EncryptedPrivateKey []byte

	CreatedAt time.Time
}

type MagicLink struct {
	ID uuid.UUID

	UserID    uuid.UUID
	Token     string
	ExpiresAt time.Time
	Used      bool

	CreatedAt time.Time
}

type Conversation struct {
	ID uuid.UUID

	Name      string
	IsGroup   bool
	CreatedBy uuid.UUID

	CreatedAt time.Time
}

type ConversationMember struct {
	ConversationID uuid.UUID
	UserID         uuid.UUID
	Role           string
	JoinedAt       time.Time
}

type ConversationKey struct {
	ConversationID uuid.UUID
	UserID         uuid.UUID
	EncryptedKey   []byte
	KeyVersion     int

	CreatedAt time.Time
}

type Message struct {
	ID             uuid.UUID
	ConversationID uuid.UUID
	SenderID       uuid.UUID
	Ciphertext     []byte
	Nonce          []byte
	Signature      []byte
	KeyVersion     int

	CreatedAt time.Time
}
