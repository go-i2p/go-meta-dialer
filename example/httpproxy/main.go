package main

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"net"
	"net/http"

	metadialer "github.com/go-i2p/go-meta-dialer"
)

// This example demonstrates how to use the MetaHTTPClient to handle HTTP requests
// over I2P and Tor networks. It sets up a simple HTTP proxy server that listens
// for incoming connections and forwards them to the appropriate dialer based on
// the URL scheme (I2P or Tor). The response is then sent back to the client.
// The MetaHTTPClient is a custom HTTP client that uses the MetaDialer to handle
// connections to I2P and Tor networks.

func main() {
	// Example usage of the MetaHTTPClient
	client := metadialer.NewMetaHTTPClient(nil)
	proxyListener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	defer proxyListener.Close()
	// Start the proxy server
	for {
		conn, err := proxyListener.Accept()
		if err != nil {
			continue
		}
		go func(c net.Conn) {
			defer c.Close()
			// read the request from the connection
			b, err := io.ReadAll(c)
			if err != nil {
				log.Println("Error reading request:", err)
				return
			}
			// parse the request
			oldreq, err := http.ReadRequest(bufio.NewReader(bytes.NewReader(b)))
			if err != nil {
				log.Println("Error parsing request:", err)
				return
			}
			// obtain the method and URL from the request
			method := oldreq.Method
			url := oldreq.URL.String()
			// recreate the body
			body := bytes.NewBuffer(b)
			// print the method and URL
			log.Println("Method:", method)
			log.Println("URL:", url)

			// Handle the connection using the MetaHTTPClient
			req, err := http.NewRequest(method, url, body)
			if err != nil {
				log.Println("Error creating request:", err)
				return
			}
			resp, err := client.Do(req)
			if err != nil {
				log.Println("Error making request:", err)
				return
			}
			defer resp.Body.Close()
			// get the response body
			respBody, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Println("Error reading response body:", err)
				return
			}
			// write the response back to the connection
			_, err = c.Write(respBody)
			if err != nil {
				log.Println("Error writing response:", err)
				return
			}
			// close the connection
			c.Close()
			// print the response status
			log.Println("Response Status:", resp.Status)
		}(conn)
	}
}
