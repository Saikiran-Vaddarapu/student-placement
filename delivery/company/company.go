package company

import (
	"StudentPlacement/models"
	"StudentPlacement/service"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type Company struct {
	serve service.Company
}

func New(companyservice service.Company) Company {
	return Company{serve: companyservice}
}

func (s Company) GetCompany(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := r.URL.Query().Get("id")

	ctx := r.Context()

	res, er := s.serve.Get(ctx, id)

	if er != nil {
		http.Error(w, er.Error(), http.StatusNotFound)
		return
	}

	c, _ := json.Marshal(res)

	if _, er = w.Write(c); er != nil {
		log.Println(er)
	}

	w.WriteHeader(http.StatusOK)
}

func (s Company) PostCompany(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	resp, err := io.ReadAll(r.Body)

	ctx := r.Context()

	if err != nil {
		http.Error(w, "Error in Reading Request", http.StatusBadRequest)
		return
	}

	var st models.Company

	if err = json.Unmarshal(resp, &st); err != nil {
		log.Println(err)
	}

	er := s.serve.Post(ctx, st)

	if er != nil {
		http.Error(w, er.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)

	if _, er = w.Write([]byte("Successfully Inserted")); er != nil {
		log.Println(er)
	}
}

func (s Company) UpdateCompany(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	resp, err := io.ReadAll(r.Body)

	ctx := r.Context()

	if err != nil {
		log.Println(err)
		return
	}

	var st models.Company

	if err = json.Unmarshal(resp, &st); err != nil {
		log.Println(err)
	}

	er := s.serve.Update(ctx, st)

	if er != nil {
		http.Error(w, er.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)

	if _, er = w.Write([]byte("Successfully Updated")); er != nil {
		log.Println(er)
	}
}

func (s Company) DeleteCompany(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := r.URL.Query().Get("id")

	ctx := r.Context()

	er := s.serve.Delete(ctx, id)

	if er != nil {
		http.Error(w, er.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)

	if _, er = w.Write([]byte("Delete Successful")); er != nil {
		log.Println(er)
	}
}
