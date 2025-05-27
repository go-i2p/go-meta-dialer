package metadialer

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// MetaHTTPClient is an HTTP client that skips TLS verification for .i2p and .onion domains
// but performs standard verification for all other domains.
type MetaHTTPClient struct {
	*http.Client
}

// NewMetaHTTPClient creates a new client with special handling for .i2p and .onion domains.
// It accepts an optional root CA pool for custom certificate authorities.
func NewMetaHTTPClient(rootCAs *x509.CertPool) *MetaHTTPClient {
	// Create a custom transport with our special TLS config
	transport := &http.Transport{
		Dial: Dial,
		TLSClientConfig: &tls.Config{
			RootCAs: rootCAs, // May be nil, which will use system default
			VerifyConnection: func(state tls.ConnectionState) error {
				// Skip verification for .onion and .i2p domains
				domain := state.ServerName
				if strings.HasSuffix(domain, ".onion") || strings.HasSuffix(domain, ".i2p") {
					// Skip verification for these special domains
					return nil
				}

				// Use standard verification for all other domains
				if len(state.PeerCertificates) == 0 {
					return fmt.Errorf("no peer certificates provided")
				}

				opts := x509.VerifyOptions{
					DNSName:       state.ServerName,
					Intermediates: x509.NewCertPool(),
				}
				for _, cert := range state.PeerCertificates[1:] {
					opts.Intermediates.AddCert(cert)
				}
				_, err := state.PeerCertificates[0].Verify(opts)
				return err
			},
		},
		// Set reasonable defaults
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: false,
	}

	return &MetaHTTPClient{
		Client: &http.Client{
			Transport: transport,
			Timeout:   30 * time.Second,
		},
	}
}

// Get is a convenience method for making GET requests
func (c *MetaHTTPClient) Get(url string) (*http.Response, error) {
	return c.Client.Get(url)
}

// Post is a convenience method for making POST requests
func (c *MetaHTTPClient) Post(url, contentType string, body interface{}) (*http.Response, error) {
	// Convert the body interface{} to io.Reader
	var bodyReader io.Reader
	if body != nil {
		switch v := body.(type) {
		case io.Reader:
			bodyReader = v
		case []byte:
			bodyReader = bytes.NewReader(v)
		case string:
			bodyReader = strings.NewReader(v)
		default:
			// For other types, convert to string then to reader
			bodyReader = strings.NewReader(fmt.Sprintf("%v", v))
		}
	}
	return c.Client.Post(url, contentType, bodyReader)
}

// PostForm is a convenience method for making POST requests with form data
func (c *MetaHTTPClient) PostForm(url string, data url.Values) (*http.Response, error) {
	return c.Client.PostForm(url, data)
}

// Do is a convenience method for making arbitrary HTTP requests
func (c *MetaHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return c.Client.Do(req)
}

// Head is a convenience method for making HEAD requests
func (c *MetaHTTPClient) Head(url string) (*http.Response, error) {
	return c.Client.Head(url)
}

// CloseIdleConnections closes any idle connections in the transport
func (c *MetaHTTPClient) CloseIdleConnections() {
	c.Client.CloseIdleConnections()
}

// HTTPClient returns the underlying http.Client.
// This is useful for accessing the raw client if needed.
// It is anticipated that this will be necessary for some use cases.
// For example, to configure the default dialer for the application:
//
//	http.DefaultClient = client.HTTPClient()
func (c *MetaHTTPClient) HTTPClient() *http.Client {
	return c.Client
}

/*
// init initializes the MetaHTTPClient and sets it as the default HTTP client.
// This would be called when the application starts up, if not for being commented out here.
// It is commented out to avoid side effects during package initialization.
You can copy this code to your main package or wherever you want to initialize the client.
func init() {
	// Initialize the MetaHTTPClient with default settings
	client := NewMetaHTTPClient(nil)
	// Set the default HTTP client to our custom client
	http.DefaultClient = client.HTTPClient()
}
*/
