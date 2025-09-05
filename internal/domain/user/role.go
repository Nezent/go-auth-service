package user

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Role struct {
	bun.BaseModel `bun:"table:roles"`
	ID            uuid.UUID     `bun:",pk,type:uuid,default:uuid_generate_v4()"`
	Name          string        `bun:",unique,notnull"`
	Description   string        `bun:",notnull"`
	Permissions   []*Permission `bun:"rel:has-many,join:id=role_id"`
}

type Permission struct {
	bun.BaseModel `bun:"table:permissions"`
	ID            uuid.UUID `bun:",pk,type:uuid,default:uuid_generate_v4()"`
	Name          string    `bun:",unique,notnull"`
	Description   string    `bun:",notnull"`
}

type RolePermission struct {
	bun.BaseModel `bun:"table:role_permissions"`
	RoleID        uuid.UUID   `bun:",notnull"`
	PermissionID  uuid.UUID   `bun:",notnull"`
	Role          *Role       `bun:"rel:belongs-to,join:role_id=id"`
	Permission    *Permission `bun:"rel:belongs-to,join:permission_id=id"`
}

type UserRole struct {
	bun.BaseModel `bun:"table:user_roles"`
	UserID        uuid.UUID `bun:",notnull"`
	RoleID        uuid.UUID `bun:",notnull"`
	User          *User     `bun:"rel:belongs-to,join:user_id=id"`
	Role          *Role     `bun:"rel:belongs-to,join:role_id=id"`
}
