package main

import "net/http"

// Service provides different kinds of network durability impacting services.
type Service struct {
	input    http.Response
	endpoint string
	port     int
	service  string
	delay    int
}

func buildService(config RuntimeConfiguration, input http.Response) Service {
	s := Service{
		endpoint: config.endpoint,
		input:    input,
		port:     config.port,
		service:  config.service,
		delay:    config.delay,
	}

	return s
}
