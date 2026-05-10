package handler

import (
	"fmt"
	"net/http"
)

type Handler struct {
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	fmt.Println("List all pods")
	w.Write([]byte("Success"))
}
