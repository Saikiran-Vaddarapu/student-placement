package student

import (
	customErr "StudentPlacement/customerrors"
	"StudentPlacement/models"
	"context"
	"database/sql"
	"log"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

// nolint : gochecknoglobals // required
var (
	db   *sql.DB
	mock sqlmock.Sqlmock
	err  error
)

func TestDataBase_ValidateCategoryBranch(t *testing.T) {
	tests := []struct {
		student models.Student
		er      error
	}{
		{models.Student{ID: "12", Name: "Arya", Phone: "8798745321", Branch: "EEE",
			DOB: "12-07-1998", Status: "REJECTED", Company: models.Company{ID: "12", Name: "EY", Category: "CORE"}},
			customErr.Validity{Message: "student Branch and Company Category Mismapped"},
		},
		{models.Student{ID: "13", Name: "Arya", Phone: "8798745321123", Branch: "ECE",
			DOB: "12-07-1998", Status: "REJECTED", Company: models.Company{ID: "19", Name: "EY", Category: "DREAM IT"}},
			customErr.Validity{Message: "student Branch and Company Category Mismapped"},
		},
		{models.Student{ID: "13", Name: "Arya", Phone: "8798745321123", Branch: "MECH",
			DOB: "12-07-1998", Status: "REJECTED",
			Company: models.Company{ID: "19", Name: "EY", Category: "OPEN DREAM"}},
			customErr.Validity{Message: "student Branch and Company Category Mismapped"},
		},
		{models.Student{ID: "13", Name: "Arya", Phone: "8798745321123", Branch: "CIVIL",
			DOB: "12-07-1998", Status: "REJECTED", Company: models.Company{ID: "19", Name: "EY", Category: "MASS"}},
			nil},
	}
	for _, val := range tests {
		db, mock, err = sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		if err != nil {
			log.Println(err)
			return
		}

		rows := sqlmock.NewRows([]string{"category"}).
			AddRow(val.student.Company.Category)

		mock.ExpectQuery("select category from company where id = ?").WithArgs(val.student.Company.ID).WillReturnRows(rows)

		a := New(db)

		er := a.ValidateCategoryBranch(&val.student)
		if !reflect.DeepEqual(er, val.er) {
			t.Errorf("Expected : %v, Got : %v", val.er, er)
		}
	}
}

func TestDataBase_GetAllStudentsWithComapany(t *testing.T) {
	db, mock, err = sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "phone", "branch", "dob", "status", "companyid", "name", "category"}).
		AddRow("21", "Sarath", "5467890321", "CSE", "2000-08-21", "PENDING", "23", "McAfee", "DREAM IT")

	mock.ExpectQuery("select s.id,s.name,s.phone,s.branch,s.dob,s.status,s.companyid,c.name,c.category from student s "+
		"join company c on s.companyid = c.id where s.name = ? and s.branch = ?").
		WithArgs("Sarath", "CSE").WillReturnRows(rows)

	a := New(db)
	st, _ := a.GetAllStudentsWithCompany(context.TODO(), "CSE", "Sarath")
	expected := []models.Student{{ID: "21", Name: "Sarath", Phone: "5467890321", Branch: "CSE", DOB: "2000-08-21",
		Status: "PENDING", Company: models.Company{ID: "23", Name: "McAfee", Category: "DREAM IT"}}}

	if !reflect.DeepEqual(st, expected) {
		t.Errorf("Expected : %v, Got : %v", expected, st)
	}
}

func TestDataBase_GetAllStudents(t *testing.T) {
	db, mock, err = sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "phone", "branch", "dob", "status", "companyid"}).
		AddRow("21", "Sarath", "5467890321", "CSE", "2000-08-21", "PENDING", "23")

	mock.ExpectQuery("select id,name,phone,branch,dob,status,companyid from student where name = ? and branch = ?").
		WithArgs("Sarath", "CSE").WillReturnRows(rows)

	a := New(db)
	st, _ := a.GetAll(context.TODO(), "Sarath", "CSE")
	expected := []models.Student{{ID: "21", Name: "Sarath", Phone: "5467890321", Branch: "CSE", DOB: "2000-08-21",
		Status: "PENDING", Company: models.Company{ID: "23", Name: "", Category: ""}}}

	if !reflect.DeepEqual(st, expected) {
		t.Errorf("Expected : %v, Got : %v", expected, st)
	}
}

func TestDataBase_GetStudent(t *testing.T) {
	db, mock, err = sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "phone", "branch", "dob", "status", "companyid"}).
		AddRow("21", "Sarath", "5467890321", "CSE", "2000-08-21", "PENDING", "23")

	mock.ExpectQuery("select id,name,phone,branch,dob,status,companyid from student where id = ?").
		WithArgs("21").WillReturnRows(rows)

	if err != nil {
		t.Fatal(err)
	}

	a := New(db)
	st, _ := a.Get(context.TODO(), "21")
	expected := models.Student{ID: "21", Name: "Sarath", Phone: "5467890321", Branch: "CSE", DOB: "2000-08-21",
		Status: "PENDING", Company: models.Company{ID: "23", Name: "", Category: ""}}

	if !reflect.DeepEqual(st, expected) {
		t.Errorf("Expected : %v, Got : %v", expected, st)
	}
}

func TestDataBase_DeleteStudent(t *testing.T) {
	db, mock, err = sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()
	mock.ExpectExec("delete from student where id = ?").WithArgs("21").WillReturnError(customErr.InternalErrors{Message: "Invalid ID"})
	a := New(db)
	er := a.Delete(context.TODO(), "21")
	expect := customErr.InternalErrors{Message: "Invalid ID"}

	if !reflect.DeepEqual(er, expect) {
		t.Errorf("Expected : %v, Got : %v", expect, er)
	}
}

func TestDataBase_PostStudent(t *testing.T) {
	st := models.Student{ID: "", Name: "Ramu", Phone: "7896578921", Branch: "EEE",
		DOB: "2000-10-22", Status: "ACCEPTED", Company: models.Company{ID: "876", Name: "", Category: ""}}
	db, mock, err = sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()
	mock.ExpectExec("insert into student values(uuid(),?,?,?,?,?,?)").
		WithArgs("Ramu", "7896578921", "EEE", "2000-10-22", "ACCEPTED", "876").WillReturnResult(sqlmock.NewResult(1, 1))

	a := New(db)
	er := a.Post(context.TODO(), &st)

	if er != nil {
		t.Errorf("Expected : %v, Got : %v", nil, er)
	}
}

func TestDataBase_UpdateStudent(t *testing.T) {
	st := models.Student{ID: "32", Name: "Ramu", Phone: "7896578921", Branch: "EEE",
		DOB: "2000-10-22", Status: "ACCEPTED", Company: models.Company{ID: "876", Name: "", Category: ""}}
	db, mock, err = sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "phone", "branch", "dob", "status", "companyid"}).
		AddRow("32", "Ramu", "7896578921", "EEE", "2000-10-22", "ACCEPTED", "876")

	mock.ExpectQuery("select id,name,phone,branch,dob,status,companyid from student where id = ?").
		WithArgs("32").WillReturnRows(rows)
	mock.ExpectExec("update student set name = ?,phone=?,branch=?,dob=?,status=?,companyid=? where id = ?").
		WithArgs("Ramu", "7896578921", "EEE", "2000-10-22", "ACCEPTED", "876", "32").WillReturnResult(sqlmock.NewResult(1, 1))

	a := New(db)
	er := a.Update(context.TODO(), &st)

	if er != nil {
		t.Errorf("Expected : %v, Got : %v", nil, er)
	}
}
