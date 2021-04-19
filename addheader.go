package addheader

import (
	"context"
	"fmt"
	"net/http"
)

// Config the plugin configuration.
type Config struct {
	Headers map[string]string `json:"headers,omitempty"`
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{
		Headers: make(map[string]string),
	}
}

// Demo a Demo plugin.
type Demo struct {
	next     http.Handler
	headers  map[string]string
	name     string
}

// New created a new Demo plugin.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	if len(config.Headers) == 0 {
		return nil, fmt.Errorf("headers cannot be empty")
	}

	return &Demo{
		headers:  config.Headers,
		next:     next,
		name:     name,
	}, nil
}

func (d *Demo) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	for k, v := range d.headers {
		req.Header.Set(k, v)
	}
	d.next.ServeHTTP(rw, req)
}
