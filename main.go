package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

// Service provides different kinds of network durability impacting services.
type Service struct {
	input    http.Response
	endpoint string
	port     int
	service  string
	delay    int
}

// RuntimeConfiguration provides a structured configuration for application invocation.
type RuntimeConfiguration struct {
	listMode  bool
	usageMode bool
	endpoint  string
	input     string
	port      int
	service   string
	delay     int
}

// Hack to avoid warning "should not use basic type as key in Context.WithValue()".
type key int

const ctxKey key = 1

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
	if err := http.ListenAndServe(fmt.Sprintf(":%d", s.port), nil); err != nil {
		log.Fatal("ERROR: Cannot start server: %s\n", err)
	}

	fmt.Printf("Stopping service...\n")
}

func getInput(u url.URL) (http.Response, error) {
	var client = &http.Client{Timeout: 3 * time.Second}
	response, err := client.Get(u.String())
	if err != nil {
		return http.Response{}, errors.New("Setting up HTTP client failed")
	}
	defer response.Body.Close()

	return *response, nil
}

// Build a HTTP handler from input.
func buildHandler(input http.Response, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), ctxKey, input)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Sends an HTTP response based on the input contents.
func writeResponse(w http.ResponseWriter, r *http.Request) {
	input := r.Context().Value(ctxKey)
	if input == nil {
		w.WriteHeader(500)
		w.Write([]byte("No message"))
	} else {
		w.WriteHeader(200)
		w.Write([]byte(input.(string)))
	}
}

func displayUsage() {
	fmt.Printf("Usage: hel OPTIONS\n")
	fmt.Printf("  --list             - Lists the provided services\n")
	fmt.Printf("  --input=URL        - Retrieve input data from URL\n")
	fmt.Printf("  --endpoint=PATH    - Serves Hel from the PATH endpoint\n")
	fmt.Printf("  --port=PORT        - Serves Hel from the PORT TCP port number\n")
	fmt.Printf("  --service=NAME     - Use the service type NAME\n")
	fmt.Printf("  --delay=SECONDS    - Delay the response for SECONDS\n")
	fmt.Printf("\n")
}

func displayServices() {
	fmt.Printf("Services:\n")
	fmt.Printf("  badrequest  - Returns HTTP 400 Bad Request\n")
	fmt.Printf("\n")
}

func buildRuntimeConfiguration() (RuntimeConfiguration, error) {
	listPtr := flag.Bool("list", false, "Lists the provided services")
	endpointPtr := flag.String("endpoint", "", "Endpoint to serve")
	inputPtr := flag.String("input", "", "Retrieve data from this URL")
	portPtr := flag.Int("port", 0, "TCP port to serve from")
	servicePtr := flag.String("service", "", "Service type, see --list")
	delayPtr := flag.Int("delay", 0, "Delay the response")
	flag.Parse()

	var usageMode = false
	if (*endpointPtr == "" || *inputPtr == "" || *portPtr == 0 || *servicePtr == "") && *listPtr == false {
		usageMode = true
	}

	runtime := RuntimeConfiguration{
		listMode:  *listPtr,
		usageMode: usageMode,
		endpoint:  *endpointPtr,
		input:     *inputPtr,
		port:      *portPtr,
		service:   *servicePtr,
		delay:     *delayPtr,
	}
	return runtime, nil
}

func main() {

	config, err := buildRuntimeConfiguration()
	if err != nil {
		log.Fatal("ERROR: Invalid runtime configuration")
		os.Exit(1)
	}

	if config.listMode == true {
		displayServices()
		os.Exit(0)
	}

	if config.usageMode == true {
		displayUsage()
		os.Exit(0)
	}
	u, err := url.Parse(config.input)
	if err != nil {
		log.Fatal("ERROR: %s is not a valid input URL\n", config.input)
		os.Exit(1)
	}

	input, err := getInput(*u)
	if err != nil {
		log.Fatal("ERROR: Cannot retrieve contents of %s", u.String())
	}

	s := Service{
		endpoint: config.endpoint,
		input:    input,
		port:     config.port,
		service:  config.service,
		delay:    config.delay,
	}

	serviceHandler := buildHandler(s.input, http.HandlerFunc(writeResponse))
	mux := new(http.ServeMux)
	mux.Handle(s.endpoint, serviceHandler)
	http.ListenAndServe(fmt.Sprintf(":%d", s.port), mux)
}
