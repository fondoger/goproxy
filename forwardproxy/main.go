package forwardproxy

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"

	"time"

	"github.com/elazarl/goproxy"
)

var client = &http.Client{}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	var addr = flag.String("addr", "", "Listen host port, eg: --addr=0.0.0.0:8000")
	var httpRelay = flag.String("http-relay", "", "(optional) Relay to http proxy, eg: --http-relay=http://127.0.0.1:8081")
	flag.Usage = func() {
		_, _ = fmt.Fprintf(os.Stderr, "Usage: goproxy --addr=0.0.0.0:8080 \n\nOptions:\n")
		flag.PrintDefaults()
	}
	flag.Parse()
	if *addr == "" {
		flag.Usage()
		return
	}

	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = true

	if *httpRelay != "" {
		proxyUrl, err := url.Parse(*httpRelay)
		if err != nil {
			log.Fatalf("failed parse http relay, err %v", err)
		}
		proxy.Tr.Proxy = http.ProxyURL(proxyUrl)
		log.Printf("Serve as a http proxy relay to %v", *httpRelay)
	}

	log.Println("Starting proxy server on", *addr)
	PrintPublicIp()
	PrintLocalIp()
	if err := http.ListenAndServe(*addr, proxy); err != nil {
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
