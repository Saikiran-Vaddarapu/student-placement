package datastore

import (
	"StudentPlacement/models"
	"context"
)

type Student interface {
	ValidateCategoryBranch(student *models.Student) error
	GetAll(ctx context.Context, name string, branch string) ([]models.Student, error)
	GetAllStudentsWithCompany(ctx context.Context, branch string, name string) ([]models.Student, error)
	Get(ctx context.Context, id string) (models.Student, error)
	Delete(ctx context.Context, id string) error
	Post(ctx context.Context, student *models.Student) error
	Update(ctx context.Context, student *models.Student) error
}

type Company interface {
	Get(ctx context.Context, id string) (models.Company, error)
	Delete(ctx context.Context, id string) error
	Post(ctx context.Context, company models.Company) error
	Update(ctx context.Context, company models.Company) error
}
