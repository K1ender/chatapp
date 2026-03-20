-- +goose Up
CREATE EXTENSION IF NOT EXISTS "citext";

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email citext NOT NULL UNIQUE,

    signing_public_key BYTEA NOT NULL, --ed25519
    encryption_public_key BYTEA NOT NULL, --x25519

    encrypted_private_key BYTEA,

    created_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE magic_links (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    t   n TEXT NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    used BOOLEAN DEFAULT FALSE,

    created_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE conversations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    name TEXT,
    is_group BOOLEAN DEFAULT FALSE,
    created_by UUID REFERENCES users(id),

    created_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE conversations_members (
    conversation_id UUID REFERENCES conversations(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    role TEXT NOT NULL DEFAULT 'member',
    joined_at TIMESTAMPTZ DEFAULT now(),

    PRIMARY KEY (conversation_id, user_id)
);

CREATE TABLE conversation_keys (
    conversation_id UUID REFERENCES conversations(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    encrypted_key BYTEA NOT NULL,
    key_version INT NOT NULL,

    created_at TIMESTAMPTZ DEFAULT now(),

    PRIMARY KEY (conversation_id, user_id, key_version)
);

CREATE TABLE messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    conversation_id UUID REFERENCES conversations(id) ON DELETE CASCADE,
    sender_id UUID REFERENCES users(id) ON DELETE CASCADE,
    ciphertext BYTEA NOT NULL,
    nonce BYTEA NOT NULL,
    signature BYTEA NOT NULL,
    key_version INT NOT NULL,

    created_at TIMESTAMPTZ DEFAULT now()
);

CREATE INDEX idx_messages_conversation_id ON messages(conversation_id);
CREATE INDEX idx_messages_created_at ON messages(created_at);
CREATE INDEX idx_conversations_members_user ON conversations_members(user_id);
CREATE INDEX idx_conversations_created_by ON conversations(created_by);

-- +goose Down
DROP TABLE users;
DROP TABLE magic_links;
DROP TABLE conversations;
DROP TABLE conversations_members;
DROP TABLE conversation_keys;
DROP TABLE messages;

DROP INDEX idx_messages_conversation_id;
DROP INDEX idx_messages_created_at;
DROP INDEX idx_conversations_members_user;
DROP INDEX idx_conversations_created_by;

DROP EXTENSION IF EXISTS "citext";
