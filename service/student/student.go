package student

import (
	customErr "StudentPlacement/customerrors"
	"StudentPlacement/datastore"
	"StudentPlacement/models"
	"context"
	"strconv"
	"time"
)

type Student struct {
	store datastore.Student
}

func New(st datastore.Student) Student {
	return Student{store: st}
}

func getAge(date string) (int, error) {
	currYear := time.Now().Year()

	birthYear, err := strconv.Atoi(date[6:])

	if err != nil {
		return 0, err
	}

	return currYear - birthYear, nil
}

func validate(student *models.Student) error {
	if len(student.Name) < models.MinimumNameLength {
		return customErr.Validity{Message: "name should be minimum 3 characters"}
	}

	if len(student.Phone) < models.MinimumPhoneLength || len(student.Phone) > models.MaximumPhoneLength {
		return customErr.Validity{Message: "phone must be 10 -12 digits"}
	}

	if !models.Contains(models.ValidBranch, student.Branch) {
		return customErr.Validity{Message: "invalid Branch"}
	}

	if !models.Contains(models.ValidStudentStatus, student.Status) {
		return customErr.Validity{Message: "invalid Status"}
	}

	age, er := getAge(student.DOB)

	if er != nil {
		return er
	}

	if age < models.MinimumAge {
		return customErr.Validity{Message: "minimum Age Required : 22"}
	}

	return nil
}

func (s Student) GetAll(ctx context.Context, name, branch, company string) ([]models.Student, error) {
	if name == "" {
		return nil, customErr.Validity{Message: "invalid Name"}
	}

	if branch == "" {
		return nil, customErr.Validity{Message: "invalid Branch"}
	}

	var (
		resp []models.Student
		er   error
	)

	if company == "true" {
		resp, er = s.store.GetAllStudentsWithCompany(ctx, branch, name)

		if er != nil {
			return nil, er
		}

		return resp, nil
	}

	resp, er = s.store.GetAll(ctx, name, branch)

	if er != nil {
		return nil, er
	}

	return resp, nil
}

func (s Student) Get(ctx context.Context, id string) (models.Student, error) {
	resp, er := s.store.Get(ctx, id)

	if er != nil {
		return models.Student{}, er
	}

	return resp, nil
}

func (s Student) Post(ctx context.Context, student *models.Student) error {
	if err := validate(student); err != nil {
		return err
	}

	if err := s.store.ValidateCategoryBranch(student); err != nil {
		return err
	}

	er := s.store.Post(ctx, student)

	if er != nil {
		return er
	}

	return nil
}

func (s Student) Update(ctx context.Context, student *models.Student) error {
	if err := validate(student); err != nil {
		return err
	}

	if err := s.store.ValidateCategoryBranch(student); err != nil {
		return err
	}

	er := s.store.Update(ctx, student)

	if er != nil {
		return er
	}

	return nil
}

func (s Student) Delete(ctx context.Context, id string) error {
	er := s.store.Delete(ctx, id)

	if er != nil {
		return er
	}

	return nil
}
