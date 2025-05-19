package main

import (
	"log"
	"net/http"
	"net/url"
)

func main() {
	c := http.Client{}
	res, rErr := c.Do(&http.Request{
		Method: "GET",
		URL:    &url.URL{Scheme: "https", Host: "www.googleapis.com", Path: "/"},
	})

	log.Printf("[trying to get] want %d got %s err=%#v\n", 404, res.Status, rErr)
}
