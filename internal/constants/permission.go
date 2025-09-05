package constants

type Permission string

const (
	// Org permissions
	OrgViewHealth  Permission = "org:view_health"
	OrgInviteUser  Permission = "org:invite_user"
	OrgViewMembers Permission = "org:view_members"
	OrgDelete      Permission = "org:delete"
	OrgUpdate      Permission = "org:update"

	// User permissions
	UserCreate           Permission = "user:create"
	UserView             Permission = "user:view"
	UserUpdate           Permission = "user:update"
	UserDelete           Permission = "user:delete"
	UserUpdatePermission Permission = "user:update_permission"
	UserUpdateRole       Permission = "user:update_role"
)
