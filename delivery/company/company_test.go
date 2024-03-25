package company

import (
	"StudentPlacement/models"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockDataStore struct{}

func (m mockDataStore) Get(ctx context.Context, id string) (models.Company, error) {
	if id == "10" {
		return models.Company{}, errors.New("Invalid ID")
	}

	return models.Company{ID: "11", Name: "TATA", Category: "MASS"}, nil
}

func (m mockDataStore) Update(ctx context.Context, company models.Company) error {
	if company.Category == "MASS" || company.Category == "DREAM IT" || company.Category == "OPEN DREAM" ||
		company.Category == "CORE" {
		return nil
	}

	return errors.New("Invalid Category")
}

func (m mockDataStore) Post(ctx context.Context, company models.Company) error {
	if company.Category == "MASS" || company.Category == "DREAM IT" || company.Category == "OPEN DREAM" ||
		company.Category == "CORE" {
		return nil
	}

	return errors.New("Invalid Category")
}

func (m mockDataStore) Delete(ctx context.Context, id string) error {
	if id == "19" {
		return errors.New("Invalid ID")
	}

	return nil
}

// nolint
func TestHandler_GetCompany(t *testing.T) {
	type test struct {
		id   string
		code int
	}

	tests := []test{
		{id: "23", code: 200},
		{id: "10", code: 404},
	}

	for _, val := range tests {
		req, er := http.NewRequest("GET", "/employee?id="+val.id, http.NoBody)
		if er != nil {
			log.Println(er)
		}

		rc := httptest.NewRecorder()

		a := New(mockDataStore{})

		a.GetCompany(rc, req)

		if rc.Code != val.code {
			t.Errorf("Expected : %v Got ; %v", val.code, rc.Code)
		}
	}
}

func TestHandler_PostCompany(t *testing.T) {
	type test struct {
		body models.Company
		code int
	}

	tests := []test{
		{body: models.Company{ID: "76", Name: "Infosys", Category: "MASS"}, code: 200},
		{body: models.Company{ID: "77", Name: "TCS", Category: ""}, code: 400},
	}

	for _, val := range tests {
		res, _ := json.Marshal(val.body)
		req, er := http.NewRequest("POST", "/employee", bytes.NewReader(res))

		if er != nil {
			log.Println(er)
		}

		rc := httptest.NewRecorder()

		a := New(mockDataStore{})

		a.PostCompany(rc, req)

		if rc.Code != val.code {
			t.Errorf("Expected : %v Got ; %v", val.code, rc.Code)
		}
	}
}

func TestHandler_UpdateCompany(t *testing.T) {
	type test struct {
		id   string
		body models.Company
		code int
	}

	tests := []test{
		{id: "21", body: models.Company{ID: "76", Name: "Infosys", Category: "MASS"}, code: 200},
		{id: "11", body: models.Company{ID: "77", Name: "TCS", Category: "MAS"}, code: 400},
	}

	for _, val := range tests {
		res, _ := json.Marshal(val.body)
		req, er := http.NewRequest("PUT", "/employee", bytes.NewReader(res))

		if er != nil {
			log.Println(er)
		}

		rc := httptest.NewRecorder()

		a := New(mockDataStore{})
		a.UpdateCompany(rc, req)

		if rc.Code != val.code {
			t.Errorf("Expected : %v Got ; %v", val.code, rc.Code)
		}
	}
}

// nolint
func TestHandler_DeleteCompany(t *testing.T) {
	type test struct {
		id   string
		code int
	}

	tests := []test{
		{id: "89", code: 200},
		{id: "19", code: 404},
	}

	for _, val := range tests {
		req, er := http.NewRequest("DELETE", "/employee?id="+val.id, http.NoBody)
		if er != nil {
			log.Println(er)
		}

		rc := httptest.NewRecorder()

		a := New(mockDataStore{})

		a.DeleteCompany(rc, req)

		if rc.Code != val.code {
			t.Errorf("Expected : %v Got ; %v", val.code, rc.Code)
		}
	}
}
