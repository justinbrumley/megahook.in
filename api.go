package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func historyHandler(w http.ResponseWriter, r *http.Request) {
	type historyResponse struct {
		Records []Record `json:"records"`
	}

	vars := mux.Vars(r)
	name := vars["name"]

	records, err := getRecords(name)
	if err != nil {
		fmt.Printf("Error fetching records: %v %v\n", name, err)
		w.WriteHeader(500)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(200)

	res, err := json.Marshal(&historyResponse{
		Records: records,
	})

	if err != nil {
		fmt.Printf("Error marshaling records: %v\n", err)
		w.WriteHeader(500)
		return
	}

	w.Write(res)
}
