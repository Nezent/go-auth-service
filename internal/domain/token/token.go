package token

import (
	"time"

	"github.com/Nezent/auth-service/internal/domain/auth"
	"github.com/Nezent/auth-service/internal/domain/user"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type AuthSession struct {
	bun.BaseModel `bun:"table:auth_sessions"`
	ID            uuid.UUID       `bun:",pk,type:uuid,default:uuid_generate_v4()"`
	UserID        uuid.UUID       `bun:",notnull"`
	LoginID       uuid.UUID       `bun:",notnull"`
	Token         string          `bun:",notnull"`
	ExpiresAt     time.Time       `bun:",notnull"`
	User          *user.User      `bun:"rel:belongs-to,join:user_id=id"`
	UserLogin     *auth.UserLogin `bun:"rel:belongs-to,join:login_id=id"`
}
