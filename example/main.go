package main

import (
	"encoding/json"
	"net/http"
)

func firstHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"key": "value",
	}

	body, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

func secondHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]int{
		"key": 1,
	}

	body, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

func main() {
	http.HandleFunc("/first", firstHandler)
	http.HandleFunc("/second", secondHandler)
	http.ListenAndServe(":8080", nil)
}
