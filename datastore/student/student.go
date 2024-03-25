package student

import (
	customErr "StudentPlacement/customerrors"
	"StudentPlacement/models"
	"context"
	"database/sql"
)

type DataBase struct {
	db *sql.DB
}

func New(db *sql.DB) DataBase {
	return DataBase{db: db}
}

// nolint
func (d DataBase) ValidateCategoryBranch(student *models.Student) error {
	id := student.Company.ID

	var company models.Company

	row := d.db.QueryRow("select category from company where id = ?", id)

	if er := row.Scan(&company.Category); er != nil {
		return customErr.InternalErrors{Message: er.Error()}
	}

	category := company.Category

	if category == models.Core && !(student.Branch == models.CivilBranch || student.Branch ==
		models.MechBranch) {
		return customErr.Validity{Message: "student Branch and Company Category Mismapped"}
	}

	if category == models.DreamIT && !(student.Branch == models.CSEBranch || student.Branch ==
		models.ISEBranch) {
		return customErr.Validity{Message: "student Branch and Company Category Mismapped"}
	}

	if category == models.OpenDream && !(student.Branch == models.CSEBranch || student.Branch ==
		models.ISEBranch || student.Branch == models.EceBranch || student.Branch == models.EeeBranch) {
		return customErr.Validity{Message: "student Branch and Company Category Mismapped"}
	}

	return nil
}

func (d DataBase) GetAllStudentsWithCompany(ctx context.Context, branch, name string) ([]models.Student, error) {
	rows, er := d.db.QueryContext(ctx, "select s.id,s.name,s.phone,s.branch,s.dob,s.status,s.companyid,c.name,c.category"+
		" from student s join company c on s.companyid = c.id where s.name = ? and s.branch = ?", name, branch)
	if er != nil {
		return nil, customErr.InternalErrors{Message: er.Error()}
	}

	var students []models.Student

	for rows.Next() {
		var student models.Student
		if err := rows.Scan(&student.ID, &student.Name, &student.Phone, &student.Branch, &student.DOB, &student.Status,
			&student.Company.ID, &student.Company.Name, &student.Company.Category); err != nil {
			return nil, customErr.InternalErrors{Message: er.Error()}
		}

		students = append(students, student)
	}

	return students, nil
}

func (d DataBase) GetAll(ctx context.Context, name, branch string) ([]models.Student, error) {
	rows, er := d.db.QueryContext(ctx, "select id,name,phone,branch,dob,status,companyid from student "+
		"where name = ? and branch = ?", name, branch)
	if er != nil {
		return nil, customErr.InternalErrors{Message: er.Error()}
	}

	var students []models.Student

	for rows.Next() {
		var student models.Student
		if err := rows.Scan(&student.ID, &student.Name, &student.Phone, &student.Branch, &student.DOB, &student.Status,
			&student.Company.ID); err != nil {
			return nil, customErr.InternalErrors{Message: er.Error()}
		}

		students = append(students, student)
	}

	return students, nil
}

func (d DataBase) Get(ctx context.Context, id string) (models.Student, error) {
	rows := d.db.QueryRowContext(ctx, "select id,name,phone,branch,dob,status,companyid from student where id = ?", id)

	var student models.Student

	er := rows.Scan(&student.ID, &student.Name, &student.Phone, &student.Branch, &student.DOB,
		&student.Status, &student.Company.ID)

	if er != nil {
		return models.Student{}, customErr.InternalErrors{Message: er.Error()}
	}

	return student, nil
}

func (d DataBase) Delete(ctx context.Context, id string) error {
	_, er := d.db.ExecContext(ctx, "delete from student where id = ? ", id)
	if er != nil {
		return customErr.InternalErrors{Message: er.Error()}
	}

	return nil
}

func (d DataBase) Post(ctx context.Context, student *models.Student) error {
	_, er := d.db.ExecContext(ctx, "insert into student values(uuid(),?,?,?,?,?,?)", student.Name, student.Phone,
		student.Branch, student.DOB, student.Status, student.Company.ID)
	if er != nil {
		return customErr.InternalErrors{Message: er.Error()}
	}

	return nil
}

func (d DataBase) Update(ctx context.Context, student *models.Student) error {
	_, er := d.Get(ctx, student.ID)

	if er != nil {
		return customErr.InternalErrors{Message: er.Error()}
	}

	_, err := d.db.ExecContext(ctx, "update student set name = ?,phone=?,branch=?,dob=?,status=?,companyid=? where id = ?",
		student.Name, student.Phone, student.Branch, student.DOB, student.Status, student.Company.ID, student.ID)

	if err != nil {
		return customErr.InternalErrors{Message: er.Error()}
	}

	return nil
}
