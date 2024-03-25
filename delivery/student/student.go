package student

import (
	"StudentPlacement/models"
	"StudentPlacement/service"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type student struct {
	service service.Student
}

func New(studentservice service.Student) student {
	return student{service: studentservice}
}

func (s student) GetAllStudents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ctx := r.Context()
	name := r.URL.Query().Get("name")
	branch := r.URL.Query().Get("branch")
	includeFlag := r.URL.Query().Get("includeCompany")
	res, er := s.service.GetAll(ctx, name, branch, includeFlag)

	if er != nil {
		w.WriteHeader(http.StatusNotFound)

		if _, er = w.Write([]byte(er.Error())); er != nil {
			log.Println(er)
		}

		return
	}

	resp, err := json.Marshal(res)

	if err != nil {
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusOK)

	if _, er = w.Write(resp); er != nil {
		log.Println(er)
	}
}

func (s student) GetStudentByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx := r.Context()

	id := r.URL.Query().Get("id")

	res, er := s.service.Get(ctx, id)

	if er != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	resp, er := json.Marshal(res)
	if er != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Error in marshaling the response : %v", er)

		return
	}

	w.WriteHeader(http.StatusOK)

	if _, er = w.Write(resp); er != nil {
		log.Println(er)
	}
}

func (s student) PostStudent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ctx := r.Context()

	resp, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(w, "Error in Reading Request", http.StatusBadRequest)
		return
	}

	var st models.Student

	err = json.Unmarshal(resp, &st)
	if err != nil {
		http.Error(w, "Error in marshaling body", http.StatusInternalServerError)
		return
	}

	er := s.service.Post(ctx, &st)
	if er != nil {
		http.Error(w, er.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)

	if _, er = w.Write([]byte("Successfully Inserted")); er != nil {
		log.Println(er)
	}
}

func (s student) UpdateStudent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	resp, err := io.ReadAll(r.Body)

	ctx := r.Context()

	if err != nil {
		log.Println(err)

		return
	}

	var st models.Student

	er := json.Unmarshal(resp, &st)

	if er != nil {
		log.Println(er)

		if _, er = w.Write([]byte("Error : " + er.Error())); er != nil {
			log.Println(er)
		}

		return
	}

	er = s.service.Update(ctx, &st)

	if er != nil {
		log.Println(er)
		w.WriteHeader(http.StatusBadRequest)

		if _, er = w.Write([]byte("Error : " + er.Error())); er != nil {
			log.Println(er)
		}

		return
	}

	w.WriteHeader(http.StatusOK)

	if _, er = w.Write([]byte("Successfully Updated")); er != nil {
		log.Println(er)
	}
}

func (s student) DeleteStudent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := r.URL.Query().Get("id")

	ctx := r.Context()

	er := s.service.Delete(ctx, id)

	if er != nil {
		http.Error(w, er.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)

	_, er = w.Write([]byte("Delete Successful"))

	if er != nil {
		log.Println(er)
	}
}
