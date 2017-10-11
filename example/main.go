package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/beme/abide/example/models"
)

func firstHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"A": map[string][]string{
			"B": []string{"foo", "bar"},
		},
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
	post := &models.Post{
		Title: "Hello World",
		Body:  "Foo Bar",
	}
	data := map[string]interface{}{
		"post": post,
		"stats": map[string]interface{}{
			"updated_at": time.Now().Unix(),
			"things": []map[string]interface{}{
				{
					"updated_at": time.Now().Unix(),
				},
			},
		},
	}

	body, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

func thirdHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func fourthHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`Hello World.`))
}

func main() {
	http.HandleFunc("/first", firstHandler)
	http.HandleFunc("/second", secondHandler)
	http.HandleFunc("/third", thirdHandler)
	http.HandleFunc("/fourth", fourthHandler)
	http.ListenAndServe(":8080", nil)
}
