package handler

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (h *Handler) GetPods(w http.ResponseWriter, r *http.Request) {
	fmt.Println("List all Pods")
	namespace := chi.URLParam(r, "namespace")
	if namespace == "" {
		namespace = "default"
	}
	podList, err := h.Client.CoreV1().Pods(namespace).List(r.Context(), metav1.ListOptions{})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	type podItem struct {
		Name      string `json:"name"`
		Namespace string `json:"namespace"`
		Phase     string `json:"phase"`
		PodIP     string `json:"podIP"`
		NodeName  string `json:"nodeName"`
	}

	result := make([]podItem, 0, len(podList.Items))
	for _, pod := range podList.Items {
		result = append(result, podItem{
			Name:      pod.Name,
			Namespace: pod.Namespace,
			Phase:     string(pod.Status.Phase),
			PodIP:     pod.Status.PodIP,
			NodeName:  pod.Spec.NodeName,
		})
	}

	writeJSON(w, http.StatusOK, result)
}
