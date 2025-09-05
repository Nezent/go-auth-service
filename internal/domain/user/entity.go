package user

import (
	"context"
	"time"

	"github.com/Nezent/auth-service/internal/domain/shared"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`

	ID        uuid.UUID `bun:",pk,type:uuid,default:uuid_generate_v4()"`
	Email     string    `bun:",unique,notnull"`
	Password  string    `bun:",notnull"`
	Salt      string    `bun:",notnull"`
	CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
}

// Compile-time interface checks for Bun hooks
var _ bun.BeforeAppendModelHook = (*User)(nil)
var _ bun.BeforeSelectHook = (*User)(nil)

// BeforeSelect implements bun.BeforeSelectHook for multi-tenancy or context-based filtering.
func (u *User) BeforeSelect(ctx context.Context, query *bun.SelectQuery) error {
	if tenantID := ctx.Value("tenant_id"); tenantID != nil {
		query.Where("tenant_id = ?", tenantID)
	}
	return nil
}

// BeforeAppendModel sets timestamps before insert/update.
func (u *User) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	now := time.Now().UTC()
	switch query.(type) {
	case *bun.InsertQuery:
		u.CreatedAt = now
		u.UpdatedAt = now
	case *bun.UpdateQuery:
		u.UpdatedAt = now
	}
	return nil
}

// UserRepository defines the methods for user data access.
type UserRepository interface {
	Create(user *User) (uuid.UUID, *shared.DomainError)
	GetByID(id uuid.UUID) (*User, *shared.DomainError)
	// Update(user *User) *shared.DomainError
	// Delete(id uuid.UUID) *shared.DomainError
}
