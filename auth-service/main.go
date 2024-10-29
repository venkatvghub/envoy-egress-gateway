package main

import (
    "log"
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
)

// LogRequest logs detailed request information
func LogRequest(c *gin.Context) {
    log.Printf("\n=== Incoming Request at %s ===", time.Now().Format(time.RFC3339))
    log.Printf("Client IP: %s", c.ClientIP())
    log.Printf("Method: %s", c.Request.Method)
    log.Printf("Original URL: %s", c.Request.URL.String())
    log.Printf("Path: %s", c.Request.URL.Path)
    log.Printf("Raw Query: %s", c.Request.URL.RawQuery)
    log.Printf("Host: %s", c.Request.Host)
    log.Printf("Headers:")
    for name, values := range c.Request.Header {
        for _, value := range values {
            log.Printf("  %s: %s", name, value)
        }
    }
    log.Printf("==================\n")
}

// AuthHandler handles all authorization requests
func AuthHandler(c *gin.Context) {
    LogRequest(c)

    // Always allow the request by setting appropriate headers
    c.Header("x-auth-header", "allowed")
    
    // Add useful debugging headers
    c.Header("x-auth-original-method", c.Request.Method)
    c.Header("x-auth-original-path", c.Request.URL.Path)
    c.Header("x-auth-original-host", c.Request.Host)
    c.Header("x-auth-original-uri", c.Request.URL.String())

    // Log the authorization decision
    log.Printf("Authorizing request to: %s %s", c.Request.Method, c.Request.URL.String())
    
    // Return OK status for authorization
    c.JSON(http.StatusOK, gin.H{
        "status": "authorized",
        "path": c.Request.URL.Path,
        "method": c.Request.Method,
        "host": c.Request.Host,
    })
}

func main() {
    // Set Gin to release mode
    gin.SetMode(gin.ReleaseMode)
    
    // Create a new Gin router with default middleware
    r := gin.New()
    
    // Add recovery middleware
    r.Use(gin.Recovery())
    
    // Add logger middleware
    r.Use(gin.LoggerWithConfig(gin.LoggerConfig{
        SkipPaths: []string{"/health"},
    }))

    // Health check endpoint
    r.GET("/health", func(c *gin.Context) {
        log.Printf("Health check received")
        c.Status(http.StatusOK)
    })

    // Main authorization check endpoint - handle all paths
    r.Any("/check", AuthHandler)  // Handle the /check endpoint specifically
    r.NoRoute(AuthHandler)        // Handle all other paths

    // Start server
    log.Printf("Starting auth service on :8081")
    log.Printf("Ready to handle authorization requests")
    
    if err := r.Run(":8081"); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}
