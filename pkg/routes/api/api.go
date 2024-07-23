package api

import (
	"fmt"
	"log"
	"net/http"
	"runtime/pprof"

	"github.com/debdut/Resoorch/lib/router"
)

// handler for /api route
func Handle(w http.ResponseWriter, r *http.Request) {
	_, r.URL.Path = router.ShiftPath(r.URL.Path)

	routes := map[string]http.HandlerFunc{
		"ping":   HandlePing,
		"memory": HandleMemory,
	}

	handler := router.Router(routes)
	handler(w, r)
}

// handler for /api/ping
func HandlePing(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "ok")
}

// handler for /api/memory
func HandleMemory(w http.ResponseWriter, r *http.Request) {
	// Gather memory allocations profile.
	profile := pprof.Lookup("allocs")

	// Write profile (human readable, via debug: 1) to HTTP response.
	err := profile.WriteTo(w, 1)
	if err != nil {
		log.Printf("Error: Failed to write allocs profile: %v", err)
	}
}
