package reverseproxy

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"time"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	var addr = flag.String("addr", "", "Listen host port, eg: --addr=0.0.0.0:8000")
	var httpRelay = flag.String("http-relay", "", " Relay to http proxy, eg: --http-relay=http://127.0.0.1:8081")
	flag.Usage = func() {
		_, _ = fmt.Fprintf(os.Stderr, "Usage: goproxy --addr=0.0.0.0:8080 \n\nOptions:\n")
		flag.PrintDefaults()
	}
	flag.Parse()
	if *addr == "" {
		flag.Usage()
		return
	}

	// Create a reverse proxy
	targetUrl, _ := url.Parse(*httpRelay)
	proxy := httputil.NewSingleHostReverseProxy(targetUrl)

	// Configure the reverse proxy to handle TLS forwarding between the client and the target server
	proxy.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	// Modify the request Host header
	proxy.Director = func(req *http.Request) {
		req.URL.Scheme = targetUrl.Scheme
		req.URL.Host = targetUrl.Host
		req.Host = targetUrl.Host
	}

	// Start the HTTP server
	server := &http.Server{
		Addr:    *addr,
		Handler: proxy,
	}

	log.Println("Starting proxy server on", *addr)
	PrintPublicIp()
	PrintLocalIp()
	if err := server.ListenAndServe(); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

// Get preferred outbound ip of this machine
func PrintLocalIp() {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	log.Printf("local ip address: %v", localAddr.IP.String())
}

func PrintPublicIp() {
	var ctx, cancel = context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	result := make(chan bool)
	go func() {
		resp, err := http.Get("http://httpbin.org/ip")
		if err != nil {
			log.Printf("get public ip error: %v", err)
		}
		buf, _ := ioutil.ReadAll(resp.Body)
		var IP = struct {
			Origin string
		}{}
		_ = json.Unmarshal(buf, &IP)
		log.Printf("public ip address: %v", IP.Origin)
		result <- true
	}()
	select {
	case <-ctx.Done():
		if ctx.Err() != nil {
			log.Printf("get public ip timeout.")
		}
	case <-result:
	}
}
