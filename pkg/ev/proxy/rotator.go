package proxy

import (
	"net/http"
	"net/url"
)

type ProxyFunc func(*http.Request) (*url.URL, error)
