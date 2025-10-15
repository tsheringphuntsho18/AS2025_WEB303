package main

import (
    "log"
    "net/http"
    "os"
    "menu-service/database"
    "menu-service/handlers"

    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"

    "fmt"
    consulapi "github.com/hashicorp/consul/api"
)

func registerWithConsul(serviceName string, port int) error {
    config := consulapi.DefaultConfig()
    config.Address = "consul:8500"

    consul, err := consulapi.NewClient(config)
    if err != nil {
        return err
    }

    hostname, _ := os.Hostname()

    registration := &consulapi.AgentServiceRegistration{
        ID:      fmt.Sprintf("%s-%s", serviceName, hostname),
        Name:    serviceName,
        Port:    port,
        Address: hostname,
        Check: &consulapi.AgentServiceCheck{
            HTTP:     fmt.Sprintf("http://%s:%d/health", hostname, port),
            Interval: "10s",
            Timeout:  "3s",
        },
    }

    return consul.Agent().ServiceRegister(registration)
}

func main() {
    // Connect to dedicated menu database
    dsn := os.Getenv("DATABASE_URL")
    if dsn == "" {
        dsn = "host=localhost user=postgres password=postgres dbname=menu_db port=5432 sslmode=disable"
    }

    if err := database.Connect(dsn); err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }

    r := chi.NewRouter()
    r.Use(middleware.Logger)

    // Add health endpoint
    r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
    })

    // Menu endpoints (note: no /api prefix)
    r.Get("/menu", handlers.GetMenu)
    r.Post("/menu", handlers.CreateMenuItem)
    r.Get("/menu/{id}", handlers.GetMenuItem)

    // Register service with Consul
    if err := registerWithConsul("menu-service", 8082); err != nil {
        log.Fatalf("Failed to register service with Consul: %v", err)
    }

    port := os.Getenv("PORT")
    if port == "" {
        port = "8082"
    }

    log.Printf("Menu service starting on :%s", port)
    http.ListenAndServe(":"+port, r)
}