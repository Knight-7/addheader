package addheader

import (
	"context"
	"fmt"
	"net/http"
)

const (
	HeaderKey    = "Authorization"
	HeaderPrefix = "Bearer"
)

// Config the plugin configuration
type Config struct {
	Headers map[string]string
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

	auth := req.Header.Get(HeaderKey)
	c, err := ParseToken(auth[len(HeaderPrefix)+1:])
	if err != nil {
		req.Header.Set("err", err.Error())
	} else {
		req.Header.Set("username", c.Username)
		req.Header.Set("passwd", c.Passwd)
		req.Header.Set("role", c.Role)
	}

	d.next.ServeHTTP(rw, req)
}
