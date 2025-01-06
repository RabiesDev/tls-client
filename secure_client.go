package tls_client

import (
	"crypto/tls"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"

	oohttp "github.com/ooni/oohttp"
)

var defaultDialer = &net.Dialer{
	Timeout:   30 * time.Second,
	KeepAlive: 30 * time.Second,
}

type TLSConnFactory func(conn net.Conn, config *tls.Config) oohttp.TLSConn

type ClientOptions struct {
	TLSFactory *SecureTLSFactory
	ProxyURL   *url.URL

	ForceAttemptHTTP2 bool
	MaxIdleConns      int
	Timeout           time.Duration
}

func DefaultOptions() *ClientOptions {
	return &ClientOptions{
		TLSFactory: &SecureTLSFactory{},
		ProxyURL:   nil,

		ForceAttemptHTTP2: true,
		MaxIdleConns:      100,
		Timeout:           time.Second * 20,
	}
}

func MustNewSecureClient(opts *ClientOptions) *http.Client {
	client, err := NewSecureClient(opts)
	if err != nil {
		panic(err)
	}
	return client
}

func NewSecureClient(opts *ClientOptions) (*http.Client, error) {
	if opts == nil {
		opts = DefaultOptions()
	}

	return NewClientWithCustomTransport(NewSecureTransport(*opts), *opts)
}

func NewClientWithCustomTransport(transport *oohttp.Transport, opts ClientOptions) (*http.Client, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	return &http.Client{
		Transport: &oohttp.StdlibTransport{
			Transport: transport,
		},
		CheckRedirect: nil,
		Timeout:       opts.Timeout,
		Jar:           jar,
	}, nil
}

func NewSecureTransport(opts ClientOptions) *oohttp.Transport {
	return &oohttp.Transport{
		Proxy: func(request *oohttp.Request) (*url.URL, error) {
			if opts.ProxyURL != nil {
				return opts.ProxyURL, nil
			}
			return oohttp.ProxyFromEnvironment(request)
		},
		DialContext:           defaultDialer.DialContext,
		TLSClientFactory:      opts.TLSFactory.NewTLSConnection,
		ForceAttemptHTTP2:     opts.ForceAttemptHTTP2,
		MaxIdleConns:          opts.MaxIdleConns,
		TLSHandshakeTimeout:   10 * time.Second,
		IdleConnTimeout:       90 * time.Second,
		ExpectContinueTimeout: time.Second,
	}
}
