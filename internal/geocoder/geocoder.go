package geocoder

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Client struct {
	BaseURL string
}

type nominatimResp struct {
	DisplayName string `json:"display_name"`
}

func (c *Client) ReverseGeocode(lat, lon float64) (string, error) {
	url := fmt.Sprintf("%s?format=jsonv2&lat=%f&lon=%f", c.BaseURL, lon, lat)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	log.Println()
	var res nominatimResp
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", err
	}
	log.Println(res.DisplayName)
	return res.DisplayName, nil
}
