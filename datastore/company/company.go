package company

import (
	customErr "StudentPlacement/customerrors"
	"StudentPlacement/models"
	"context"
	"database/sql"
)

type Company struct {
	db *sql.DB
}

func New(db *sql.DB) Company {
	return Company{db: db}
}

func (c Company) Get(ctx context.Context, id string) (models.Company, error) {
	rows := c.db.QueryRowContext(ctx, "select id,name,category from company where id = ?", id)

	var company models.Company

	er := rows.Scan(&company.ID, &company.Name, &company.Category)

	if er != nil {
		return models.Company{}, customErr.InternalErrors{Message: er.Error()}
	}

	return company, nil
}

func (c Company) Delete(ctx context.Context, id string) error {
	_, er := c.db.ExecContext(ctx, "delete from company where id = ? ", id)

	if er != nil {
		return customErr.InternalErrors{Message: er.Error()}
	}

	return nil
}

func (c Company) Post(ctx context.Context, company models.Company) error {
	_, er := c.db.ExecContext(ctx, "insert into company values(uuid(),?,?)", company.Name, company.Category)

	if er != nil {
		return customErr.InternalErrors{Message: er.Error()}
	}

	return nil
}

func (c Company) Update(ctx context.Context, company models.Company) error {
	_, er := c.Get(ctx, company.ID)

	if er != nil {
		return customErr.InternalErrors{Message: er.Error()}
	}

	_, err := c.db.ExecContext(ctx, "update company set name = ?,category = ? where id = ?",
		company.Name, company.Category, company.ID)

	if err != nil {
		return customErr.InternalErrors{Message: er.Error()}
	}

	return nil
}
