package nsq

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// PublishHTTP publishes a message to nsqd via HTTP with optional defer (delay in ms)
func PublishHTTP(nsqdHTTPAddr, topic string, message interface{}, deferMs int) error {
	b, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	if topic == "" {
		return fmt.Errorf("NSQ topic is empty")
	}

	url := fmt.Sprintf("%s/pub?topic=%s", nsqdHTTPAddr, topic)
	if deferMs > 0 {
		url = fmt.Sprintf("%s&defer=%d", url, deferMs)
	}
	resp, err := http.Post(url, "application/octet-stream", bytes.NewReader(b))
	if err != nil {
		return fmt.Errorf("failed to publish message via HTTP: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("nsqd returned non-200 status: %s", resp.Status)
	}
	return nil
}
