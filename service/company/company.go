package company

import (
	customErr "StudentPlacement/customerrors"
	"StudentPlacement/datastore"
	"StudentPlacement/models"
	"context"
)

type Company struct {
	store datastore.Company
}

func New(store datastore.Company) Company {
	return Company{store: store}
}

func validate(company models.Company) error {
	category := company.Category

	if !models.Contains(models.ValidCategory, category) {
		return &customErr.Validity{Message: "invalid Category"}
	}

	return nil
}

func (c Company) Get(ctx context.Context, id string) (models.Company, error) {
	res, er := c.store.Get(ctx, id)

	if er != nil {
		return models.Company{}, er
	}

	return res, nil
}

func (c Company) Update(ctx context.Context, company models.Company) error {
	if er := validate(company); er != nil {
		return er
	}

	if er := c.store.Update(ctx, company); er != nil {
		return er
	}

	return nil
}

func (c Company) Post(ctx context.Context, company models.Company) error {
	if er := validate(company); er != nil {
		return er
	}

	if er := c.store.Post(ctx, company); er != nil {
		return er
	}

	return nil
}

func (c Company) Delete(ctx context.Context, id string) error {
	if er := c.store.Delete(ctx, id); er != nil {
		return er
	}

	return nil
}
