package application

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"k8s.io/client-go/kubernetes"

	"github.com/anand-aman/kubepilot/k8s"
)

type App struct {
	router http.Handler
	config Config
	client *kubernetes.Clientset
}

func New(config Config) *App {
	app := &App{
		config: config,
	}

	return app
}

func (a *App) Start(ctx context.Context) error {
	// Initialize K8s client at startup
	client, err := k8s.NewClient()
	if err != nil {
		return fmt.Errorf("failed to create kubernetes client: %w", err)
	}
	a.client = client
	log.Println("✓ Kubernetes client initialized successfully")

	a.loadRoutes()

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", a.config.ServerPort),
		Handler: a.router,
	}

	var serverErr error
	ch := make(chan error, 1)
	go func() {
		serverErr = server.ListenAndServe()
		if serverErr != nil {
			ch <- fmt.Errorf("Failed to start server: %w", serverErr)
		}
		close(ch)
	}()

	select {
	case err = <-ch:
		return err
	case <-ctx.Done():
		log.Println("Shutting down the server...")
		timeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		err = server.Shutdown(timeout)
	}
	return err
}
