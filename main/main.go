package main

import (
    "net/http"
    "MindAssistantBackend/projects"
    "MindAssistantBackend/data"
)

func main() {
    // Projects handlers
    http.HandleFunc("/projects", projects.List)
    http.HandleFunc("/projects/show", projects.Item)
    http.HandleFunc("/projects/create", projects.Create)
    http.HandleFunc("/projects/update", projects.Update)
    http.HandleFunc("/projects/delete", projects.Delete)
    // Data handlers
    http.HandleFunc("/data", data.List)
    http.HandleFunc("/data/project", data.ListByProject)
    http.HandleFunc("/data/show", data.Item)
    http.HandleFunc("/data/create", data.Create)
    http.HandleFunc("/data/update", data.Update)
    http.HandleFunc("/data/delete", data.Delete)
    // Start server
    http.ListenAndServe(":3000", nil)
}