package company

import (
	customErr "StudentPlacement/customerrors"
	"StudentPlacement/models"
	"context"
	"errors"
	"reflect"
	"testing"
)

func TestValidateCompany(t *testing.T) {
	type test struct {
		company models.Company
		er      error
	}

	tests := []test{
		{company: models.Company{ID: "65", Name: "Hyundai", Category: "OPEN DREAM"}, er: nil},
		{company: models.Company{ID: "78", Name: "TVS", Category: "CORE SERVICE"},
			er: &customErr.Validity{Message: "invalid Category"}},
	}
	for _, val := range tests {
		er := validate(val.company)

		if !reflect.DeepEqual(er, val.er) {
			t.Errorf("Expected : %v Got : %v", val.er, er)
		}
	}
}

type mockdatastore struct{}

func (m mockdatastore) Get(ctx context.Context, id string) (models.Company, error) {
	if id == "67" {
		return models.Company{ID: "32", Name: "Infosys", Category: "MASS"}, nil
	}

	return models.Company{}, errors.New("Invalid ID")
}

func TestService_GetCompany(t *testing.T) {
	type test struct {
		id   string
		resp models.Company
		er   error
	}

	tests := []test{
		{id: "67", resp: models.Company{ID: "32", Name: "Infosys", Category: "MASS"},
			er: nil},
		{id: "69", resp: models.Company{}, er: errors.New("Invalid ID")},
	}
	for _, val := range tests {
		a := New(mockdatastore{})

		res, er := a.Get(context.TODO(), val.id)

		if er != nil && !reflect.DeepEqual(er, val.er) {
			t.Errorf("Expected : %v Got : %v", val.er, er)

			return
		}

		if !reflect.DeepEqual(res, val.resp) {
			t.Errorf("Expected : %v Got : %v", val.resp, res)
		}
	}
}

func (m mockdatastore) Delete(ctx context.Context, id string) error {
	if id == "69" {
		return errors.New("Invalid ID")
	}

	return nil
}

func TestService_DeleteCompany(t *testing.T) {
	type test struct {
		id string
		er error
	}

	tests := []test{
		{id: "67", er: nil},
		{id: "69", er: errors.New("Invalid ID")},
	}
	for _, val := range tests {
		a := New(mockdatastore{})

		er := a.Delete(context.TODO(), val.id)

		if !reflect.DeepEqual(er, val.er) {
			t.Errorf("Expected : %v Got : %v", val.er, er)
		}
	}
}

func (m mockdatastore) Post(ctx context.Context, company models.Company) error {
	return nil
}

// nolint
func TestService_PostCompany(t *testing.T) {
	type test struct {
		resp models.Company
		er   error
	}

	tests := []test{
		{resp: models.Company{ID: "32", Name: "Infosys", Category: "MASS"},
			er: nil},
		{resp: models.Company{ID: "72", Name: "Tricon", Category: "MASS IT"},
			er: &customErr.Validity{Message: "invalid Category"}},
	}
	for _, val := range tests {
		a := New(mockdatastore{})

		er := a.Post(context.TODO(), val.resp)

		if !reflect.DeepEqual(er, val.er) {
			t.Errorf("Expected : %v Got : %v", val.er, er)
		}
	}
}

func (m mockdatastore) Update(ctx context.Context, company models.Company) error {
	return nil
}

// nolint
func TestService_UpdateCompany(t *testing.T) {
	type test struct {
		resp models.Company
		er   error
	}

	tests := []test{
		{resp: models.Company{ID: "43", Name: "META", Category: "DREAM IT"},
			er: nil},
		{resp: models.Company{ID: "98", Name: "Nutanix", Category: "DREAM BIG"},
			er: &customErr.Validity{Message: "invalid Category"}},
	}
	for _, val := range tests {
		a := New(mockdatastore{})

		er := a.Update(context.TODO(), val.resp)

		if !reflect.DeepEqual(er, val.er) {
			t.Errorf("Expected : %v Got : %v", val.er, er)
		}
	}
}
