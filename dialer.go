package metadialer

import (
	"context"
	"crypto/rand"
	"fmt"
	"net"
	"net/url"
	"strings"

	"github.com/go-i2p/onramp"
)

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

// randomString generates a random string of 4 characters.
// It uses the crypto/rand package to generate a secure random byte slice,
// which is then converted to a string using a custom alphabet.
// The generated string is suitable for use as a unique identifier or token.
func randomString() string {
	// Define the alphabet for the random string
	const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	// Create a byte slice to hold the random bytes
	b := make([]byte, 4)
	// Generate secure random bytes
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	// Convert the random bytes to a string using the custom alphabet
	for i := range b {
		b[i] = alphabet[b[i]%byte(len(alphabet))]
	}
	return string(b)
}

// ANON is a flag to indicate whether to use the onion dialer for all non-I2P connections.
// If true, all non-I2P connections will be routed through the onion dialer.
// If false, regular connection will be made directly.
// Default is true.
var ANON = true

// Dial is a custom dialer that handles .i2p and .onion domains differently.
// It uses the garlic dialer for .i2p domains and the onion dialer for .onion domains.
// For all other domains, it uses the default dialer.
// If ANON is true, it will use the onion dialer for all non-I2P connections.
// It returns a net.Conn interface for the connection.
// If the address is invalid or the connection fails, it returns an error.
// The network parameter is ignored for onion connections.
var Dial = func(network, addr string) (net.Conn, error) {
	return dialHelper(network, addr)
}

// DialContext is a custom dialer that handles .i2p and .onion domains differently.
// It uses the garlic dialer for .i2p domains and the onion dialer for .onion domains.
// For all other domains, it uses the default dialer.
// If ANON is true, it will use the onion dialer for all non-I2P connections.
// It returns a net.Conn interface for the connection.
// If the address is invalid or the connection fails, it returns an error.
// The network parameter is ignored for onion connections.
// It accepts a context.Context parameter for cancellation and timeout.
// The context is ignored.
var DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
	return dialHelper(network, addr)
}

func dialHelper(network, addr string) (net.Conn, error) {
	// convert the addr to a URL
	tld, err := GetTLD(addr)
	if err != nil {
		return nil, err
	}
	switch tld {
	case "i2p":
		if GarlicErr != nil {
			return nil, GarlicErr
		}
		// I2P is a special case, we need to use the garlic dialer
		return Garlic.Dial(network, addr)
	case "onion":
		if OnionErr != nil {
			return nil, OnionErr
		}
		// make sure it is a TCP connection
		if network != "tcp" && network != "tcp4" && network != "tcp6" && network != "onion" {
			return nil, net.InvalidAddrError("only TCP connections are supported")
		}
		// Onion is a special case, we need to use the onion dialer
		return Onion.Dial("onion", addr)
	default:
		// If ANON is true, we need to use the onion dialer
		if ANON {
			if OnionErr != nil {
				return nil, OnionErr
			}
			// make sure it is a TCP connection
			if network != "tcp" {
				return nil, net.InvalidAddrError("only TCP connections are supported")
			}
			// ANON is a special case, we need to use the onion dialer
			return Onion.Dial(network, addr)
		} else {
			// For everything else, we can use the default dialer
			return net.Dial(network, addr)
		}
	}
}

// GetTLD is a helper function that returns the top-level domain of the given address.
// It takes a string address as input, which can be a fully qualified domain name or a URL.
// If the address does not include a scheme, "http://" is added by default.
// It returns the top-level domain as a string or an error if the address is invalid.
// The function also checks if the domain is an IP address and returns "ip" in that case.
// If there is no top-level domain found, it returns the entire domain as the TLD.
// The function is useful for determining the type of domain (I2P, onion, or regular) for routing purposes.
// It uses the net package to parse the address and extract the hostname.
func GetTLD(addr string) (string, error) {
	// Add a default scheme if missing
	if !strings.Contains(addr, "://") {
		addr = "http://" + addr
	}
	url, err := url.Parse(addr)
	if err != nil {
		return "", err
	}
	domain := url.Hostname()
	if domain == "" {
		return "", net.InvalidAddrError("invalid address: no hostname found")
	}
	// Check if the domain is an IP address
	if net.ParseIP(domain) != nil {
		return "ip", nil
	}
	lastDot := strings.LastIndex(domain, ".")
	if lastDot == -1 || lastDot == len(domain)-1 {
		return domain, nil
	}
	tld := domain[lastDot+1:]
	return tld, nil
}
