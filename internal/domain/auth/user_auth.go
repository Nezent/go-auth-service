package auth

import (
	"time"

	"github.com/Nezent/auth-service/internal/domain/user"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type UserLogin struct {
	bun.BaseModel `bun:"table:user_logins"`
	ID            uuid.UUID  `bun:",pk,type:uuid,default:uuid_generate_v4()"`
	UserID        uuid.UUID  `bun:",notnull"`
	Count         int        `bun:",notnull,default:0"`
	DeviceID      string     `bun:",notnull"`
	LastLoginAt   time.Time  `bun:",notnull"`
	User          *user.User `bun:"rel:belongs-to,join:user_id=id"`
}
