package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func JSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(Response{Success: status < 400, Data: data})
}

func Error(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(Response{Success: false, Error: msg})
}

func ParseID(r *http.Request) (int, error) {
	return strconv.Atoi(r.PathValue("id"))
}

func Decode(r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}
