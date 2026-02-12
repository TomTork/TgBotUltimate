package helper

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

func Get(ctx context.Context, url string, headers map[string]string, out interface{}) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	client := &http.Client{Timeout: 120 * time.Second}
	response, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(response.Body)
	if response.StatusCode != http.StatusOK {
		_, err := io.ReadAll(response.Body)
		return err
	}
	if err := json.NewDecoder(response.Body).Decode(out); err != nil {
		return err
	}
	return nil
}
