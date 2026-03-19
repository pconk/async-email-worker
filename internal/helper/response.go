package helper

import (
	"encoding/json"
	"net/http"
)

type PaginationMeta struct {
	CurrentPage int   `json:"current_page"`
	TotalItems  int64 `json:"total_items"`
	TotalPages  int   `json:"total_pages"`
	Limit       int   `json:"limit"`
}

type WebResponse struct {
	Code    int         `json:"code"`
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type WebResponsePaging struct {
	Code       int             `json:"code"`
	Status     string          `json:"status"`
	Message    string          `json:"message"`
	Data       interface{}     `json:"data,omitempty"`
	Pagination *PaginationMeta `json:"pagination,omitempty"`
}

func SendResponse(w http.ResponseWriter, code int, status, msg string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(WebResponse{
		Code:    code,
		Status:  status,
		Message: msg,
		Data:    data,
	})
}

func SendResponseWithPaging(w http.ResponseWriter, code int, status, msg string, data interface{}, pagination *PaginationMeta) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(WebResponsePaging{
		Code:       code,
		Status:     status,
		Message:    msg,
		Data:       data,
		Pagination: pagination, // Akan hilang di JSON jika nil
	})
}
