package main

import (
	"crypto/rand"
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
	w.WriteHeader(http.StatusOK)

	res, err := json.Marshal(&historyResponse{
		Records: records,
	})

	if err != nil {
		fmt.Printf("Error marshaling records: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(res)
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	type registerPayload struct {
		Namespace string `json:"namespace"`
	}

	type responsePayload struct {
		Success        bool            `json:"success"`
		Error          string          `json:"error,omitempty"`
		TokenNamespace *TokenNamespace `json:"data,omitempty"`
	}

	decoder := json.NewDecoder(r.Body)
	payload := &registerPayload{}
	decoder.Decode(&payload)

	namespace := payload.Namespace

	if namespace == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Add("Content-Type", "application/json")

	// Check if namespace already exists first
	ns, err := lookupNamespace(namespace)
	if err != nil {
		fmt.Printf("Error looking up namespace: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if ns != nil {
		res, err := json.Marshal(&responsePayload{
			Success: false,
			Error:   "Namespace already registered",
		})

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusBadRequest)
		w.Write(res)
		return
	}

	// Create namespace
	token := make([]byte, 12)
	rand.Read(token)
	ns, err = createNamespace(fmt.Sprintf("%x", token), namespace)
	if err != nil {
		fmt.Printf("Error creating namespace: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(&responsePayload{
		Success:        true,
		TokenNamespace: ns,
	})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
