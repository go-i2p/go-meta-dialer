# metadialer
--
    import "github.com/go-i2p/go-meta-dialer"


## Usage

```go
var (
	// Garlic and Onion are the dialers for I2P and onion connections respectively.
	// Garlic is used for I2P connections and Onion is used for onion connections.
	// GarlicErr and OnionErr are the errors returned by the dialers.
	// It is important to `defer close` the dialers when you include them in your code.
	// Otherwise your SAMv3 or Tor sessions may leak. onramp tries to fix it for you but do it anyway.
	// in your `main` function, do:
	// defer Garlic.Close()
	// defer Onion.Close()
	Garlic, GarlicErr = onramp.NewGarlic(fmt.Sprintf("metadialer-%s", randomString()), "127.0.0.1:7656", onramp.OPT_DEFAULTS)
	Onion, OnionErr   = onramp.NewOnion(fmt.Sprintf("metadialer-%s", randomString()))
)
```

```go
var ANON = true
```
ANON is a flag to indicate whether to use the onion dialer for all non-I2P
connections. If true, all non-I2P connections will be routed through the onion
dialer. If false, regular connection will be made directly. Default is true.

```go
var Dial = func(network, addr string) (net.Conn, error) {
	return dialHelper(network, addr)
}
```
Dial is a custom dialer that handles .i2p and .onion domains differently. It
uses the garlic dialer for .i2p domains and the onion dialer for .onion domains.
For all other domains, it uses the default dialer. If ANON is true, it will use
the onion dialer for all non-I2P connections. It returns a net.Conn interface
for the connection. If the address is invalid or the connection fails, it
returns an error. The network parameter is ignored for onion connections.

```go
var DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
	return dialHelper(network, addr)
}
```
DialContext is a custom dialer that handles .i2p and .onion domains differently.
It uses the garlic dialer for .i2p domains and the onion dialer for .onion
domains. For all other domains, it uses the default dialer. If ANON is true, it
will use the onion dialer for all non-I2P connections. It returns a net.Conn
interface for the connection. If the address is invalid or the connection fails,
it returns an error. The network parameter is ignored for onion connections. It
accepts a context.Context parameter for cancellation and timeout. The context is
ignored.

#### func  GetTLD

```go
func GetTLD(addr string) (string, error)
```
GetTLD is a helper function that returns the top-level domain of the given
address. It takes a string address as input, which can be a fully qualified
domain name or a URL. If the address does not include a scheme, "http://" is
added by default. It returns the top-level domain as a string or an error if the
address is invalid. The function also checks if the domain is an IP address and
returns "ip" in that case. If there is no top-level domain found, it returns the
entire domain as the TLD. The function is useful for determining the type of
domain (I2P, onion, or regular) for routing purposes. It uses the net package to
parse the address and extract the hostname.

#### type MetaHTTPClient

```go
type MetaHTTPClient struct {
	*http.Client
}
```

MetaHTTPClient is an HTTP client that skips TLS verification for .i2p and .onion
domains but performs standard verification for all other domains.

#### func  NewMetaHTTPClient

```go
func NewMetaHTTPClient(rootCAs *x509.CertPool) *MetaHTTPClient
```
NewMetaHTTPClient creates a new client with special handling for .i2p and .onion
domains. It accepts an optional root CA pool for custom certificate authorities.

#### func (*MetaHTTPClient) CloseIdleConnections

```go
func (c *MetaHTTPClient) CloseIdleConnections()
```
CloseIdleConnections closes any idle connections in the transport

#### func (*MetaHTTPClient) Do

```go
func (c *MetaHTTPClient) Do(req *http.Request) (*http.Response, error)
```
Do is a convenience method for making arbitrary HTTP requests

#### func (*MetaHTTPClient) Get

```go
func (c *MetaHTTPClient) Get(url string) (*http.Response, error)
```
Get is a convenience method for making GET requests

#### func (*MetaHTTPClient) HTTPClient

```go
func (c *MetaHTTPClient) HTTPClient() *http.Client
```
HTTPClient returns the underlying http.Client. This is useful for accessing the
raw client if needed. It is anticipated that this will be necessary for some use
cases. For example, to configure the default dialer for the application:

    http.DefaultClient = client.HTTPClient()

#### func (*MetaHTTPClient) Head

```go
func (c *MetaHTTPClient) Head(url string) (*http.Response, error)
```
Head is a convenience method for making HEAD requests

#### func (*MetaHTTPClient) Post

```go
func (c *MetaHTTPClient) Post(url, contentType string, body interface{}) (*http.Response, error)
```
Post is a convenience method for making POST requests

#### func (*MetaHTTPClient) PostForm

```go
func (c *MetaHTTPClient) PostForm(url string, data url.Values) (*http.Response, error)
```
PostForm is a convenience method for making POST requests with form data
