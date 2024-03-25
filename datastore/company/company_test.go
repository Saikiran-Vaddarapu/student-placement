package company

import (
	customErr "StudentPlacement/customerrors"
	"StudentPlacement/models"
	"context"
	"database/sql"
	"errors"
	"log"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

// nolint : gochecknoglobals //required
var (
	db   *sql.DB
	mock sqlmock.Sqlmock
	err  error
)

func TestDataBase_GetCompany(t *testing.T) {
	db, mock, err = sqlmock.New()

	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "category"}).
		AddRow("89", "MindTree", "MASS")

	mock.ExpectQuery("select id,name,category from company where id = ?").WithArgs("89").WillReturnRows(rows)

	a := New(db)

	got, er := a.Get(context.TODO(), "89")

	expected := models.Company{ID: "89", Name: "MindTree", Category: "MASS"}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected : %v, Got : %v %v", expected, got, er)
	}

	mock.ExpectQuery("select id,name,category, from company where id = ?").
		WithArgs("98").WillReturnError(errors.New("Invalid ID"))

	got, er = a.Get(context.TODO(), "98")

	expected = models.Company{}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected : %v, Got : %v %v", expected, got, er)
	}
}

func TestDataBase_DeleteCompany(t *testing.T) {
	db, mock, err = sqlmock.New()

	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()

	mock.ExpectExec("delete from company where id = ?").WithArgs("88").
		WillReturnError(errors.New("Invalid ID"))

	a := New(db)

	er := a.Delete(context.TODO(), "88")

	expect := customErr.InternalErrors{Message: "Invalid ID"}

	if !reflect.DeepEqual(er, expect) {
		t.Errorf("Expected : %v, Got : %v", expect, er)
	}

	mock.ExpectExec("delete from companyy where id = ?").WithArgs("98").
		WillReturnError(errors.New("ExecQuery: could not match actual sql: \\\"delete from companyy where id " +
			"= ?\\\" with\" +\n\t\t\" expected regexp \\\"delete from company where id = ?\\\""))

	er = a.Delete(context.TODO(), "98")

	expect = customErr.InternalErrors{Message: "ExecQuery: could not match actual sql: \"delete from company where id = ?\" with" +
		" expected regexp \"delete from companyy where id = ?\""}

	if !reflect.DeepEqual(er, expect) {
		t.Errorf("Expected : %v, Got : %v", expect, er)
	}
}

func TestDataBase_PostCompany(t *testing.T) {
	company := models.Company{Name: "Quantiphi", Category: "DREAM IT"}

	db, mock, err = sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()

	mock.ExpectExec("insert into company values(uuid(),?,?)").
		WithArgs("Quantiphi", "DREAM IT").WillReturnResult(sqlmock.NewResult(1, 1))

	a := New(db)

	er := a.Post(context.TODO(), company)

	if er != nil {
		t.Errorf("Expected : %v Got : %v", nil, er)
	}

	mock.ExpectExec("insert into company values(uuid(),,?)").
		WithArgs("Quantiphi", "DREAM IT").WillReturnError(errors.New("Syntax error"))

	er = a.Post(context.TODO(), company)

	if reflect.DeepEqual(er, customErr.InternalErrors{Message: " ExecQuery: actual sql: \"insert into Company values(uuid(),?,?)\"" +
		" does not equal to expected \"insert into Company values(uuid(),,?)"}) {
		t.Errorf("Expected : %v Got : %v", nil, er)
	}
}

func TestDataBase_UpdateCompany(t *testing.T) {
	company := models.Company{ID: "89", Name: "HCL", Category: "MASS"}

	db, mock, err = sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "category"}).
		AddRow("89", "HCL", "MASS")

	mock.ExpectQuery("select id,name,category from company where id = ?").WithArgs("89").
		WillReturnRows(rows)

	mock.ExpectExec("update company set name = ?,category = ? where id = ?").
		WithArgs("HCL", "MASS", "89").WillReturnResult(sqlmock.NewResult(1, 1))

	a := New(db)

	er := a.Update(context.TODO(), company)

	if er != nil {
		t.Errorf("Expected : %v, Got : %v", nil, er)
	}
}
