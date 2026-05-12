package handler

import (
	"fmt"
	"net/http"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (h *Handler) GetNamespaces(w http.ResponseWriter, r *http.Request) {
	fmt.Println("List all namespaces")
	nsList, err := h.Client.CoreV1().Namespaces().List(r.Context(), metav1.ListOptions{})
	if err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Errorf("failed to list namespaces: %w", err))
		return
	}

	type nsItem struct {
		Name   string `json:"name"`
		Status string `json:"status"`
	}

	result := make([]nsItem, 0, len(nsList.Items))
	for _, ns := range nsList.Items {
		result = append(result, nsItem{
			Name:   ns.Name,
			Status: string(ns.Status.Phase),
		})
	}

	writeJSON(w, http.StatusOK, result)
}
