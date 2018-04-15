package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

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

func main() {
	config, err := buildRuntimeConfiguration()
	if err != nil {
		log.Fatal("ERROR: Invalid runtime configuration")
		os.Exit(1)
	}

	shouldExit, code := parseRuntimeConfiguration(config)
	if shouldExit == true {
		os.Exit(code)
	}

	// Once implememented, the proxy service pattern replaces getInput().
	u, _ := url.Parse(config.input)
	input, err := getInput(*u)
	if err != nil {
		log.Fatal("ERROR: Cannot retrieve contents of %s", u.String())
	}

	s := buildService(config, input)

	serviceHandler := buildHandler(s.input, http.HandlerFunc(writeResponse))
	mux := new(http.ServeMux)
	mux.Handle(s.endpoint, serviceHandler)
	http.ListenAndServe(fmt.Sprintf(":%d", s.port), mux)
}
