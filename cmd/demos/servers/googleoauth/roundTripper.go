package main

import (
	"bytes"
	"fmt"
	"golang-server/cmd/product/makanplace/httplog"
	"io"
	"net/http"
	"time"
)

type LoggingRoundTripper struct {
	rt http.RoundTripper
}

func (lrt *LoggingRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	start := time.Now()

	req = httplog.HttpRequestWithValues(req)

	prefix := httplog.SPrintHttpRequestPrefix(req)
	fmt.Printf("%s received", prefix)
	if req.Body != nil {
		bodyBytes, _ := io.ReadAll(req.Body)
		req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes)) // Reset body
		fmt.Printf("[%s] body:", bodyBytes)
	}

	// Perform the actual request
	resp, err := lrt.rt.RoundTrip(req)
	if err != nil {
		fmt.Printf("%s Round Trip Error: [%s]", prefix, err)
		return nil, err
	}

	// Log response
	fmt.Printf("[HTTP Response] Status: %s in %v\n", resp.Status, time.Since(start))
	if resp.Body != nil {
		bodyBytes, _ := io.ReadAll(resp.Body)
		resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes)) // Reset body
		fmt.Printf("%s Response Body: [%s]", prefix, bodyBytes)
	}

	return resp, nil
}
