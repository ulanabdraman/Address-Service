package geocoder

import (
	"AddressService/internal/domains/message/model"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

var client = &http.Client{Timeout: 5 * time.Second}

func GetAddress(ctx context.Context, pos model.Pos) (model.Address, error) {
	url := fmt.Sprintf("http://labauto.kz:8080/nominatim/reverse.php?format=jsonv2&lat=%f&lon=%f", pos.X, pos.Y)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return model.Address{}, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return model.Address{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return model.Address{}, fmt.Errorf("geocoder returned status %d", resp.StatusCode)
	}

	var decoded struct {
		Address model.Address `json:"address"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&decoded); err != nil {
		return model.Address{}, err
	}
	log.Println("Geocoded address:", decoded.Address)

	return decoded.Address, nil
}
