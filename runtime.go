package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
)

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

func parseRuntimeConfiguration(r RuntimeConfiguration) (bool, int) {
	if r.listMode == true {
		displayServices()
		return true, 0
	}

	if r.usageMode == true {
		displayUsage()
		return true, 0
	}

	_, err := url.Parse(r.input)
	if err != nil {
		log.Fatal("ERROR: %s is not a valid input URL\n", r.input)
		return true, 1
	}

	return true, 0
}
