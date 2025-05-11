package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"
)

type LoggingRoundTripper struct {
	rt http.RoundTripper
}

func (lrt *LoggingRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	start := time.Now()

	// Log request
	fmt.Printf("[HTTP Request] %s %s\n", req.Method, req.URL.String())
	if req.Body != nil {
		bodyBytes, _ := io.ReadAll(req.Body)
		req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes)) // Reset body
		fmt.Printf("[Request Body] %s\n", string(bodyBytes))
	}

	// Perform the actual request
	resp, err := lrt.rt.RoundTrip(req)
	if err != nil {
		fmt.Printf("[HTTP Error] %v\n", err)
		return nil, err
	}

	// Log response
	fmt.Printf("[HTTP Response] Status: %s in %v\n", resp.Status, time.Since(start))
	if resp.Body != nil {
		bodyBytes, _ := io.ReadAll(resp.Body)
		resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes)) // Reset body
		fmt.Printf("[Response Body] %s\n", string(bodyBytes))
	}

	return resp, nil
}
