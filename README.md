# go-meta-dialer

A dialer which will reach clearnet, onion, and I2P sites, and an HTTP Client which has specific TLS behavior for Onion and I2P domains.

## Installation

```bash
go get github.com/go-i2p/go-meta-dialer
```

## Features

- Multi-network dialer supporting regular internet, Tor (.onion), and I2P (.i2p) connections
- Automatic routing based on destination address
- Optional anonymity mode to route all non-I2P connections through Tor
- HTTP client with special TLS handling for darknet services

## Usage

### Basic Dialer

```go
package main

import (
    "fmt"
    "io/ioutil"
    "net/http"
    
    metadialer "github.com/go-i2p/go-meta-dialer"
)

func main() {
    // IMPORTANT: Close the dialers when done to prevent SAM/Tor session leaks
    defer metadialer.Garlic.Close()
    defer metadialer.Onion.Close()
    
    // Use the dialer directly
    conn, err := metadialer.Dial("tcp", "example.i2p:80")
    if err != nil {
        panic(err)
    }
    defer conn.Close()
    
    // Do something with the connection...
}
```

### HTTP Client

```go
package main

import (
    "fmt"
    "io/ioutil"
    
    metadialer "github.com/go-i2p/go-meta-dialer"
)

func main() {
    // IMPORTANT: Close the dialers when done
    defer metadialer.Garlic.Close()
    defer metadialer.Onion.Close()
    
    // Create a new HTTP client
    client := metadialer.NewMetaHTTPClient(nil)
    
    // Make requests to any network
    resp, err := client.Get("https://example.com")           // Normal website
    // resp, err := client.Get("http://example.onion")       // Tor onion service
    // resp, err := client.Get("http://example.i2p")         // I2P site
    
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()
    
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        panic(err)
    }
    
    fmt.Println(string(body))
}
```

### Make it the Default HTTP Client

```go
// Set the MetaHTTPClient as the default for all HTTP requests
client := metadialer.NewMetaHTTPClient(nil)
http.DefaultClient = client.HTTPClient()
```

## Configuration

- `metadialer.ANON = true`: (Default) Routes all non-I2P connections through Tor
- `metadialer.ANON = false`: Only routes .onion domains through Tor, direct connection for regular domains

## Notes

- TLS verification is skipped for .i2p and .onion domains but performs standard verification for all other domains
- The dialer automatically routes traffic to the appropriate network based on the TLD
- Requires a running SAM bridge (default: 127.0.0.1:7656) for I2P connections
- Requires a running Tor SOCKS proxy for .onion connections