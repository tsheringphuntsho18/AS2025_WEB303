package main

import (
    "fmt"
    "log"
    "net/http"
    "net/http/httputil"
    "net/url"
    "strings"

    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
    consulapi "github.com/hashicorp/consul/api"
)

func discoverService(serviceName string) (string, error) {
    config := consulapi.DefaultConfig()
    config.Address = "consul:8500"

    consul, err := consulapi.NewClient(config)
    if err != nil {
        return "", err
    }

    services, _, err := consul.Health().Service(serviceName, "", true, nil)
    if err != nil {
        return "", err
    }

    if len(services) == 0 {
        return "", fmt.Errorf("no healthy instances of %s", serviceName)
    }

    service := services[0].Service
    return fmt.Sprintf("http://%s:%d", service.Address, service.Port), nil
}

func proxyToService(serviceName, stripPrefix string) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Discover service dynamically
        targetURL, err := discoverService(serviceName)
        if err != nil {
            http.Error(w, err.Error(), http.StatusServiceUnavailable)
            return
        }

        target, _ := url.Parse(targetURL)
        proxy := httputil.NewSingleHostReverseProxy(target)

        r.URL.Path = strings.TrimPrefix(r.URL.Path, "/api")
        log.Printf("Proxying to %s at %s", serviceName, targetURL)
        proxy.ServeHTTP(w, r)
    }
}

func main() {
    r := chi.NewRouter()
    r.Use(middleware.Logger)

    r.HandleFunc("/api/users*", proxyToService("user-service", "/users"))
    r.HandleFunc("/api/menu*", proxyToService("menu-service", "/menu"))
    r.HandleFunc("/api/orders*", proxyToService("order-service", "/orders"))

    log.Println("API Gateway starting on :8080")
    http.ListenAndServe(":8080", r)
}