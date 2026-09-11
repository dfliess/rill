package odoo

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// requestTimeout bounds one JSON-2 call. A page of search_read is small; a call slower than
// this is a server in trouble, not a slow page.
const requestTimeout = 2 * time.Minute

// configProperties is the connector configuration (connectors/<name>.yaml).
type configProperties struct {
	URL      string `mapstructure:"url"`
	Database string `mapstructure:"database"`
	APIKey   string `mapstructure:"api_key"`
	Lang     string `mapstructure:"lang"`
}

// client calls Odoo's JSON-2 API: POST {url}/json/2/{model}/{method} with the method's
// arguments as a JSON object, authenticated with an API key sent as a bearer token.
type client struct {
	baseURL  string
	database string
	apiKey   string
	http     *http.Client
}

func newClient(conf *configProperties) (*client, error) {
	if conf.URL == "" {
		return nil, errors.New(`odoo: property "url" is required`)
	}
	if conf.APIKey == "" {
		return nil, errors.New(`odoo: property "api_key" is required`)
	}
	u, err := url.Parse(conf.URL)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
		return nil, fmt.Errorf("odoo: invalid url %q: it must look like https://host", conf.URL)
	}
	return &client{
		baseURL:  strings.TrimRight(conf.URL, "/"),
		database: conf.Database,
		apiKey:   conf.APIKey,
		http:     &http.Client{Timeout: requestTimeout},
	}, nil
}

// odooError is the body Odoo sends with a non-200 status. It also carries "debug", a server
// traceback, which is deliberately left out of our errors.
type odooError struct {
	Name    string `json:"name"`
	Message string `json:"message"`
}

// call runs model.method with args and decodes the result into out. Numbers are decoded as
// json.Number so integer ids survive intact.
func (c *client) call(ctx context.Context, model, method string, args map[string]any, out any) error {
	body, err := json.Marshal(args)
	if err != nil {
		return fmt.Errorf("odoo: %s.%s: encoding arguments: %w", model, method, err)
	}
	endpoint := c.baseURL + "/json/2/" + url.PathEscape(model) + "/" + url.PathEscape(method)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("odoo: %s.%s: %w", model, method, err)
	}
	req.Header.Set("Authorization", "bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	if c.database != "" {
		req.Header.Set("X-Odoo-Database", c.database)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("odoo: %s.%s: %w", model, method, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		raw, _ := io.ReadAll(io.LimitReader(resp.Body, 64<<10))
		var oe odooError
		if json.Unmarshal(raw, &oe) == nil && oe.Message != "" {
			return fmt.Errorf("odoo: %s.%s: HTTP %d: %s: %s", model, method, resp.StatusCode, oe.Name, oe.Message)
		}
		return fmt.Errorf("odoo: %s.%s: HTTP %d", model, method, resp.StatusCode)
	}

	dec := json.NewDecoder(resp.Body)
	dec.UseNumber()
	if err := dec.Decode(out); err != nil {
		return fmt.Errorf("odoo: %s.%s: decoding response: %w", model, method, err)
	}
	return nil
}
