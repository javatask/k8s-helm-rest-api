package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
	"helm.sh/helm/v3/pkg/cli"
)

// Server represents the REST API server
type Server struct {
	Router      *mux.Router
	HelmEnv     *cli.EnvSettings
	KubeConfig  string
	RegistryURL string
}

// NewServer initializes and returns a new Server instance
func NewServer() *Server {
	// Default to home directory/.kube/config if KUBECONFIG is not set
	kubeconfig := os.Getenv("KUBECONFIG")
	if kubeconfig == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			log.Printf("Failed to get user home directory: %v", err)
		} else {
			kubeconfig = filepath.Join(home, ".kube", "config")
		}
	}

	helmEnv := cli.New()

	s := &Server{
		Router:      mux.NewRouter(),
		HelmEnv:     helmEnv,
		KubeConfig:  kubeconfig,
		RegistryURL: os.Getenv("HELM_REGISTRY_URL"),
	}

	s.setupRoutes()
	return s
}

// setupRoutes configures all API routes
func (s *Server) setupRoutes() {
	api := s.Router.PathPrefix("/api/v1").Subrouter()

	// Health check endpoint
	api.HandleFunc("/health", s.healthHandler).Methods("GET")

	// Chart installation and management
	api.HandleFunc("/charts/install", s.installChartHandler).Methods("POST")
	api.HandleFunc("/charts/upgrade", s.upgradeChartHandler).Methods("PUT")
	api.HandleFunc("/charts/uninstall", s.uninstallChartHandler).Methods("DELETE")

	// Release information
	api.HandleFunc("/releases", s.listReleasesHandler).Methods("GET")
	api.HandleFunc("/releases/{name}", s.getReleaseHandler).Methods("GET")
	api.HandleFunc("/releases/{name}/history", s.getReleaseHistoryHandler).Methods("GET")
	api.HandleFunc("/releases/{name}/status", s.getReleaseStatusHandler).Methods("GET")

	// Add global middleware
	s.Router.Use(logMiddleware)
}

// healthHandler provides a simple health check endpoint
func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, map[string]string{"status": "healthy"})
}

// logMiddleware logs all incoming requests
func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
