package service

import (
	"StudentPlacement/models"
	"context"
)

type Company interface {
	Get(ctx context.Context, id string) (models.Company, error)
	Update(ctx context.Context, company models.Company) error
	Post(ctx context.Context, company models.Company) error
	Delete(ctx context.Context, id string) error
}

type Student interface {
	GetAll(ctx context.Context, name string, branch string, flag string) ([]models.Student, error)
	Get(ctx context.Context, id string) (models.Student, error)
	Update(ctx context.Context, student *models.Student) error
	Post(ctx context.Context, student *models.Student) error
	Delete(ctx context.Context, id string) error
}
