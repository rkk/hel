package main

import (
	"fmt"
	"net/http"
	"log"
)

// Service provides different kinds of network durability impactings services.
type Service struct {
	serviceType string
	replyWait int
}

func (s *Service) badRequestHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		fmt.Printf("Getting HTTP GET request\n")
		return
	case "POST":
		fmt.Printf("Getting HTTP POST request\n")
	    return
	default:
		fmt.Printf("Getting unknown type request\n")
	}
	http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
}

func serveBadRequest() {
	s := &Service{
		serviceType: "badrequest",
		replyWait: 0,
	}

	fmt.Printf("Starting %s service...\n", s.serviceType)
	http.HandleFunc("/", s.badRequestHandler)
	if err := http.ListenAndServe(":9000", nil); err != nil {
		log.Fatal("ERROR: Cannot start server: %s\n", err)
	}

	fmt.Printf("Stopping service...\n")
}

func main() {
	serveBadRequest()
}
