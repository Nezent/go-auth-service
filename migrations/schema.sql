-- Users & credentials
CREATE TABLE users (
id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
email CITEXT UNIQUE NOT NULL,
email_verified BOOLEAN NOT NULL DEFAULT false,
created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
disabled BOOLEAN NOT NULL DEFAULT false
);

CREATE TABLE user_credentials (
user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
password_hash BYTEA NOT NULL, -- Argon2id hash
password_set_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- MFA
CREATE TABLE totp_secrets (
user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
secret BYTEA NOT NULL,
created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE webauthn_credentials (
id TEXT PRIMARY KEY,
user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
public_key BYTEA NOT NULL,
sign_count BIGINT NOT NULL DEFAULT 0,
created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Sessions & tokens
CREATE TABLE refresh_tokens (
id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
token_hash BYTEA NOT NULL, -- store hash only
client_id TEXT NOT NULL, -- app id/device id
created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
expires_at TIMESTAMPTZ NOT NULL,
revoked_at TIMESTAMPTZ,
UNIQUE(user_id, client_id, token_hash)
);

-- Service-to-service clients (Client Credentials)
CREATE TABLE service_clients (
id TEXT PRIMARY KEY,
name TEXT NOT NULL,
secret_hash BYTEA NOT NULL,
scopes TEXT[] NOT NULL DEFAULT '{}',
created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
disabled BOOLEAN NOT NULL DEFAULT false
);

-- RBAC
CREATE TABLE roles (
id SERIAL PRIMARY KEY,
name TEXT UNIQUE NOT NULL
);

CREATE TABLE permissions (
);