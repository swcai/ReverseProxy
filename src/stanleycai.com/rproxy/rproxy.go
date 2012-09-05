package rproxy

import (
	"io"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/http/httputil"
	"time"
)

type Client struct {
	Addr   string
	Scheme string
}

var Clients []Client

type Request struct {
	done     chan bool
	request  http.Request
	response []byte
}

func selectClient() Client {
	// make it simiple. Occam's Razor ;)
	return Clients[rand.Int()%len(Clients)]
}

/* no matter what request, we update and forward it to web server instance */
func defaultHandler(w http.ResponseWriter, r *http.Request) {
	outreq := new(http.Request)
	*outreq = *r

	client := selectClient()
	outreq.URL.Scheme = client.Scheme
	outreq.URL.Host = client.Addr
	outreq.Host = outreq.URL.Host

	outreq.Proto = "HTTP/1.1"
	outreq.ProtoMajor = 1
	outreq.ProtoMinor = 1
	outreq.Close = false

	if outreq.Header.Get("Connectoin") != "" {
		outreq.Header = make(http.Header)
		copyHeader(outreq.Header, r.Header)
		delete(outreq.Header, "Connection")
	}

	if clientIp, _, err := net.SplitHostPort(outreq.RemoteAddr); err != nil {
		outreq.Header.Add("X-Forwarded-For", clientIp)
	}

	dmp, _ := httputil.DumpRequest(outreq, false)
	log.Println(string(dmp))

	transport := http.DefaultTransport
	res, err := transport.RoundTrip(outreq)
	if err != nil {
		log.Printf("proxy error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	copyHeader(w.Header(), res.Header)
	dmp2, _ := httputil.DumpResponse(res, false)
	log.Println(string(dmp2))

	w.WriteHeader(res.StatusCode)
	if res.Body != nil {
		var rw io.Writer = w
		io.Copy(rw, res.Body)
	}
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

func ListenAndServe(addr string, cs []string) {
	rand.Seed(time.Now().UnixNano())

	Clients = make([]Client, len(cs))
	for i, v := range cs {
		Clients[i] = Client{Addr: v, Scheme: "http"}
	}

	http.HandleFunc("/", defaultHandler)
	http.ListenAndServe(addr, nil)
}
