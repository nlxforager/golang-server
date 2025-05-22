package onemap

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"
)

type LatLong struct {
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
}

type Client struct {
	postalCodeToLatLong map[string]*LatLong
	l                   sync.RWMutex
}

func NewClient() *Client {
	return &Client{
		postalCodeToLatLong: make(map[string]*LatLong),
	}
}

type SearchResultJSON struct {
	Found         int `json:"found"`
	TotalNumPages int `json:"totalNumPages"`
	PageNum       int `json:"pageNum"`
	Results       []struct {
		SEARCHVAL string `json:"SEARCHVAL"`
		BLKNO     string `json:"BLK_NO"`
		ROADNAME  string `json:"ROAD_NAME"`
		BUILDING  string `json:"BUILDING"`
		ADDRESS   string `json:"ADDRESS"`
		POSTAL    string `json:"POSTAL"`
		X         string `json:"X"`
		Y         string `json:"Y"`
		LATITUDE  string `json:"LATITUDE"`
		LONGITUDE string `json:"LONGITUDE"`
	} `json:"results"`
}

func (c *Client) GetLatLong(postalCode string) (*LatLong, error) {
	if len(postalCode) != 6 {
		return nil, fmt.Errorf("postal code=%s is invalid", postalCode)
	}
	url := fmt.Sprintf("https://www.onemap.gov.sg/api/common/elastic/search?searchVal=%s&returnGeom=Y&getAddrDetails=N&pageNum=1", postalCode)

	c.l.RLock()
	latlong, ok := c.postalCodeToLatLong[postalCode]
	c.l.RUnlock()
	if ok {
		return latlong, nil
	}
	// Make the GET request
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var ll *LatLong
	defer c.Set(postalCode, ll)
	// Parse the JSON response
	var data SearchResultJSON
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	if len(data.Results) == 0 {
		return nil, errors.New("no results found")
	}

	ll = &LatLong{
		Latitude:  data.Results[0].LATITUDE,
		Longitude: data.Results[0].LONGITUDE,
	}

	return ll, nil
}
func (c *Client) Set(postalCode string, long *LatLong) {
	c.l.Lock()
	defer c.l.Unlock()
	c.postalCodeToLatLong[postalCode] = long
}
