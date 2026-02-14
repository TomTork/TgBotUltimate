package helper

import (
	"TgBotUltimate/errors"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
)

var client = &http.Client{Timeout: 120 * time.Second}

func Get(ctx context.Context, url string, headers map[string]string, out interface{}) error {
	if out == nil {
		log.Println(errors.OutInterfaceIsNil)
		return http.ErrNotSupported
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
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

func Post(ctx context.Context, url string, headers map[string]string, body interface{}, out interface{}) error {
	if out == nil {
		log.Println(errors.OutInterfaceIsNil)
		return http.ErrNotSupported
	}
	__body, err := json.Marshal(body)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(__body))
	if err != nil {
		return err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	req.Header.Set("Content-Type", "application/json")
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
