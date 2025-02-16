package seeders

import "gorm.io/gorm"

type Registry struct {
	db *gorm.DB
}

type ISeedRegistry interface {
	Run()
}

func NewSeederRegistry(db *gorm.DB) ISeedRegistry {
	return &Registry{db: db}
}

func (r *Registry) Run() {
	RunRoleSeeder(r.db)
	RunUserSeeder(r.db)
}
