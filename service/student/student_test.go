package student

import (
	customErr "StudentPlacement/customerrors"
	"StudentPlacement/models"
	"context"
	"errors"
	"reflect"
	"testing"
)

func TestGetAge(t *testing.T) {
	tests := []struct {
		input  string
		output int
		err    error
	}{
		{input: "12-09-2000", output: 22, err: nil},
		{input: "18-11-1981", output: 41, err: nil},
		{input: "28-02-2021", output: 1, err: nil},
		{input: "28-02-20o1", output: 0, err: errors.New("strconv.Atoi: parsing \"20o1\": invalid syntax ")},
	}
	for _, val := range tests {
		got, er := getAge(val.input)

		if er != nil && val.output != 0 {
			t.Errorf("Expected : %v Got : %v", val.err, er)
			return
		}

		if got != val.output {
			t.Errorf("Expected : %v Got : %v", val.output, val.input)
		}
	}
}

func TestValidateStudent(t *testing.T) {
	type test struct {
		input models.Student
		er    error
	}

	tests := []test{
		{input: models.Student{ID: "11", Name: "VS", Phone: "8760965543", Branch: "ECE", DOB: "12-08-1998",
			Status: "ACCEPTED", Company: models.Company{ID: "98", Name: "TCS", Category: "MASS"}},
			er: customErr.Validity{Message: "name should be minimum 3 characters"}},
		{input: models.Student{ID: "14", Name: "Namata", Phone: "8760965", Branch: "ECE", DOB: "12-07-1998",
			Status: "ACCEPTED", Company: models.Company{ID: "98", Name: "TCS", Category: "MASS"}},
			er: customErr.Validity{Message: "phone must be 10 -12 digits"}},
		{input: models.Student{ID: "19", Name: "Helen", Phone: "8760965987", Branch: "ECE", DOB: "13-08-1998",
			Status: "ACCEPTED", Company: models.Company{ID: "98", Name: "TCS", Category: "MASS"}}},
		{input: models.Student{ID: "91", Name: "Nasar", Phone: "9874563213", Branch: "IT", DOB: "13-08-1998",
			Status: "ACCEPTED", Company: models.Company{ID: "98", Name: "TCS", Category: "MASS"}},
			er: customErr.Validity{Message: "invalid Branch"}},
		{input: models.Student{ID: "46", Name: "Teena", Phone: "6798543218", Branch: "ECE", DOB: "13-08-1998",
			Status: "INREVIEW", Company: models.Company{ID: "98", Name: "TCS", Category: "MASS"}},
			er: customErr.Validity{Message: "invalid Status"}},
		{input: models.Student{ID: "23", Name: "Radha", Phone: "8760965987", Branch: "ECE", DOB: "13-08-2001",
			Status: "ACCEPTED", Company: models.Company{ID: "98", Name: "TCS", Category: "MASS"}},
			er: customErr.Validity{Message: "minimum Age Required : 22"}},
	}
	for _, val := range tests {
		res := validate(&val.input)

		if !reflect.DeepEqual(res, val.er) {
			t.Errorf("Expected : %v Got : %v ", val.er.Error(), res.Error())
		}
	}
}

type mockdatastore struct{}

func (m mockdatastore) GetAllStudentsWithCompany(ctx context.Context, branch, name string) ([]models.Student,
	error) {
	return []models.Student{{ID: "12", Name: "Ramesh", Phone: "9876094532", Branch: "ECE",
		DOB: "12-07-1999", Status: "REJECTED", Company: models.Company{ID: "13", Name: "LandT", Category: "MASS"}}}, nil
}

func (m mockdatastore) GetAll(ctx context.Context, name, branch string) ([]models.Student, error) {
	return []models.Student{{ID: "11", Name: "Arya", Phone: "8798745321", Branch: "EEE", DOB: "12-07-1998",
		Status:  "REJECTED",
		Company: models.Company{ID: "12"}}}, nil
}

func TestStore_GetAllStudents(t *testing.T) {
	tests := []struct {
		name           string
		branch         string
		includecompany string
		resp           []models.Student
		er             error
	}{
		{name: "Arya", branch: "EEE", includecompany: "false", resp: []models.Student{{ID: "11", Name: "Arya",
			Phone: "8798745321", Branch: "EEE",
			DOB: "12-07-1998", Status: "REJECTED", Company: models.Company{ID: "12"}}}},
		{name: "Ramesh", branch: "ECE", includecompany: "true", resp: []models.Student{{ID: "12", Name: "Ramesh",
			Phone: "9876094532", Branch: "ECE",
			DOB: "12-07-1999", Status: "REJECTED", Company: models.Company{ID: "13", Name: "LandT", Category: "MASS"}}}},
		{name: "Helen", includecompany: "false", er: customErr.Validity{Message: "invalid Branch"}},
		{branch: "MECH", includecompany: "false", er: customErr.Validity{Message: "invalid Name"}},
	}
	for _, val := range tests {
		a := New(mockdatastore{})

		res, er := a.GetAll(context.TODO(), val.name, val.branch, val.includecompany)

		if er != nil && !reflect.DeepEqual(er, val.er) {
			t.Errorf("Expected : %v Got : %v", val.er, er)
			return
		}

		if !reflect.DeepEqual(res, val.resp) {
			t.Errorf("Expected : %v Got : %v", val.resp, res)
		}
	}
}

