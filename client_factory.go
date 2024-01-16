package tls_client

import (
	"crypto/tls"
	oohttp "github.com/ooni/oohttp"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"
)

// DefaultNetDialer is the default [net.Dialer].
var DefaultNetDialer = &net.Dialer{
	Timeout:   30 * time.Second,
	KeepAlive: 30 * time.Second,
}

type TLSFactoryFunc func(conn net.Conn, config *tls.Config) oohttp.TLSConn

func NewHttpClient(proxy *url.URL, factory *FactoryWithParrot) (*http.Client, error) {
	if factory == nil {
		factory = &FactoryWithParrot{}
	}
	return NewHttpClientWithTransport(NewTransport(proxy, factory.NewTLSConn))
}

func NewHttpClientWithTransport(transport *oohttp.Transport) (*http.Client, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	return &http.Client{
		Transport: &oohttp.StdlibTransport{
			Transport: transport,
		},
		Timeout:       time.Second * 20,
		CheckRedirect: nil,
		Jar:           jar,
	}, nil
}

func NewTransport(proxy *url.URL, tlsFactoryFunc TLSFactoryFunc) *oohttp.Transport {
	return &oohttp.Transport{
		Proxy: func(request *oohttp.Request) (*url.URL, error) {
			if proxy != nil {
				return proxy, nil
			}
			return oohttp.ProxyFromEnvironment(request)
		},
		DialContext:           DefaultNetDialer.DialContext,
		TLSClientFactory:      tlsFactoryFunc,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		TLSHandshakeTimeout:   10 * time.Second,
		IdleConnTimeout:       90 * time.Second,
		ExpectContinueTimeout: time.Second,
	}
}
