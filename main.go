package main

import (
	"StudentPlacement/configs"
	companydatastore "StudentPlacement/datastore/company"
	studentdatastore "StudentPlacement/datastore/student"
	companydelivery "StudentPlacement/delivery/company"
	studentdelivery "StudentPlacement/delivery/student"
	"StudentPlacement/models"
	companyservice "StudentPlacement/service/company"
	studentservice "StudentPlacement/service/student"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

// nolint
func main() {
	dataConfig, er := configs.LoadConfig()

	if er != nil {
		log.Println(er)
		return
	}

	db, er := sql.Open(dataConfig.Driver, dataConfig.DataSource)

	if er != nil {
		log.Fatal(er)
		return
	}

	defer db.Close()

	database := studentdatastore.New(db)
	service := studentservice.New(database)
	student := studentdelivery.New(service)

	http.HandleFunc("/student", student.GetStudentByID)

	router := mux.NewRouter()

	router.HandleFunc("/allstudents", authorize(student.GetAllStudents)).Methods("GET")
	router.HandleFunc("/student", authorize(student.GetStudentByID)).Methods("GET")
	router.HandleFunc("/student", authorize(student.PostStudent)).Methods("POST")
	router.HandleFunc("/student", authorize(student.DeleteStudent)).Methods("DELETE")
	router.HandleFunc("/student", authorize(student.UpdateStudent)).Methods("PUT")

	companyDatabase := companydatastore.New(db)
	companyService := companyservice.New(companyDatabase)
	company := companydelivery.New(companyService)

	router.HandleFunc("/company", authorize(company.GetCompany)).Methods("GET")
	router.HandleFunc("/company", authorize(company.PostCompany)).Methods("POST")
	router.HandleFunc("/company", authorize(company.DeleteCompany)).Methods("DELETE")
	router.HandleFunc("/company", authorize(company.UpdateCompany)).Methods("PUT")

	server := &http.Server{
		Addr:              ":8081",
		ReadHeaderTimeout: models.ReadoutTimer * time.Second,
		Handler:           router,
	}

	if er := server.ListenAndServe(); er != nil {
		log.Println(er)
	}
}

// nolint
func authorize(orgHandler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if er := validateAPIKey(r.Header.Get(models.APIKey)); er != nil {
			w.WriteHeader(http.StatusBadRequest)

			if _, er := w.Write([]byte("Invalid Header")); er != nil {
				log.Println(er)
			}

			return
		}

		if r.Header.Get("Content-Type") != models.ContentType {
			w.WriteHeader(http.StatusUnsupportedMediaType)

			if _, er := w.Write([]byte("Header Content-Type is incorrect")); er != nil {
				log.Println(er)
			}

			return
		}

		orgHandler.ServeHTTP(w, r)
	}
}

func validateAPIKey(authKey string) error {
	dataConfig, er := configs.LoadConfig()

	if er != nil {
		log.Println(er)
		return er
	}

	if authKey != dataConfig.Key1 && authKey != dataConfig.Key2 {
		return errors.New("failed to authenticate request")
	}

	return nil
}