func (m mockdatastore) Get(ctx context.Context, id string) (models.Student, error) {
	if id == "11" {
		return models.Student{ID: "11", Name: "Arya", Phone: "8798745321", Branch: "EEE", DOB: "12-07-1998",
			Status:  "REJECTED",
			Company: models.Company{ID: "12", Name: "EY", Category: "DREAM IT"}}, nil
	}

	return models.Student{}, customErr.Validity{Message: "invalid id"}
}

func TestService_GetStudent(t *testing.T) {
	tests := []struct {
		id   string
		resp models.Student
		er   error
	}{
		{id: "11", resp: models.Student{ID: "11", Name: "Arya", Phone: "8798745321", Branch: "EEE",
			DOB: "12-07-1998", Status: "REJECTED", Company: models.Company{ID: "12", Name: "EY", Category: "DREAM IT"}}},
		{id: "12", er: customErr.Validity{Message: "invalid id"}},
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
	if id == "11" {
		return nil
	}

	return customErr.Validity{Message: "invalid id"}
}

func TestService_DeleteStudent(t *testing.T) {
	type test struct {
		id string
		er error
	}

	tests := []test{
		{id: "11", er: nil},
		{id: "12", er: customErr.Validity{Message: "invalid id"}},
	}
	for _, val := range tests {
		a := New(mockdatastore{})

		er := a.Delete(context.TODO(), val.id)

		if !reflect.DeepEqual(er, val.er) {
			t.Errorf("Expected : %v Got : %v", val.er, er)
		}
	}
}

func (m mockdatastore) Update(ctx context.Context, student *models.Student) error {
	if er := validate(student); er != nil {
		return er
	}

	if student.ID == "11" {
		return nil
	}

	return customErr.Validity{Message: "invalid ID"}
}

func TestService_UpdateStudent(t *testing.T) {
	type test struct {
		resp models.Student
		er   error
	}

	tests := []test{
		{resp: models.Student{ID: "11", Name: "Arya", Phone: "8798745321", Branch: "EEE",
			DOB: "12-07-1998", Status: "REJECTED", Company: models.Company{ID: "12", Name: "EY", Category: "DREAM IT"}}},
		{resp: models.Student{ID: "13", Name: "Arya", Phone: "8798745321123", Branch: "EEE",
			DOB: "12-07-1998", Status: "REJECTED", Company: models.Company{ID: "19", Name: "EY", Category: "DREAM IT"}},
			er: customErr.Validity{Message: "phone must be 10 -12 digits"}},
	}
	for _, val := range tests {
		a := New(mockdatastore{})

		er := a.Update(context.TODO(), &val.resp)

		if !reflect.DeepEqual(er, val.er) {
			t.Errorf("Expected : %v Got : %v", val.er, er)
			return
		}
	}
}

func (m mockdatastore) ValidateCategoryBranch(student *models.Student) error {
	if student.Company.ID == "18" {
		return customErr.Validity{Message: "student Branch and Company Category Mismapped"}
	}

	return nil
}

func (m mockdatastore) Post(ctx context.Context, student *models.Student) error {
	return nil
}

func TestService_PostStudent(t *testing.T) {
	type test struct {
		resp models.Student
		er   error
	}

	tests := []test{
		{resp: models.Student{ID: "18", Name: "Arya", Phone: "8798745321", Branch: "EEE",
			DOB: "12-07-1998", Status: "REJECTED", Company: models.Company{ID: "12", Name: "EY", Category: "DREAM IT"}}},
		{resp: models.Student{ID: "13", Name: "Arya", Phone: "8798745321123", Branch: "EEE",
			DOB: "12-07-1998", Status: "REJECTED", Company: models.Company{ID: "19", Name: "EY", Category: "DREAM IT"}},
			er: customErr.Validity{Message: "phone must be 10 -12 digits"}},
		{resp: models.Student{ID: "19", Name: "Arya", Phone: "8798745321", Branch: "CSE",
			DOB: "12-07-1998", Status: "REJECTED", Company: models.Company{ID: "18", Name: "TASL", Category: "CORE"}},
			er: customErr.Validity{Message: "student Branch and Company Category Mismapped"}},
	}
	for _, val := range tests {
		a := New(mockdatastore{})

		er := a.Post(context.TODO(), &val.resp)

		if !reflect.DeepEqual(er, val.er) {
			t.Errorf("Expected : %v Got : %v", val.er, er)
			return
		}
	}
}
