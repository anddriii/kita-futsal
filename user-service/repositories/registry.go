package repositories

import (
	repositories "github.com/anddriii/kita-futsal/user-service/repositories/user"
	"gorm.io/gorm"
)

type Registry struct {
	db *gorm.DB
}

type IRepoRegistry interface {
	GetUser() repositories.IUserRepo
}

func NewRepoRegistry(db *gorm.DB) IRepoRegistry {
	return &Registry{db: db}
}

// GetUser implements IRepoRegistry.
func (r *Registry) GetUser() repositories.IUserRepo {
	return repositories.NewUserRepo(r.db)
}
