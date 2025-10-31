package repository

import (
	"quizizz.com/internal/domain"
)

// UserRepository defines the interface for user data access
type UserRepository interface {
	GetByID(id string) (*domain.User, error)
	List() ([]*domain.User, error)
	Create(user *domain.User) error
	Update(user *domain.User) error
	Delete(id string) error
}

// This would typically be implemented with a concrete type, like:
// type sqlUserRepository struct {
//     db *sql.DB
// }
//
// func NewUserRepository(db *sql.DB) UserRepository {
//     return &sqlUserRepository{
//         db: db,
//     }
// }
//
// Then implement all the interface methods using the DB
