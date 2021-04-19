package addheader

import (
	"context"
	"fmt"
	"net/http"

	"github.com/casbin/casbin/v2"
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
	name     string
	enforcer *casbin.Enforcer
	headers  map[string]string
}

// New created a new Demo plugin.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	if len(config.Headers) == 0 {
		return nil, fmt.Errorf("headers cannot be empty")
	}

	enforcer, err := casbin.NewEnforcer("./model.conf", "./policy.csv")
	if err != nil {
		return nil, err
	}

	return &Demo{
		headers:  config.Headers,
		enforcer: enforcer,
		next:     next,
		name:     name,
	}, nil
}

func (d *Demo) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	for k, v := range d.headers {
		req.Header.Set(k, v)
	}

	ok, err := d.enforcer.Enforce("bob", "/dataset2/resource2", "GET")
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	if !ok {
		rw.WriteHeader(http.StatusUnauthorized)
		return
	}

	d.next.ServeHTTP(rw, req)
}
