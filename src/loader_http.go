// HTTPRuleProvider implements RuleProvider that fetches rules from a remote URL via HTTP.

package gorulesengine

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// HTTPRuleProvider implements RuleProvider for fetching rules from a URL.
type HTTPRuleProvider struct {
	URL        string
	Client     *http.Client
	LastUpdate time.Time
	ETag       string
}

// NewHTTPRuleProvider creates a new HTTP rule provider.
func NewHTTPRuleProvider(url string) *HTTPRuleProvider {
	return &HTTPRuleProvider{
		URL: url,
		Client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// FetchRules fetches rules from the configured URL.
// It expects a JSON array of rules.
func (p *HTTPRuleProvider) FetchRules(ctx context.Context) ([]*Rule, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", p.URL, nil)
	if err != nil {
		return nil, &RuleEngineError{
			Type: ErrLoader,
			Msg:  fmt.Sprintf("failed to create request: %v", err),
			Err:  err,
		}
	}

	if p.ETag != "" {
		req.Header.Set("If-None-Match", p.ETag)
	}

	resp, err := p.Client.Do(req)
	if err != nil {
		return nil, &RuleEngineError{
			Type: ErrLoader,
			Msg:  fmt.Sprintf("failed to fetch rules: %v", err),
			Err:  err,
		}
	}

	if resp.StatusCode == http.StatusNotModified {
		return nil, nil // No changes
	}

	if resp.StatusCode != http.StatusOK {
		return nil, &RuleEngineError{
			Type: ErrLoader,
			Msg:  fmt.Sprintf("unexpected status code %d", resp.StatusCode),
		}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &RuleEngineError{
			Type: ErrLoader,
			Msg:  "failed to read response body",
			Err:  err,
		}
	}

	var rules []*Rule
	if err := json.Unmarshal(body, &rules); err != nil {
		return nil, &RuleEngineError{
			Type: ErrLoader,
			Msg:  "failed to unmarshal rules",
			Err:  err,
		}
	}

	p.ETag = resp.Header.Get("ETag")
	p.LastUpdate = time.Now()

	return rules, nil
}
