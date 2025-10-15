package geocoder

import (
	"AddressService/internal/domains/message/model"
	"bytes"
	"context"
	"fmt"
	"net/http"
	"time"

	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigFastest

type Geocoder struct {
	baseURL string
	client  *http.Client
	batch   int
}

func New(baseURL string, timeoutMs int, maxConns int) *Geocoder {
	return &Geocoder{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: time.Duration(timeoutMs) * time.Millisecond,
			Transport: &http.Transport{
				MaxIdleConns:        maxConns,
				MaxIdleConnsPerHost: maxConns,
				MaxConnsPerHost:     maxConns,
				IdleConnTimeout:     90 * time.Second,
				ForceAttemptHTTP2:   true,
			},
		},
		batch: 100,
	}
}

func (g *Geocoder) GetAddresses(ctx context.Context, positions []model.Pos) ([]string, error) {
	results := make([]string, 0, len(positions))

	for start := 0; start < len(positions); start += g.batch {
		end := start + g.batch
		if end > len(positions) {
			end = len(positions)
		}

		batch := positions[start:end]
		addrs, err := g.getBatch(ctx, batch)
		if err != nil {
			return nil, fmt.Errorf("batch %d-%d failed: %w", start, end, err)
		}
		results = append(results, addrs...)
	}

	return results, nil
}

func (g *Geocoder) getBatch(ctx context.Context, positions []model.Pos) ([]string, error) {
	body, err := json.Marshal(positions)
	if err != nil {
		return nil, fmt.Errorf("marshal batch: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, g.baseURL+"/reverse_batch", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("create req: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := g.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do req: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("geocache batch status %d", resp.StatusCode)
	}

	var addrs []string
	if err := json.NewDecoder(resp.Body).Decode(&addrs); err != nil {
		return nil, fmt.Errorf("decode resp: %w", err)
	}

	return addrs, nil
}
