package clickhouse

import (
	"auth_service/helper"
	"bytes"
	"fmt"
	"net/http"
	"time"
)

var httpClient = &http.Client{Timeout: 3 * time.Second}

func Exec(query string) error {
	url := helper.ENV("CLICKHOUSE_URL") // http://localhost:8123

	req, err := http.NewRequest("POST", url, bytes.NewBufferString(query))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "text/plain")

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("clickhouse error: %d", resp.StatusCode)
	}
	return nil
}
