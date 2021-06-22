package main

import (
//    "fmt"
    "net/http"
    "encoding/json"
//    "reflect"
)


func Respond(w http.ResponseWriter, httpStatus int, data map[string] interface{})  {
	//w.WriteHeader(httpStatus)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	json.NewEncoder(w).Encode(data)
}

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		mess := map[string]interface{}{
			"status" : "false",
		}
		Respond(w, http.StatusNotFound, mess)
	})

	http.HandleFunc("/user/create", func(w http.ResponseWriter, r *http.Request) {
		mess := map[string]interface{}{
			"status" : "ok",
			"message" : "user create",
		}
		Respond(w, http.StatusOK, mess)
	})


	http.HandleFunc("/user/login", func(w http.ResponseWriter, r *http.Request) {
		mess := map[string]interface{}{
			"status" : "ok",
			"message" : "user login",
		}
		Respond(w, http.StatusOK, mess)
	})


	http.HandleFunc("/btcRate", func(w http.ResponseWriter, r *http.Request) {
		cost := Cost("BTCUAH")
		mess := map[string]interface{}{
			"status" : "fail",
			"BTCUAH" : cost,
		}
		Respond(w, http.StatusInternalServerError, mess)
	})

	http.ListenAndServe(":9990", nil)
}

