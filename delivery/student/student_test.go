package student

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

func TestHandler_GetAllStudents(t *testing.T) {
	type test struct {
		id   string
		resp models.Student
		code int
	}

	tests := []test{
		{id: "11",
			resp: models.Student{ID: "11", Name: "Arun", Phone: "8768890981", Branch: "CSE", DOB: "12-12-1999",
				Status: "ACCEPTED", Company: models.Company{ID: "98", Name: "TCS", Category: "MASS"}},
			code: 200},
	}

	for _, val := range tests {
		req, er := http.NewRequest("GET", "/allstudents", http.NoBody)

		if er != nil {
			log.Println(er)
		}

		rc := httptest.NewRecorder()

		a := New(mockDataStore{})

		a.GetAllStudents(rc, req)

		if rc.Code != val.code {
			t.Errorf("Expected : %v Got : %v ", val.code, rc.Code)
		}
	}
}

func TestHandler_GetById(t *testing.T) {
	type test struct {
		id   string
		resp models.Student
		code int
	}

	tests := []test{
		{id: "11", resp: models.Student{ID: "11", Name: "Arun", Phone: "8768890981", Branch: "CSE", DOB: "12-12-1999",
			Status: "ACCEPTED", Company: models.Company{ID: "98", Name: "TCS", Category: "MASS"}}, code: 200},
		{id: "1", code: 404},
	}

	for _, val := range tests {
		req, er := http.NewRequest("GET", "/student?id="+val.id, http.NoBody)

		if er != nil {
			log.Println(er)
		}

		rc := httptest.NewRecorder()
		a := New(mockDataStore{})

		a.GetStudentByID(rc, req)

		if rc.Code != val.code {
			t.Errorf("Expected : %v Got : %v ", val.code, rc.Code)

			return
		}
	}
}

func TestHandler_Post(t *testing.T) {
	type test struct {
		body       models.Student
		statuscode int
	}

	tests := []test{
		{body: models.Student{ID: "22", Name: "Raju", Phone: "8887765431", Branch: "CSE",
			DOB: "16-09-1999", Status: "PENDING",
			Company: models.Company{ID: "89", Name: "Deloitte", Category: "DREAM IT"}},
			statuscode: 200},
	}

	for _, val := range tests {
		res, _ := json.Marshal(val.body)
		req, er := http.NewRequest("POST", "/student", bytes.NewReader(res))

		if er != nil {
			log.Println(er)
		}

		rc := httptest.NewRecorder()
		a := New(mockDataStore{})

		a.PostStudent(rc, req)

		if rc.Code != val.statuscode {
			t.Errorf("Expected : %v Got : %v ", val.statuscode, rc.Code)
		}
	}
}

func TestHandler_Update(t *testing.T) {
	type test struct {
		id         string
		body       models.Student
		statuscode int
	}

	tests := []test{
		{id: "12", body: models.Student{ID: "12", Name: "SAMY", Phone: "8764532091", Branch: "CIVIL",
			DOB: "12-07-1998", Status: "ACCEPTED",
			Company: models.Company{ID: "98", Name: "TCS", Category: "MASS"}},
			statuscode: 200},
		{id: "76", body: models.Student{ID: "76", Name: "HEMA", Phone: "9125437690", Branch: "MECH",
			DOB: "19-07-1997", Status: "PENDING",
			Company: models.Company{ID: "101", Name: "AMAZON", Category: "DREAMIT"}},
			statuscode: 400},
	}

	for _, val := range tests {
		resp, _ := json.Marshal(val.body)
		req, er := http.NewRequest("PUT", "/student", bytes.NewReader(resp))

		if er != nil {
			log.Println(er)
		}

		rc := httptest.NewRecorder()
		a := New(mockDataStore{})

		a.UpdateStudent(rc, req)

		if rc.Code != val.statuscode {
			t.Errorf("Expected : %v Got : %v ", val.statuscode, rc.Code)
		}
	}
}

func TestHandler_Delete(t *testing.T) {
	type test struct {
		url  string
		code int
	}

	tests := []test{
		{url: "/student?id=11", code: 200},
		{url: "/student?id=2", code: 404},
		{url: "/student", code: 404},
	}

	for _, val := range tests {
		req, er := http.NewRequest("DELETE", val.url, http.NoBody)
		if er != nil {
			log.Println(er)
		}

		rc := httptest.NewRecorder()
		a := New(mockDataStore{})

		a.DeleteStudent(rc, req)

		if rc.Code != val.code {
			t.Errorf("Expected : %v Got : %v ", val.code, rc.Code)
		}
	}
}

func (m mockDataStore) GetAll(ctx context.Context, name, branch, flag string) ([]models.Student, error) {
	return []models.Student{
		{ID: "11", Name: "Arun", Phone: "8768890981", Branch: "CSE", DOB: "12-12-1999",
			Status: "ACCEPTED", Company: models.Company{ID: "98", Name: "TCS", Category: "MASS"}}}, nil
}

func (m mockDataStore) Get(ctx context.Context, id string) (models.Student, error) {
	if id == "11" {
		return models.Student{ID: "11", Name: "Arun", Phone: "8768890981", Branch: "CSE", DOB: "12-12-1999",
			Status: "ACCEPTED", Company: models.Company{ID: "98", Name: "TCS", Category: "MASS"}}, nil
	}

	return models.Student{}, errors.New("Invalid ID")
}

func (m mockDataStore) Post(ctx context.Context, student *models.Student) error {
	if len(student.Name) < 3 {
		return errors.New("Name should be atleast 3 characters")
	}

	if len(student.Phone) < 10 || len(student.Phone) > 12 {
		return errors.New("Phone must be 10 - 12 diits")
	}

	return nil
}

func (m mockDataStore) Delete(ctx context.Context, id string) error {
	if id == "11" {
		return nil
	}

	return errors.New("Invalid ID")
}

func (m mockDataStore) Update(ctx context.Context, student *models.Student) error {
	if student.ID == "76" {
		return errors.New("Invalid ID")
	}

	if len(student.Name) < 3 {
		return errors.New("Name should be atleast 3 characters")
	}

	if len(student.Phone) < 10 || len(student.Phone) > 12 {
		return errors.New("Phone must be 10 - 12 diits")
	}

	return nil
}
