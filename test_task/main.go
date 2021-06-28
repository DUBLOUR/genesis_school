package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

const logFile string = "/var/log/btc-requester/responces.log"
const dbFile string = "/etc/btc-requester/users.csv"

func Respond(w http.ResponseWriter, r *http.Request, httpStatus int, data map[string]interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	json.NewEncoder(w).Encode(data)

	log.Println("REQ:", r.URL.String(), "\nSTATUS:", httpStatus, "BODY:", data)
}

func main() {

	loger, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer loger.Close()
	log.SetOutput(loger)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		mess := map[string]interface{}{
			"status": "Fail",
		}
		Respond(w, r, http.StatusNotFound, mess)
	})

	http.HandleFunc("/user/create", func(w http.ResponseWriter, r *http.Request) {
		email := r.URL.Query().Get("email")
		password := r.URL.Query().Get("password")

		httpStatus, err := UserRegister(email, password)

		status := "Ok"
		if err != nil {
			status = err.Error()
		}

		mess := map[string]interface{}{
			"status": status,
		}
		Respond(w, r, httpStatus, mess)
	})

	http.HandleFunc("/user/login", func(w http.ResponseWriter, r *http.Request) {
		email := r.URL.Query().Get("email")
		password := r.URL.Query().Get("password")

		token, httpStatus, err := UserLogin(email, password)

		mess := map[string]interface{}{}
		if err == nil {
			mess = map[string]interface{}{
				"status": "Ok",
				"token":  token,
			}
		} else {
			mess = map[string]interface{}{
				"status": err.Error(),
			}
		}
		Respond(w, r, httpStatus, mess)
	})

	http.HandleFunc("/btcRate", func(w http.ResponseWriter, r *http.Request) {
		var urlToken, headerToken, token string
		urlToken = r.URL.Query().Get("token")
		headerToken = r.Header.Get("X-API-Key")

		if headerToken != "" {
			token = headerToken
		} else {
			token = urlToken
		}

		log.Println(r.URL.String())
		log.Println(token)

		if IsAvaiableToken(token) {
			cost, httpStatus, err := Cost("BTCUAH")

			mess := map[string]interface{}{}
			if err == nil {
				mess = map[string]interface{}{
					"status": "Ok",
					"BTCUAH": cost,
				}
			} else {
				mess = map[string]interface{}{
					"status": err.Error(),
				}
			}
			Respond(w, r, httpStatus, mess)
			return
		}

		if token == "" {
			mess := map[string]interface{}{
				"status": "Missing token",
			}
			Respond(w, r, http.StatusBadRequest, mess)
			return
		}

		mess := map[string]interface{}{
			"status": "Invalid token",
		}
		Respond(w, r, http.StatusForbidden, mess)

	})

	http.ListenAndServe(":9990", nil)
}
