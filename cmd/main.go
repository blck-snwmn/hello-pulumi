package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/pulumi/pulumi/sdk/v3/go/auto"
)

func main() {
	ctx := context.Background()
	w, err := auto.NewLocalWorkspace(ctx)
	if err != nil {
		fmt.Printf("Failed to setup and run http server: %v\n", err)
		os.Exit(1)
	}
	err = w.InstallPlugin(ctx, "aws", "v5.39.0")
	if err != nil {
		fmt.Printf("Failed to install program plugins: %v\n", err)
		os.Exit(1)
	}

	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Request: %s %s\n", r.Method, r.URL.Path)
	switch r.Method {
	case http.MethodGet:
		getHandler(w, r)
	case http.MethodPost:
		postHandler(w, r)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ws, err := auto.NewLocalWorkspace(ctx, auto.WorkDir("."))
	if err != nil {
		log.Printf("Failed to setup and run http server: %v\n", err)
		http.Error(w, "Failed to setup and run http server", http.StatusInternalServerError)
		return
	}
	stacks, err := ws.ListStacks(ctx)
	if err != nil {
		log.Printf("Failed to list stacks: %v\n", err)
		http.Error(w, "Failed to list stacks", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(stacks)
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, World!")
}
