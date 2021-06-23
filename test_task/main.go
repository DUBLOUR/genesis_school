package main

import (
	// "fmt"
	"encoding/json"
	"log"
	"net/http"
	"os"
	//    "reflect"
)

func Respond(w http.ResponseWriter, httpStatus int, data map[string]interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	json.NewEncoder(w).Encode(data)
}

func main() {

	f, err := os.OpenFile("/home/admin/go/genesis_school/test_task/log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	log.SetOutput(f)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		mess := map[string]interface{}{
			"status": "Fail",
		}
		Respond(w, http.StatusNotFound, mess)
	})

	http.HandleFunc("/user/create", func(w http.ResponseWriter, r *http.Request) {
		email := r.URL.Query().Get("email")
		password := r.URL.Query().Get("password")

		err := UserRegister(email, password)
		log.Println(r.URL.String())
		status := "Ok"
		if err != nil {
			status = err.Error()
		}

		mess := map[string]interface{}{
			"status": status,
		}
		Respond(w, http.StatusOK, mess)
	})

	http.HandleFunc("/user/login", func(w http.ResponseWriter, r *http.Request) {
		email := r.URL.Query().Get("email")
		password := r.URL.Query().Get("password")

		token, err := UserLogin(email, password)
		status := "Ok"
		if err != nil {
			status = err.Error()
		}

		mess := map[string]interface{}{
			"status": status,
			"token":  token,
		}
		Respond(w, http.StatusOK, mess)
	})

	http.HandleFunc("/btcRate", func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("token")

		log.Println(r.URL.String())
		log.Println(token)

		if IsAvaiableToken(token) {
			cost, err := Cost("BTCUAH")

			if err == nil {
				mess := map[string]interface{}{
					"status": "Ok",
					"BTCUAH": cost,
				}
				Respond(w, http.StatusOK, mess)
			} else {
				mess := map[string]interface{}{
					"status": "Internal error",
				}
				Respond(w, http.StatusInternalServerError, mess)
			}

			return
		}

		if token == "" {
			mess := map[string]interface{}{
				"status": "Missing token",
			}
			Respond(w, http.StatusBadRequest, mess)
			return
		}

		mess := map[string]interface{}{
			"status": "Invalid token",
		}
		Respond(w, http.StatusForbidden, mess)

	})

	http.ListenAndServe(":9990", nil)
}
