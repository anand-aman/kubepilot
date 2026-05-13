package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
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

func (h *Handler) CreateNamespace(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Name string `json:"name"`
	}

	var req request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, fmt.Errorf("invalid request body: %w", err))
		return
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		writeError(w, http.StatusBadRequest, fmt.Errorf("namespace name cannot be empty"))
		return
	}

	_, err := h.Client.CoreV1().Namespaces().Create(r.Context(), &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}, metav1.CreateOptions{})

	if err != nil {
		if k8serrors.IsAlreadyExists(err) {
			writeError(w, http.StatusConflict, fmt.Errorf("namespace already exists"))
			return
		}
		writeError(w, http.StatusInternalServerError, fmt.Errorf("failed to create namespace: %w", err))
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{"message": "Namespace created successfully"})
}
