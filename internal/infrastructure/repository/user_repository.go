package repository

import (
	"context"
	"time"

	"github.com/Nezent/auth-service/internal/domain/shared"
	"github.com/Nezent/auth-service/internal/domain/user"
	"github.com/Nezent/auth-service/internal/infrastructure/persistence"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type UserRepositoryImpl struct {
	db *persistence.Database
}

// Compile-time interface check
var _ user.UserRepository = (*UserRepositoryImpl)(nil)

func NewUserRepository(db *persistence.Database) *UserRepositoryImpl {
	return &UserRepositoryImpl{
		db: db,
	}
}

func (r *UserRepositoryImpl) Create(user *user.User) (uuid.UUID, *shared.DomainError) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Start transaction
	err := r.db.DB.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		// Insert user
		_, err := tx.NewInsert().Model(user).Exec(ctx)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return uuid.Nil, shared.NewDomainError("CREATE_FAILED", 500, err.Error())
	}
	return user.ID, nil
}

func (r *UserRepositoryImpl) GetByID(id uuid.UUID) (*user.User, *shared.DomainError) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var user user.User
	err := r.db.DB.NewSelect().Model(&user).Where("id = ?", id).Scan(ctx)
	if err != nil {
		return nil, shared.NewDomainError("GET_BY_ID_FAILED", 500, err.Error())
	}
	return &user, nil
}
