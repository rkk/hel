package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

// Service provides different kinds of network durability impacting services.
type Service struct {
	endpoint string
	port     int
	service  string
	delay    int
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
		fmt.Printf("Getting HTTP unknown type request\n")
	}

	if s.delay > 0 {
		time.Sleep(time.Duration(s.delay) * 1000)
	}
	http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
}

func serveBadRequest(s Service) {

	fmt.Printf("Starting %s service...\n", s.service)
	http.HandleFunc(s.endpoint, s.badRequestHandler)
	if err := http.ListenAndServe(":9000", nil); err != nil {
		log.Fatal("ERROR: Cannot start server: %s\n", err)
	}

	fmt.Printf("Stopping service...\n")
}

func displayUsage() {
	fmt.Printf("Usage: hel OPTIONS\n")
	fmt.Printf("  --list             Lists the provided services\n")
	fmt.Printf("  --endpoint=PATH    Serves Hel from the PATH endpoint\n")
	fmt.Printf("  --port=PORT        Serves Hel from the PORT TCP port number\n")
	fmt.Printf("  --service=NAME     Use the service type NAME\n")
	fmt.Printf("  --delay=SECONDS    Delay the response for SECONDS\n")
	fmt.Printf("\n")
}

func displayServices() {
	fmt.Printf("Services:\n")
	fmt.Printf("        badrequest Returns HTTP 400 Bad Request\n")
	fmt.Printf("\n")
}

func main() {
	listPtr := flag.Bool("list", false, "Lists the provided services")
	endpointPtr := flag.String("endpoint", "", "Endpoint to serve")
	portPtr := flag.Int("port", 0, "TCP port to serve from")
	servicePtr := flag.String("service", "", "Service type, see --list")
	delayPtr := flag.Int("delay", 0, "Delay the response")
	flag.Parse()

	if (*endpointPtr == "" || *portPtr == 0 || *servicePtr == "") && *listPtr == false {
		displayUsage()
		os.Exit(0)
	}

	if *listPtr == true {
		displayServices()
		os.Exit(0)
	}

	s := Service{
		endpoint: *endpointPtr,
		port:     *portPtr,
		service:  *servicePtr,
		delay:    *delayPtr,
	}
	serveBadRequest(s)
}
